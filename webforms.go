package fbd

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (c *Client) postForm(ctx context.Context, url string, data url.Values) (*http.Response, error) {
	postReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create post request: %w", err)
	}
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.httpClient.Do(postReq)
}

func (c *Client) getWebFormsFields(ctx context.Context, pageURL string) (url.Values, error) {
	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create get request: %w", err)
	}
	resp, err := c.httpClient.Do(getReq)
	if err != nil {
		return nil, fmt.Errorf("get page %s: %w", pageURL, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read page %s: %w", pageURL, err)
	}
	body := string(bodyBytes)

	fields := url.Values{
		"__EVENTTARGET":     {""},
		"__EVENTARGUMENT":   {""},
		"__SCROLLPOSITIONX": {"0"},
		"__SCROLLPOSITIONY": {"0"},
	}

	for _, name := range []string{"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"} {
		pattern := fmt.Sprintf(`(?is)<input[^>]*name=["']%s["'][^>]*value=["'](.*?)["']`, regexp.QuoteMeta(name))
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) < 2 {
			continue
		}

		if value := html.UnescapeString(matches[1]); value != "" {
			fields.Set(name, value)
		}
	}

	return fields, nil
}
