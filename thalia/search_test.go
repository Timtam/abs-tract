package thalia

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeBookPrefersMoreDescriptiveFormat(t *testing.T) {
	book := mergeBook(
		Book{Format: "Hörbuch-Download (MP3)"},
		Book{Format: "MP3"},
	)

	require.Equal(t, "Hörbuch-Download (MP3)", book.Format)
}
