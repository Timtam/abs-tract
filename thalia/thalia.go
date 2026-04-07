package thalia

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ahobsonsayers/abs-tract/utils"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"golang.org/x/net/html"
)

const DefaultThaliaURL = "https://www.thalia.de"

var (
	defaultThaliaURL = lo.Must(url.Parse(DefaultThaliaURL))

	DefaultClient = &Client{
		client:    req.C().ImpersonateChrome(),
		thaliaURL: utils.CloneURL(defaultThaliaURL),
	}
)

type Client struct {
	client    *req.Client
	thaliaURL *url.URL
}

// URL returns a clone of the thalia url used by the client.
func (c *Client) URL() *url.URL { return utils.CloneURL(c.thaliaURL) }

func (c *Client) get(
	ctx context.Context,
	path string,
	parameters map[string]string,
) (*html.Node, error) {
	queryParams := url.Values{}
	for key, value := range parameters {
		queryParams.Add(key, value)
	}

	requestURL := c.URL()
	if path != "" {
		requestURL = requestURL.JoinPath(path)
	}
	requestURL.RawQuery = queryParams.Encode()

	response, err := c.client.R().SetContext(ctx).Get(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	if !response.IsSuccessState() {
		return nil, fmt.Errorf("%s: %s", response.GetStatus(), response.String())
	}

	htmlResponse, err := html.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return htmlResponse, nil
}
