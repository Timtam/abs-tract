package thalia

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// Search searches for a book by title and optionally reorders results by author similarity.
func (c *Client) Search(ctx context.Context, title string, author *string) ([]Book, error) {
	normalisedTitle := normaliseString(title)
	normalisedAuthor := normaliseString(lo.FromPtr(author))

	query := strings.TrimSpace(title)
	switch {
	case normalisedTitle != "":
	case normalisedAuthor != "":
		query = lo.FromPtr(author)
	default:
		return nil, nil
	}

	books, err := c.searchBooks(ctx, query)
	if err != nil {
		return nil, err
	}

	if normalisedAuthor != "" {
		sortBooksByAuthorSimilarity(books, normalisedAuthor)
	}

	if len(books) > 20 {
		books = books[:20]
	}

	return c.enrichBooks(ctx, books), nil
}

func (c *Client) searchBooks(ctx context.Context, query string) ([]Book, error) {
	htmlResponse, err := c.get(ctx, "suche", map[string]string{
		"sq":  query,
		"asn": "true",
	})
	if err != nil {
		return nil, err
	}

	return BooksFromHTML(htmlResponse)
}

func (c *Client) enrichBooks(ctx context.Context, books []Book) []Book {
	enrichedBooks := make([]Book, len(books))
	copy(enrichedBooks, books)

	var wg sync.WaitGroup
	for idx, book := range books {
		wg.Add(1)

		go func(idx int, book Book) {
			defer wg.Done()

			detailedBook, err := c.getBook(ctx, book.Path)
			if err != nil {
				return
			}

			enrichedBooks[idx] = mergeBook(book, detailedBook)
		}(idx, book)
	}

	wg.Wait()

	return enrichedBooks
}

func (c *Client) getBook(ctx context.Context, path string) (Book, error) {
	htmlResponse, err := c.get(ctx, path, nil)
	if err != nil {
		return Book{}, err
	}

	book, err := BookDetailsFromHTML(htmlResponse)
	if err != nil {
		return Book{}, err
	}
	if book == nil {
		return Book{}, nil
	}

	return *book, nil
}

func mergeBook(base Book, detail Book) Book {
	if detail.ID != "" {
		base.ID = detail.ID
	}
	if detail.Path != "" {
		base.Path = detail.Path
	}
	if detail.Title != "" {
		base.Title = detail.Title
	}
	if detail.Subtitle != "" {
		base.Subtitle = detail.Subtitle
	}
	if detail.Author != "" {
		base.Author = detail.Author
	}
	if detail.Narrator != "" {
		base.Narrator = detail.Narrator
	}
	if detail.Abridged != nil {
		base.Abridged = detail.Abridged
	}
	base.Format = mergeFormat(base.Format, detail.Format)
	if detail.Cover != "" {
		base.Cover = detail.Cover
	}
	if detail.Description != "" {
		base.Description = detail.Description
	}
	if detail.Publisher != "" {
		base.Publisher = detail.Publisher
	}
	if detail.Language != "" {
		base.Language = detail.Language
	}
	if detail.ISBN != "" {
		base.ISBN = detail.ISBN
	}
	if detail.Duration != nil {
		base.Duration = detail.Duration
	}
	if detail.Series != "" {
		base.Series = detail.Series
	}
	if detail.Sequence != "" {
		base.Sequence = detail.Sequence
	}
	if len(detail.Tags) != 0 {
		base.Tags = detail.Tags
	}
	if detail.PublishDate != nil {
		base.PublishDate = detail.PublishDate
	}

	return base
}

func mergeFormat(baseFormat string, detailFormat string) string {
	switch {
	case detailFormat == "":
		return baseFormat
	case baseFormat == "":
		return detailFormat
	}

	baseFormatLower := strings.ToLower(baseFormat)
	detailFormatLower := strings.ToLower(detailFormat)
	switch {
	case strings.Contains(baseFormatLower, detailFormatLower):
		return baseFormat
	case strings.Contains(detailFormatLower, baseFormatLower):
		return detailFormat
	case len(detailFormat) > len(baseFormat):
		return detailFormat
	default:
		return baseFormat
	}
}

func sortBooksByAuthorSimilarity(books []Book, author string) {
	normalisedAuthor := normaliseString(author)

	slices.SortStableFunc(books, func(i, j Book) int {
		authorSimilarity1 := strutil.Similarity(
			normaliseString(i.Author),
			normalisedAuthor,
			metrics.NewJaroWinkler(),
		)
		authorSimilarity2 := strutil.Similarity(
			normaliseString(j.Author),
			normalisedAuthor,
			metrics.NewJaroWinkler(),
		)
		if authorSimilarity1 > authorSimilarity2 {
			return -1
		}
		if authorSimilarity1 < authorSimilarity2 {
			return 1
		}
		return 0
	})
}

var (
	spaceRegex        = regexp.MustCompile(`\s+`)
	alphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]`)
)

func normaliseString(s string) string {
	s = spaceRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = alphanumericRegex.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "the ")

	return s
}
