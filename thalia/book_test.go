package thalia

import (
	"strings"
	"testing"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/stretchr/testify/require"
)

func TestBooksFromHTML(t *testing.T) {
	searchHTML := `
	<html>
		<body>
			<ul class="tm-produktliste">
				<li class="tm-produktliste__eintrag artikel element-list-divider" data-productpage="1" data-tm-index="1">
					<a class="element-link-toplevel tm-produkt-link" href="/shop/home/artikeldetails/A1023034631">The Hobbit</a>
					<picture class="element-card-product tm-artikelbild-wrapper">
						<img class="tm-artikelbild-wrapper__artikelbild" src="https://images.thalia.media/03/-/39b7351cb4eb4891a2abe1b171b02776/the-hobbit-gebundene-ausgabe-j-r-r-tolkien-englisch.jpeg" />
					</picture>
					<div class="tm-artikeldetails">
						<p class="element-text-standard tm-artikeldetails__autor autoren-wrapper">J. R. R. Tolkien</p>
						<strong class="element-text-standard-strong tm-artikeldetails__titel">The Hobbit</strong>
						<p class="element-text-standard tm-artikeldetails__formatbezeichnung">
							<span>Buch (Gebundene Ausgabe)</span>
							<span>+ weitere</span>
						</p>
					</div>
				</li>
				<li class="tm-produktliste__eintrag artikel element-list-divider" data-productpage="1" data-tm-index="2">
					<a class="element-link-toplevel tm-produkt-link" href="/shop/home/artikeldetails/A0000000002">Ein Junge namens Weihnacht. Die Serien-Box</a>
					<picture class="element-card-product tm-artikelbild-wrapper">
						<img class="tm-artikelbild-wrapper__artikelbild" src="https://images.thalia.media/03/-/900301430e25462492c753fbccf5bdfe/ein-junge-namens-weihnacht-die-serien-box-mp3-rufus-beck.jpeg" />
					</picture>
					<div class="tm-artikeldetails">
						<p class="element-text-standard tm-artikeldetails__autor autoren-wrapper">Matt Haig</p>
						<strong class="element-text-standard-strong tm-artikeldetails__titel">Ein Junge namens Weihnacht. Die Serien-Box</strong>
						<p class="element-text-standard tm-artikeldetails__formatbezeichnung">
							<span>Hörbuch-Download (MP3)</span>
						</p>
					</div>
				</li>
				<li class="tm-produktliste__eintrag artikel element-list-divider" data-productpage="1" data-tm-index="3">
					<a class="element-link-toplevel tm-produkt-link" href="/shop/home/artikeldetails/A0000000003">SUPERFAN Socken</a>
					<div class="tm-artikeldetails">
						<p class="element-text-standard tm-artikeldetails__autor autoren-wrapper"></p>
						<strong class="element-text-standard-strong tm-artikeldetails__titel">SUPERFAN Socken</strong>
						<p class="element-text-standard tm-artikeldetails__formatbezeichnung">Bunt, 36-41</p>
					</div>
				</li>
			</ul>
		</body>
	</html>`

	searchNode, err := htmlquery.Parse(strings.NewReader(searchHTML))
	require.NoError(t, err)

	books, err := BooksFromHTML(searchNode)
	require.NoError(t, err)
	require.Len(t, books, 2)

	book := books[0]
	require.Equal(t, "A1023034631", book.ID)
	require.Equal(t, "/shop/home/artikeldetails/A1023034631", book.Path)
	require.Equal(t, "The Hobbit", book.Title)
	require.Equal(t, "J. R. R. Tolkien", book.Author)
	require.Equal(t, "Buch (Gebundene Ausgabe)", book.Format)
	require.Equal(t, "https://images.thalia.media/03/-/39b7351cb4eb4891a2abe1b171b02776/the-hobbit-gebundene-ausgabe-j-r-r-tolkien-englisch.jpeg", book.Cover)

	audiobook := books[1]
	require.Equal(t, "A0000000002", audiobook.ID)
	require.Equal(t, "Ein Junge namens Weihnacht. Die Serien-Box", audiobook.Title)
	require.Equal(t, "Matt Haig", audiobook.Author)
	require.Equal(t, "Hörbuch-Download (MP3)", audiobook.Format)
}

