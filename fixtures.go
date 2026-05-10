package fbd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
)

type fixturesConfig struct {
	url    string
	client *Client
	fields url.Values
}

type fixturesConfigBuilder struct {
	config fixturesConfig
}

func (c *Client) NewFixturesConfigBuilder() *fixturesConfigBuilder {
	return &fixturesConfigBuilder{
		config: fixturesConfig{
			url:    "https://www.football-bet-data.com/history/",
			client: c,
			fields: url.Values{},
		},
	}
}

func (b *fixturesConfigBuilder) WithLeagues(leagues map[Country][]string) *fixturesConfigBuilder {
	for country, leagueCodes := range ValidLeagues {
		for _, leagueCode := range leagueCodes {
			key := string(country) + leagueCode
			b.config.fields.Del("ctl00$ContentPlaceHolder2$" + key)
		}
	}

	if len(leagues) == 0 {
		leagues = ValidLeagues
	}

	for country, leagueCodes := range leagues {
		for _, leagueCode := range leagueCodes {
			key := string(country) + leagueCode
			b.config.fields.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", key), key)
		}
	}

	return b
}

func (b *fixturesConfigBuilder) WithSummerSeasons(seasons []int) *fixturesConfigBuilder {
	const prefix = "ctl00$ContentPlaceHolder2$"
	for key := range b.config.fields {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		name := strings.TrimPrefix(key, prefix)
		if len(name) == 4 && isDigits(name) {
			b.config.fields.Del(key)
		}
	}

	for _, year := range seasons {
		if year <= 0 {
			continue
		}
		key := strconv.Itoa(year)
		b.config.fields.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", key), key)
	}
	return b
}

// WithWinterSeasons expects season start years, e.g. 2024 for 2024-2025 season
func (b *fixturesConfigBuilder) WithWinterSeasons(seasons []int) *fixturesConfigBuilder {
	const prefix = "ctl00$ContentPlaceHolder2$"
	for key := range b.config.fields {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		name := strings.TrimPrefix(key, prefix)
		if len(name) == 5 && name[2] == '-' && isDigits(name[:2]) && isDigits(name[3:]) {
			b.config.fields.Del(key)
		}
	}

	for _, startYear := range seasons {
		if startYear <= 0 {
			continue
		}
		short := startYear % 100
		next := (startYear + 1) % 100
		key := fmt.Sprintf("%02d-%02d", short, next)
		b.config.fields.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", key), key)
	}
	return b
}

func (b *fixturesConfigBuilder) Build() *fixturesConfig {
	return &b.config
}

func (c *fixturesConfig) ExportToExcel(ctx context.Context) ([]byte, error) {
	if err := c.client.ensureReady(); err != nil {
		return nil, err
	}

	slog.Info("Exporting fixtures")

	fields, err := c.client.getWebFormsFields(ctx, c.url)
	if err != nil {
		return nil, err
	}

	fields.Set("ctl00$ContentPlaceHolder2$ButtonEX2", "Export to Excel")

	for key, values := range c.fields {
		fields.Del(key)
		for _, value := range values {
			fields.Add(key, value)
		}
	}

	resp, err := c.client.postForm(ctx, c.url, fields)
	if err != nil {
		return nil, fmt.Errorf("post fixtures form: %w", err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	downloadSucceeded := false
	for _, cookie := range cookies {
		if cookie.Name == "downloadStarted" && cookie.Value == "1" {
			downloadSucceeded = true
			break
		}
	}

	if !downloadSucceeded {
		return nil, fmt.Errorf("download did not succeed: missing or invalid downloadStarted cookie. Cookies: %v. Response body length: %d", cookies, resp.ContentLength)
	}

	return io.ReadAll(resp.Body)
}
