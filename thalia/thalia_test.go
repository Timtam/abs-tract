package thalia_test

import (
	"context"
	"testing"

	"github.com/ahobsonsayers/abs-tract/thalia"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestSearchBook(t *testing.T) {
	t.Skip("Fails in CI")

	books, err := thalia.DefaultClient.Search(
		context.Background(),
		"The Hobbit",
		lo.ToPtr("J. R. R. Tolkien"),
	)
	require.NoError(t, err)
	require.NotEmpty(t, books)

	book := books[0]
	require.Equal(t, "The Hobbit", book.Title)
	require.Equal(t, "J. R. R. Tolkien", book.Author)
	require.NotEmpty(t, book.Cover)
	require.NotEmpty(t, book.ISBN)
	require.Equal(t, "HarperCollins", book.Publisher)
}
