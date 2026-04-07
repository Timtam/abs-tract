package thalia

import (
	"encoding/json"
	stdhtml "html"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	mapset "github.com/deckarep/golang-set/v2"
	nethtml "golang.org/x/net/html"
)

const publishDateLayout = "02.01.2006"

var (
	searchResultsExpr = xpath.MustCompile(`//li[contains(concat(" ", normalize-space(@class), " "), " tm-produktliste__eintrag ")]`)
	searchLinkExpr    = xpath.MustCompile(`.//a[contains(concat(" ", normalize-space(@class), " "), " tm-produkt-link ")]`)
	searchTitleExpr   = xpath.MustCompile(`.//strong[contains(concat(" ", normalize-space(@class), " "), " tm-artikeldetails__titel ")]`)
	searchAuthorExpr  = xpath.MustCompile(`.//p[contains(concat(" ", normalize-space(@class), " "), " tm-artikeldetails__autor ")]`)
	searchFormatExpr  = xpath.MustCompile(`.//p[contains(concat(" ", normalize-space(@class), " "), " tm-artikeldetails__formatbezeichnung ")]`)
	searchCoverExpr   = xpath.MustCompile(`.//img[contains(concat(" ", normalize-space(@class), " "), " tm-artikelbild-wrapper__artikelbild ")]`)

	detailJSONLDExpr         = xpath.MustCompile(`//script[@type="application/ld+json"]`)
	detailAuthorExpr         = xpath.MustCompile(`//a[contains(concat(" ", normalize-space(@class), " "), " autor-name ")]`)
	detailSeriesExpr         = xpath.MustCompile(`//a[contains(concat(" ", normalize-space(@class), " "), " serie-link ")]`)
	detailSpeakerExpr        = xpath.MustCompile(`//a[contains(concat(" ", normalize-space(@class), " "), " sprecher-name ")]`)
	detailSubtitleExpr       = xpath.MustCompile(`//*[contains(concat(" ", normalize-space(@class), " "), " untertitel ")]`)
	detailDescriptionExpr    = xpath.MustCompile(`//template[@data-id="zusatztexte"]//*[contains(concat(" ", normalize-space(@class), " "), " zusatztexte ")]`)
	detailShortDescExpr      = xpath.MustCompile(`//*[contains(concat(" ", normalize-space(@class), " "), " kurzbeschreibung ")]`)
	detailMetaDescExpr       = xpath.MustCompile(`//meta[@property="og:description" or @name="description"]`)
	detailPageViewExpr       = xpath.MustCompile(`//dl-pageview`)
	detailProductDetailsExpr = xpath.MustCompile(`//div[contains(concat(" ", normalize-space(@class), " "), " details-default ")]//section[contains(concat(" ", normalize-space(@class), " "), " artikeldetail ")]`)
	detailLabelExpr          = xpath.MustCompile(`.//*[contains(concat(" ", normalize-space(@class), " "), " detailbezeichnung ")]`)
	detailValueExpr          = xpath.MustCompile(`.//*[contains(concat(" ", normalize-space(@class), " "), " value ")]`)
)

type Book struct {
	ID          string
	Path        string
	Title       string
	Subtitle    string
	Author      string
	Narrator    string
	Format      string
	Cover       string
	Description string
	Publisher   string
	Language    string
	ISBN        string
	Duration    *int
	Series      string
	Sequence    string
	Tags        []string
	PublishDate *time.Time
}

type metadataJSON struct {
	Type        string   `json:"@type"`
	Name        string   `json:"name"`
	Image       []string `json:"image"`
	Description string   `json:"description"`
	ISBN        string   `json:"isbn"`
	GTIN13      string   `json:"gtin13"`
}