func TestBookDetailsFromHTML(t *testing.T) {
	detailHTML := `
	<html>
		<head>
			<script type="application/ld+json">
			{
				"@context": "https://schema.org/",
				"@type": "Book",
				"isbn": "978-0-00-748730-1",
				"name": "The Hobbit",
				"image": ["https://images.thalia.media/-/BF2000-2000/39b7351cb4eb4891a2abe1b171b02776/the-hobbit-gebundene-ausgabe-j-r-r-tolkien-englisch.jpeg"],
				"description": "A special collector's film tie-in hardback edition."
			}
			</script>
		</head>
		<body>
			<a class="element-link-standard autor-name">J. R. R. Tolkien</a>
			<dl-pageview produkttitel="The Hobbit" hersteller="HarperCollins" coverbild="https://images.thalia.media/00/-/39b7351cb4eb4891a2abe1b171b02776/the-hobbit-gebundene-ausgabe-j-r-r-tolkien-englisch.jpeg" form="Gebundene Ausgabe"></dl-pageview>
			<div class="details-default">
				<div class="artikeldetails">
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Einband</h3>
						<p class="element-text-standard value">Gebundene Ausgabe</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Erscheinungsdatum</h3>
						<p class="element-text-standard value">08.11.2012</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Verlag</h3>
						<a class="element-link-standard value">HarperCollins</a>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Sprache</h3>
						<p class="element-text-standard value">Englisch</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">ISBN</h3>
						<p class="element-text-standard value">978-0-00-748730-1</p>
					</section>
				</div>
			</div>
		</body>
	</html>`

	detailNode, err := htmlquery.Parse(strings.NewReader(detailHTML))
	require.NoError(t, err)

	book, err := BookDetailsFromHTML(detailNode)
	require.NoError(t, err)
	require.NotNil(t, book)

	expectedPublishDate := time.Date(2012, time.November, 8, 0, 0, 0, 0, time.UTC)
	require.Equal(t, "The Hobbit", book.Title)
	require.Equal(t, "J. R. R. Tolkien", book.Author)
	require.Equal(t, "Gebundene Ausgabe", book.Format)
	require.Equal(t, "https://images.thalia.media/-/BF2000-2000/39b7351cb4eb4891a2abe1b171b02776/the-hobbit-gebundene-ausgabe-j-r-r-tolkien-englisch.jpeg", book.Cover)
	require.Equal(t, "A special collector's film tie-in hardback edition.", book.Description)
	require.Equal(t, "HarperCollins", book.Publisher)
	require.Equal(t, "Englisch", book.Language)
	require.Equal(t, "978-0-00-748730-1", book.ISBN)
	require.Equal(t, expectedPublishDate, *book.PublishDate)
}

func TestBookDetailsFromHTML_Audiobook(t *testing.T) {
	detailHTML := `
	<html>
		<head>
			<script type="application/ld+json">
			{
				"@context": "https://schema.org/",
				"@type": "Product",
				"gtin13": "9783742432193",
				"name": "Folge 238: Falsche Schuld",
				"image": ["https://images.thalia.media/-/BF2000-2000/d2a4bb547a9246ac9ec97c44e9c0a5e4/folge-238-falsche-schuld-mp3.jpeg"],
				"description": "Eine geteilte Insel, eine spirituelle Sekte und ein Einbruch."
			}
			</script>
		</head>
		<body>
			<div class="autoren-wrapper">
				<a class="element-link-standard autor-name">Andre Minninger</a>
				<a class="element-link-standard autor-name">Ben Nevis</a>
			</div>
			<div class="serien-wrapper">
				<a class="element-link-standard serie-link">Die drei ???</a>
			</div>
			<div class="sprecher-wrapper">
				<a class="element-link-standard sprecher-name">Rufus Beck</a>
			</div>
			<span class="element-text-standard untertitel">Ungekürzte Lesung mit Rufus Beck</span>
			<dl-pageview produkttitel="Folge 238: Falsche Schuld" hersteller="EUROPA/Sony Music Family Entertainment" coverbild="https://images.thalia.media/00/-/d2a4bb547a9246ac9ec97c44e9c0a5e4/folge-238-falsche-schuld-mp3.jpeg"></dl-pageview>
			<div class="details-default">
				<div class="artikeldetails">
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Gesprochen von</h3>
						<a class="element-link-standard value">Rufus Beck</a>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Spieldauer</h3>
						<p class="element-text-standard value">1 Stunde und 21 Minuten</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Erscheinungsdatum</h3>
						<p class="element-text-standard value">20.03.2026</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Hörtyp</h3>
						<p class="element-text-standard value">Hörspiel</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Fassung</h3>
						<p class="element-text-standard value">gekürzt</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Medium</h3>
						<p class="element-text-standard value">MP3</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Anzahl Dateien</h3>
						<p class="element-text-standard value">28</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Verlag</h3>
						<a class="element-link-standard value">EUROPA/Sony Music Family Entertainment</a>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Sprache</h3>
						<p class="element-text-standard value">Deutsch</p>
					</section>
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">EAN</h3>
						<p class="element-text-standard value">9783742432193</p>
					</section>
				</div>
			</div>
		</body>
	</html>`

	detailNode, err := htmlquery.Parse(strings.NewReader(detailHTML))
	require.NoError(t, err)

	book, err := BookDetailsFromHTML(detailNode)
	require.NoError(t, err)
	require.NotNil(t, book)

	expectedPublishDate := time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC)
	expectedDuration := 81 * 60
	require.Equal(t, "Folge 238: Falsche Schuld", book.Title)
	require.Equal(t, "Ungekürzte Lesung mit Rufus Beck", book.Subtitle)
	require.Equal(t, "Andre Minninger, Ben Nevis", book.Author)
	require.Equal(t, "Rufus Beck", book.Narrator)
	require.Equal(t, "https://images.thalia.media/-/BF2000-2000/d2a4bb547a9246ac9ec97c44e9c0a5e4/folge-238-falsche-schuld-mp3.jpeg", book.Cover)
	require.Equal(t, "Eine geteilte Insel, eine spirituelle Sekte und ein Einbruch.", book.Description)
	require.Equal(t, "EUROPA/Sony Music Family Entertainment", book.Publisher)
	require.Equal(t, "Deutsch", book.Language)
	require.Equal(t, "9783742432193", book.ISBN)
	require.Equal(t, expectedPublishDate, *book.PublishDate)
	require.Equal(t, expectedDuration, *book.Duration)
	require.Equal(t, "Die drei ???", book.Series)
	require.Equal(t, "238", book.Sequence)
	require.Equal(t, []string{
		"Hörtyp: Hörspiel",
		"Fassung: gekürzt",
		"Medium: MP3",
		"Anzahl Dateien: 28",
	}, book.Tags)
}

