package thalia_test

import (
	"testing"

	"github.com/ahobsonsayers/abs-tract/thalia"
	"github.com/stretchr/testify/require"
)

func TestNewClient_DefaultURL(t *testing.T) {
	client, err := thalia.NewClient(nil)
	require.NoError(t, err)
	require.Equal(t, thalia.DefaultThaliaURL, client.URL().String())
}

func TestNewClient_CustomURL(t *testing.T) {
	thaliaURL := "https://example.com/shop/"

	client, err := thalia.NewClient(&thaliaURL)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/shop", client.URL().String())
}