// BooksFromHTML parses the books from the html of a Thalia search results page.
func BooksFromHTML(searchNode *nethtml.Node) ([]Book, error) {
	resultNodes := htmlquery.QuerySelectorAll(searchNode, searchResultsExpr)

	books := make([]Book, 0, len(resultNodes))
	seenIDs := mapset.NewSet[string]()
	for _, resultNode := range resultNodes {
		book := BookFromHTML(resultNode)
		if book == nil || !isSupportedFormat(book.Format) || seenIDs.Contains(book.ID) {
			continue
		}

		books = append(books, *book)
		seenIDs.Add(book.ID)
	}

	return books, nil
}

// BookFromHTML parses a book from a Thalia search result node.
// If a result is not a book, nil is returned.
func BookFromHTML(bookNode *nethtml.Node) *Book {
	linkNode := htmlquery.QuerySelector(bookNode, searchLinkExpr)
	path := htmlquery.SelectAttr(linkNode, "href")
	if path == "" {
		return nil
	}

	id := strings.TrimPrefix(path, "/shop/home/artikeldetails/")
	if id == "" {
		return nil
	}

	title := cleanText(htmlquery.InnerText(htmlquery.QuerySelector(bookNode, searchTitleExpr)))
	if title == "" {
		return nil
	}

	format := searchFormat(bookNode)
	if format == "" {
		return nil
	}

	coverNode := htmlquery.QuerySelector(bookNode, searchCoverExpr)

	return &Book{
		ID:     id,
		Path:   path,
		Title:  title,
		Author: cleanText(htmlquery.InnerText(htmlquery.QuerySelector(bookNode, searchAuthorExpr))),
		Format: format,
		Cover:  htmlquery.SelectAttr(coverNode, "src"),
	}
}

// BookDetailsFromHTML parses a book detail page.
func BookDetailsFromHTML(detailNode *nethtml.Node) (*Book, error) {
	book := &Book{}

	if metadata := metadataFromHTML(detailNode); metadata != nil {
		book.Title = cleanText(metadata.Name)
		book.Description = cleanText(metadata.Description)
		book.ISBN = cleanText(metadata.ISBN)
		if book.ISBN == "" {
			book.ISBN = cleanText(metadata.GTIN13)
		}
		if len(metadata.Image) != 0 {
			book.Cover = cleanText(metadata.Image[0])
		}
	}
	if description := descriptionFromHTML(detailNode); description != "" {
		book.Description = description
	}

	book.Subtitle = nodeText(htmlquery.QuerySelector(detailNode, detailSubtitleExpr))

	authors := authorNamesFromHTML(detailNode)
	if len(authors) != 0 {
		book.Author = strings.Join(authors, ", ")
	}

	speakers := speakerNamesFromHTML(detailNode)
	if len(speakers) != 0 {
		book.Narrator = strings.Join(speakers, ", ")
	}

	book.Series = nodeText(htmlquery.QuerySelector(detailNode, detailSeriesExpr))

	pageView := htmlquery.QuerySelector(detailNode, detailPageViewExpr)
	if pageView != nil {
		if book.Title == "" {
			book.Title = cleanText(htmlquery.SelectAttr(pageView, "produkttitel"))
		}
		if book.Cover == "" {
			book.Cover = cleanText(htmlquery.SelectAttr(pageView, "coverbild"))
		}
		if book.Publisher == "" {
			book.Publisher = cleanText(htmlquery.SelectAttr(pageView, "hersteller"))
		}
		if book.Format == "" {
			book.Format = cleanText(htmlquery.SelectAttr(pageView, "form"))
		}
		if book.ISBN == "" {
			book.ISBN = cleanText(htmlquery.SelectAttr(pageView, "ean"))
		}
	}

	productDetails := productDetailsFromHTML(detailNode)
	if value, ok := productDetails["Einband"]; ok {
		book.Format = value
	}
	if value, ok := productDetails["Gesprochen von"]; ok && book.Narrator == "" {
		book.Narrator = value
	}
	if value, ok := productDetails["Spieldauer"]; ok {
		book.Duration = parseDuration(value)
	}
	if value, ok := productDetails["Erscheinungsdatum"]; ok {
		publishDate, err := time.Parse(publishDateLayout, value)
		if err == nil {
			book.PublishDate = &publishDate
		}
	}
	if value, ok := productDetails["Verlag"]; ok {
		book.Publisher = value
	}
	if value, ok := productDetails["Sprache"]; ok {
		book.Language = value
	}
	if value, ok := productDetails["ISBN"]; ok {
		book.ISBN = value
	}
	if value, ok := productDetails["EAN"]; ok && book.ISBN == "" {
		book.ISBN = value
	}

	book.Tags = audioTagsFromProductDetails(productDetails)
	if book.Narrator == "" {
		book.Narrator = narratorFromSubtitle(book.Subtitle)
	}
	if book.Series != "" {
		book.Sequence = seriesSequenceFromTitle(book.Title)
	}

	return book, nil
}

