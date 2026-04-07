package server

import (
	"testing"
	"time"

	"github.com/ahobsonsayers/abs-tract/thalia"
	"github.com/stretchr/testify/require"
)

func TestThaliaBookToBookMetadata_Audiobook(t *testing.T) {
	publishDate := time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC)
	duration := 81
	abridged := true

	book := thalia.Book{
		Title:       "Folge 238: Falsche Schuld",
		Subtitle:    "Ungekürzte Lesung mit Rufus Beck",
		Author:      "Andre Minninger, Ben Nevis",
		Narrator:    "Rufus Beck",
		Abridged:    &abridged,
		Format:      "Hörbuch-Download (MP3)",
		Cover:       "https://images.thalia.media/example.jpeg",
		Description: "Eine geteilte Insel, eine spirituelle Sekte und ein Einbruch.",
		Publisher:   "EUROPA/Sony Music Family Entertainment",
		Language:    "Deutsch",
		ISBN:        "9783742432193",
		Duration:    &duration,
		Series:      "Die drei ???",
		Sequence:    "238",
		PublishDate: &publishDate,
	}

	metadata := thaliaBookToBookMetadata(book)
	require.Equal(t, "Folge 238: Falsche Schuld", metadata.Title)
	require.NotNil(t, metadata.Subtitle)
	require.Equal(t, "Ungekürzte Lesung mit Rufus Beck", *metadata.Subtitle)
	require.NotNil(t, metadata.Author)
	require.Equal(t, "Andre Minninger, Ben Nevis", *metadata.Author)
	require.NotNil(t, metadata.Narrator)
	require.Equal(t, "Rufus Beck", *metadata.Narrator)
	require.NotNil(t, metadata.Abridged)
	require.True(t, *metadata.Abridged)
	require.NotNil(t, metadata.Duration)
	require.Equal(t, duration, *metadata.Duration)
	require.NotNil(t, metadata.Series)
	require.Len(t, *metadata.Series, 1)
	require.Equal(t, "Die drei ???", (*metadata.Series)[0].Series)
	require.NotNil(t, (*metadata.Series)[0].Sequence)
	require.Equal(t, "238", *(*metadata.Series)[0].Sequence)
	require.Nil(t, metadata.Tags)
	require.NotNil(t, metadata.PublishedYear)
	require.Equal(t, "2026", *metadata.PublishedYear)
}

func TestThaliaBookToBookMetadata_PreservesEANAsIsbn(t *testing.T) {
	book := thalia.Book{
		Title: "Die Lichtschöpferin",
		ISBN:  "4069829540988",
	}

	metadata := thaliaBookToBookMetadata(book)
	require.NotNil(t, metadata.Isbn)
	require.Equal(t, "4069829540988", *metadata.Isbn)
}
