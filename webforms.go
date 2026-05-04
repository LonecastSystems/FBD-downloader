package fbd

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"regexp"
)

func (c *Client) getWebFormsFields(pageURL string) (url.Values, error) {
	resp, err := c.httpClient.Get(pageURL)
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