func TestBookDetailsFromHTML_DecodesEntitiesAndPrefersDescriptionTemplate(t *testing.T) {
	detailHTML := `
	<html>
		<head>
			<script type="application/ld+json">
			{
				"@context": "https://schema.org/",
				"@type": "Product",
				"gtin13": "4069829540988",
				"name": "Die Lichtsch&ouml;pferin",
				"image": ["https://images.thalia.media/-/BF2000-2000/e3d3a7fb4eab486caf5d3c3cc7365501/die-lichtschoepferin-mp3-leoni-oeffinger.jpeg"],
				"description": "Sie ist seine Hoffnung und er ihr Verh&auml;ngnis."
			}
			</script>
			<meta property="og:description" content="Kurze Beschreibung mit &ouml;"/>
		</head>
		<body>
			<template data-id="zusatztexte">
				<div class="zusatztexte">
					Sie ist seine Hoffnung und er ihr Verhängnis.<br><br>Amber Montclair gehört keinem der großen Hexenzirkel an.
				</div>
			</template>
		</body>
	</html>`

	detailNode, err := htmlquery.Parse(strings.NewReader(detailHTML))
	require.NoError(t, err)

	book, err := BookDetailsFromHTML(detailNode)
	require.NoError(t, err)
	require.NotNil(t, book)
	require.Equal(t, "Die Lichtschöpferin", book.Title)
	require.Equal(t, "Sie ist seine Hoffnung und er ihr Verhängnis. Amber Montclair gehört keinem der großen Hexenzirkel an.", book.Description)
}

func TestBookDetailsFromHTML_DescriptionFallsBackToMetaDescription(t *testing.T) {
	detailHTML := `
	<html>
		<head>
			<meta property="og:description" content="Kurzbeschreibung mit &ouml; und &auml;"/>
		</head>
		<body></body>
	</html>`

	detailNode, err := htmlquery.Parse(strings.NewReader(detailHTML))
	require.NoError(t, err)

	book, err := BookDetailsFromHTML(detailNode)
	require.NoError(t, err)
	require.NotNil(t, book)
	require.Equal(t, "Kurzbeschreibung mit ö und ä", book.Description)
}

func TestBookDetailsFromHTML_AudiobookNarratorFallback(t *testing.T) {
	detailHTML := `
	<html>
		<body>
			<span class="element-text-standard untertitel">Ungekürzte Lesung mit Rufus Beck (4 CDs)</span>
			<div class="details-default">
				<div class="artikeldetails">
					<section class="artikeldetail">
						<h3 class="element-text-standard-strong detailbezeichnung">Spieldauer</h3>
						<p class="element-text-standard value">4 Std. 38 Min.</p>
					</section>
				</div>
			</div>
		</body>
	</html>`

	detailNode, err := htmlquery.Parse(strings.NewReader(detailHTML))
	require.NoError(t, err)

	book, err := BookDetailsFromHTML(detailNode)
	require.NoError(t, err)
	require.NotNil(t, book)
	require.Equal(t, "Rufus Beck", book.Narrator)
	require.NotNil(t, book.Duration)
	require.Equal(t, 4*3600+38*60, *book.Duration)
}

func TestParseDuration(t *testing.T) {
	duration := parseDuration("4 Stunden und 38 Minuten")
	require.NotNil(t, duration)
	require.Equal(t, 4*3600+38*60, *duration)

	duration = parseDuration("57 Minuten")
	require.NotNil(t, duration)
	require.Equal(t, 57*60, *duration)

	duration = parseDuration("4 Std. 38 Min.")
	require.NotNil(t, duration)
	require.Equal(t, 4*3600+38*60, *duration)
}