func metadataFromHTML(detailNode *nethtml.Node) *metadataJSON {
	var productMetadata *metadataJSON

	for _, metadataNode := range htmlquery.QuerySelectorAll(detailNode, detailJSONLDExpr) {
		var metadata metadataJSON
		err := json.Unmarshal([]byte(htmlquery.InnerText(metadataNode)), &metadata)
		if err != nil {
			continue
		}
		if metadata.Type == "Book" {
			return &metadata
		}
		if metadata.Type == "Product" && productMetadata == nil {
			productMetadata = &metadata
		}
	}

	return productMetadata
}

func authorNamesFromHTML(detailNode *nethtml.Node) []string {
	authorNodes := htmlquery.QuerySelectorAll(detailNode, detailAuthorExpr)
	return nodeTexts(authorNodes)
}

func speakerNamesFromHTML(detailNode *nethtml.Node) []string {
	speakerNodes := htmlquery.QuerySelectorAll(detailNode, detailSpeakerExpr)
	return nodeTexts(speakerNodes)
}

func nodeTexts(nodes []*nethtml.Node) []string {
	values := make([]string, 0, len(nodes))
	seenValues := mapset.NewSet[string]()
	for _, node := range nodes {
		value := nodeText(node)
		if value == "" || seenValues.Contains(value) {
			continue
		}

		values = append(values, value)
		seenValues.Add(value)
	}

	return values
}

func nodeText(node *nethtml.Node) string {
	if node == nil {
		return ""
	}

	return cleanText(htmlquery.InnerText(node))
}

func nodeTextWithBreaks(node *nethtml.Node) string {
	if node == nil {
		return ""
	}

	var builder strings.Builder
	var walk func(*nethtml.Node)
	walk = func(current *nethtml.Node) {
		switch current.Type {
		case nethtml.TextNode:
			builder.WriteString(current.Data)
		case nethtml.ElementNode:
			if current.Data == "br" {
				builder.WriteByte(' ')
			}
		}

		for child := current.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}

		if current.Type == nethtml.ElementNode {
			switch current.Data {
			case "div", "p", "li", "section":
				builder.WriteByte(' ')
			}
		}
	}

	walk(node)

	return cleanText(builder.String())
}

func descriptionFromHTML(detailNode *nethtml.Node) string {
	descriptionNodes := []*nethtml.Node{
		htmlquery.QuerySelector(detailNode, detailDescriptionExpr),
		htmlquery.QuerySelector(detailNode, detailShortDescExpr),
	}
	for _, descriptionNode := range descriptionNodes {
		description := nodeTextWithBreaks(descriptionNode)
		if description != "" {
			return description
		}
	}

	metaDescriptionNode := htmlquery.QuerySelector(detailNode, detailMetaDescExpr)
	return cleanText(htmlquery.SelectAttr(metaDescriptionNode, "content"))
}

