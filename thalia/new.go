package thalia

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/imroc/req/v3"
)

// NewClient creates a new thalia client.
// If thalia url is nil or unset, the default thalia url will be used.
// Will return an error if the thalia url is invalid.
func NewClient(thaliaURL *string) (*Client, error) {
	thaliaURLStruct := defaultThaliaURL
	if thaliaURL != nil && *thaliaURL != "" {
		parsedThaliaURL, err := url.Parse(strings.Trim(*thaliaURL, "/"))
		if err != nil {
			return nil, fmt.Errorf("invalid thalia url: %w", err)
		}
		thaliaURLStruct = parsedThaliaURL
	}

	return &Client{
		client:    req.C().ImpersonateChrome(),
		thaliaURL: utils.CloneURL(thaliaURLStruct),
	}, nil
}