func productDetailsFromHTML(detailNode *nethtml.Node) map[string]string {
	productDetailNodes := htmlquery.QuerySelectorAll(detailNode, detailProductDetailsExpr)
	productDetails := make(map[string]string, len(productDetailNodes))
	for _, productDetailNode := range productDetailNodes {
		label := cleanText(htmlquery.InnerText(htmlquery.QuerySelector(productDetailNode, detailLabelExpr)))
		if label == "" {
			continue
		}

		valueNodes := htmlquery.QuerySelectorAll(productDetailNode, detailValueExpr)
		values := nodeTexts(valueNodes)
		if len(values) == 0 {
			continue
		}

		productDetails[label] = strings.Join(values, ", ")
	}

	return productDetails
}

func searchFormat(bookNode *nethtml.Node) string {
	formatNode := htmlquery.QuerySelector(bookNode, searchFormatExpr)
	if formatNode == nil {
		return ""
	}

	spanNodes := htmlquery.Find(formatNode, "./span")
	if len(spanNodes) != 0 {
		return cleanText(htmlquery.InnerText(spanNodes[0]))
	}

	format := cleanText(htmlquery.InnerText(formatNode))
	format = strings.TrimSuffix(format, "+ weitere")
	return cleanText(format)
}

func isSupportedFormat(format string) bool {
	switch {
	case strings.HasPrefix(format, "Buch"):
		return true
	case strings.HasPrefix(format, "Hörbuch"):
		return true
	case strings.HasPrefix(format, "Hörbuch-Download"):
		return true
	default:
		return false
	}
}

var (
	hourRegex     = regexp.MustCompile(`(?i)(\d+)\s*(?:stunde|stunden|std|h)\.?\b`)
	minuteRegex   = regexp.MustCompile(`(?i)(\d+)\s*(?:minute|minuten|min)\.?\b`)
	narratorRegex = regexp.MustCompile(`(?i)\b(?:mit|gelesen von|gesprochen von)\s+(.+?)(?:\s*\([^)]*\))?$`)
	sequenceRegex = regexp.MustCompile(`(?i)^(?:band|folge|teil|episode)\s+(\d+)\b`)
)

func parseDuration(value string) *int {
	value = cleanText(value)
	if value == "" {
		return nil
	}

	var durationSeconds int
	if hourMatches := hourRegex.FindStringSubmatch(value); len(hourMatches) == 2 {
		hours, err := strconv.Atoi(hourMatches[1])
		if err == nil {
			durationSeconds += hours * 3600
		}
	}
	if minuteMatches := minuteRegex.FindStringSubmatch(value); len(minuteMatches) == 2 {
		minutes, err := strconv.Atoi(minuteMatches[1])
		if err == nil {
			durationSeconds += minutes * 60
		}
	}

	if durationSeconds == 0 {
		return nil
	}

	return &durationSeconds
}

func audioTagsFromProductDetails(productDetails map[string]string) []string {
	tagKeys := []string{
		"Hörtyp",
		"Fassung",
		"Medium",
		"Family Sharing",
		"Abo-Fähigkeit",
		"Anzahl Dateien",
		"Anzahl",
		"Altersempfehlung",
		"Übersetzt von",
	}

	tags := make([]string, 0, len(tagKeys))
	for _, key := range tagKeys {
		value, ok := productDetails[key]
		if !ok || value == "" {
			continue
		}
		tags = append(tags, key+": "+value)
	}

	return tags
}

func narratorFromSubtitle(subtitle string) string {
	subtitle = cleanText(subtitle)
	if subtitle == "" {
		return ""
	}

	matches := narratorRegex.FindStringSubmatch(subtitle)
	if len(matches) != 2 {
		return ""
	}

	return cleanText(matches[1])
}

func seriesSequenceFromTitle(title string) string {
	matches := sequenceRegex.FindStringSubmatch(title)
	if len(matches) != 2 {
		return ""
	}

	return matches[1]
}

func cleanText(s string) string {
	s = stdhtml.UnescapeString(s)
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
