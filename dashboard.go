package fbd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

type Country string

const (
	ARGENTINA    Country = "AR"
	AUSTRIA      Country = "AU"
	BELGIUM      Country = "B"
	BOLIVIA      Country = "BO"
	BRAZIL       Country = "BR"
	BULGARIA     Country = "BU"
	CHILE        Country = "CH"
	CHINA        Country = "CN"
	COLOMBIA     Country = "CO"
	CROATIA      Country = "CR"
	CZECH        Country = "CZ"
	GERMANY      Country = "D"
	DENMARK      Country = "DE"
	ENGLAND      Country = "E"
	FINLAND      Country = "FI"
	FRANCE       Country = "FR"
	GREECE       Country = "G"
	HUNGARY      Country = "HU"
	ICELAND      Country = "IC"
	IRELAND      Country = "IR"
	ITALY        Country = "IT"
	JAPAN        Country = "J"
	MEXICO       Country = "MX"
	NETHERLANDS  Country = "N"
	NORWAY       Country = "NO"
	PARAGUAY     Country = "PA"
	POLAND       Country = "PL"
	PORTUGAL     Country = "PT"
	ROMANIA      Country = "RO"
	RUSSIA       Country = "RU"
	SOUTH_AFRICA Country = "SA"
	SCOTLAND     Country = "SC"
	SINGAPORE    Country = "SI"
	SLOVAKIA     Country = "SL"
	SPAIN        Country = "SP"
	SWITZERLAND  Country = "SU"
	SLOVENIA     Country = "SV"
	SWEDEN       Country = "SW"
	TURKEY       Country = "T"
	UKRAINE      Country = "UK"
	USA          Country = "US"
	WALES        Country = "WA"
)

var ValidLeagues = map[Country][]string{
	ARGENTINA:    {"1"},
	AUSTRIA:      {"1", "2"},
	BELGIUM:      {"1", "2"},
	BRAZIL:       {"1"},
	BULGARIA:     {"1"},
	CHILE:        {"1"},
	CHINA:        {"1"},
	COLOMBIA:     {"1"},
	CROATIA:      {"1"},
	CZECH:        {"1", "2"},
	GERMANY:      {"1", "2", "3"},
	DENMARK:      {"1", "2"},
	ENGLAND:      {"0", "1", "2", "3", "C", "CN", "CS"},
	FINLAND:      {"1", "2"},
	FRANCE:       {"1", "2", "3"},
	GREECE:       {"1", "2"},
	HUNGARY:      {"1"},
	ICELAND:      {"1"},
	IRELAND:      {"1"},
	ITALY:        {"1", "2", "3"},
	JAPAN:        {"1", "2"},
	MEXICO:       {"1"},
	NETHERLANDS:  {"1", "2"},
	NORWAY:       {"1", "2"},
	PARAGUAY:     {"1"},
	POLAND:       {"1", "2"},
	PORTUGAL:     {"1", "2"},
	ROMANIA:      {"1"},
	RUSSIA:       {"1"},
	SOUTH_AFRICA: {"1"},
	SCOTLAND:     {"0", "1", "2", "3"},
	SINGAPORE:    {"1"},
	SLOVAKIA:     {"1", "2"},
	SPAIN:        {"1", "2"},
	SWITZERLAND:  {"1", "2"},
	SLOVENIA:     {"1"},
	SWEDEN:       {"1", "2"},
	TURKEY:       {"1", "2"},
	UKRAINE:      {"1"},
	USA:          {"1"},
	WALES:        {"1"},
}

type OddsType int

const (
	OddsTypeAverage OddsType = 1
	OddsTypeBetfair OddsType = 2
)

type Range struct {
	Min int
	Max int
}

type Month string

const (
	MonthJan Month = "fidJan"
	MonthFeb Month = "fidFeb"
	MonthMar Month = "fidMar"
	MonthApr Month = "fidApr"
	MonthMay Month = "fidMay"
	MonthJun Month = "fidJun"
	MonthJul Month = "fidJul"
	MonthAug Month = "fidAug"
	MonthSep Month = "fidSep"
	MonthOct Month = "fidOct"
	MonthNov Month = "fidNov"
	MonthDec Month = "fidDec"
)

var AllMonths = []Month{
	MonthJan, MonthFeb, MonthMar, MonthApr, MonthMay, MonthJun,
	MonthJul, MonthAug, MonthSep, MonthOct, MonthNov, MonthDec,
}

var monthValues = map[Month]string{
	MonthJan: "1",
	MonthFeb: "2",
	MonthMar: "3",
	MonthApr: "4",
	MonthMay: "5",
	MonthJun: "6",
	MonthJul: "7",
	MonthAug: "8",
	MonthSep: "9",
	MonthOct: "10",
	MonthNov: "11",
	MonthDec: "12",
}

type ScoreLine string

const (
	ScoreLine00 ScoreLine = "0 - 0"
	ScoreLine11 ScoreLine = "1 - 1"
	ScoreLine22 ScoreLine = "2 - 2"
	ScoreLine33 ScoreLine = "3 - 3"
	ScoreLine44 ScoreLine = "4 - 4"
	ScoreLine55 ScoreLine = "5 - 5"
	ScoreLine10 ScoreLine = "1 - 0"
	ScoreLine20 ScoreLine = "2 - 0"
	ScoreLine21 ScoreLine = "2 - 1"
	ScoreLine30 ScoreLine = "3 - 0"
	ScoreLine31 ScoreLine = "3 - 1"
	ScoreLine40 ScoreLine = "4 - 0"
	ScoreLine32 ScoreLine = "3 - 2"
	ScoreLine41 ScoreLine = "4 - 1"
	ScoreLine50 ScoreLine = "5 - 0"
	ScoreLine42 ScoreLine = "4 - 2"
	ScoreLine51 ScoreLine = "5 - 1"
	ScoreLine43 ScoreLine = "4 - 3"
	ScoreLine52 ScoreLine = "5 - 2"
	ScoreLine53 ScoreLine = "5 - 3"
	ScoreLine54 ScoreLine = "5 - 4"
	ScoreLine01 ScoreLine = "0 - 1"
	ScoreLine02 ScoreLine = "0 - 2"
	ScoreLine12 ScoreLine = "1 - 2"
	ScoreLine03 ScoreLine = "0 - 3"
	ScoreLine13 ScoreLine = "1 - 3"
	ScoreLine04 ScoreLine = "0 - 4"
	ScoreLine23 ScoreLine = "2 - 3"
	ScoreLine14 ScoreLine = "1 - 4"
	ScoreLine05 ScoreLine = "0 - 5"
	ScoreLine24 ScoreLine = "2 - 4"
	ScoreLine15 ScoreLine = "1 - 5"
	ScoreLine34 ScoreLine = "3 - 4"
	ScoreLine25 ScoreLine = "2 - 5"
	ScoreLine35 ScoreLine = "3 - 5"
	ScoreLine45 ScoreLine = "4 - 5"
)

var AllScoreLines = []ScoreLine{
	ScoreLine00, ScoreLine11, ScoreLine22, ScoreLine33, ScoreLine44, ScoreLine55,
	ScoreLine10, ScoreLine20, ScoreLine21, ScoreLine30, ScoreLine31, ScoreLine40,
	ScoreLine32, ScoreLine41, ScoreLine50, ScoreLine42, ScoreLine51, ScoreLine43,
	ScoreLine52, ScoreLine53, ScoreLine54,
	ScoreLine01, ScoreLine02, ScoreLine12, ScoreLine03, ScoreLine13, ScoreLine04,
	ScoreLine23, ScoreLine14, ScoreLine05, ScoreLine24, ScoreLine15, ScoreLine34,
	ScoreLine25, ScoreLine35, ScoreLine45,
}

func (s ScoreLine) FieldKey() string {
	return "s" + strings.ReplaceAll(strings.ReplaceAll(string(s), " ", ""), "-", "")
}

type DashboardRange string

const (
	RangeHomeWinOdds           DashboardRange = "home_win_odds"
	RangeDrawOdds              DashboardRange = "draw_odds"
	RangeAwayWinOdds           DashboardRange = "away_win_odds"
	RangeGGOdds                DashboardRange = "gg_odds"
	RangeUnder25Odds           DashboardRange = "under_2_5_odds"
	RangeOver25Odds            DashboardRange = "over_2_5_odds"
	RangeAsianHomeOdds         DashboardRange = "asian_home_odds"
	RangeAsianDrawOdds         DashboardRange = "asian_draw_odds"
	RangeAsianAwayOdds         DashboardRange = "asian_away_odds"
	RangeDoubleChance1X        DashboardRange = "double_chance_1x_odds"
	RangeDoubleChance12        DashboardRange = "double_chance_12_odds"
	RangeDoubleChanceX2        DashboardRange = "double_chance_x2_odds"
	RangeDNBHomeOdds           DashboardRange = "draw_no_bet_home_odds"
	RangeDNBAwayOdds           DashboardRange = "draw_no_bet_away_odds"
	RangeOver05Odds            DashboardRange = "over_0_5_odds"
	RangeUnder05Odds           DashboardRange = "under_0_5_odds"
	RangeOver15Odds            DashboardRange = "over_1_5_odds"
	RangeUnder15Odds           DashboardRange = "under_1_5_odds"
	RangeOver35Odds            DashboardRange = "over_3_5_odds"
	RangeUnder35Odds           DashboardRange = "under_3_5_odds"
	RangeOver45Odds            DashboardRange = "over_4_5_odds"
	RangeUnder45Odds           DashboardRange = "under_4_5_odds"
	RangeBothTeamsScoreYesOdds DashboardRange = "bts_yes_odds"
	RangeBothTeamsScoreNoOdds  DashboardRange = "bts_no_odds"
)

var rangeFieldPrefix = map[DashboardRange]string{
	RangeHomeWinOdds:           "userPho",
	RangeDrawOdds:              "userPdo",
	RangeAwayWinOdds:           "userPao",
	RangeGGOdds:                "userPggo",
	RangeUnder25Odds:           "userPu25o",
	RangeOver25Odds:            "userPo25o",
	RangeAsianHomeOdds:         "userAho",
	RangeAsianDrawOdds:         "userAdo",
	RangeAsianAwayOdds:         "userAao",
	RangeDoubleChance1X:        "userdc1xo",
	RangeDoubleChance12:        "userdc12o",
	RangeDoubleChanceX2:        "userdcx2o",
	RangeDNBHomeOdds:           "userdnbho",
	RangeDNBAwayOdds:           "userdnbao",
	RangeOver05Odds:            "userov05o",
	RangeUnder05Odds:           "userun05o",
	RangeOver15Odds:            "userov15o",
	RangeUnder15Odds:           "userun15o",
	RangeOver35Odds:            "userov35o",
	RangeUnder35Odds:           "userun35o",
	RangeOver45Odds:            "userov45o",
	RangeUnder45Odds:           "userun45o",
	RangeBothTeamsScoreYesOdds: "userbtsyo",
	RangeBothTeamsScoreNoOdds:  "userbtsno",
}

type dashboardConfig struct {
	url    string
	client *Client
	fields url.Values
}

type dashboardConfigBuilder struct {
	config dashboardConfig
}

func (c *Client) NewDashboardConfigBuilder() *dashboardConfigBuilder {
	cfg := dashboardConfig{
		url:    "https://www.football-bet-data.com/dashboard/",
		client: c,
		fields: url.Values{},
	}

	cfg.fields.Set("ctl00$ContentPlaceHolder2$ddlsystem", "8827")
	cfg.fields.Set("ctl00$ContentPlaceHolder2$rDate", "all")

	builder := dashboardConfigBuilder{
		config: cfg,
	}
	builder.WithRanges(map[DashboardRange]Range{})
	builder.WithOddsType(OddsTypeAverage)
	builder.WithMonths(AllMonths)

	return &builder
}

func (b *dashboardConfigBuilder) WithLeagues(leagues map[Country][]string) *dashboardConfigBuilder {
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

func (b *dashboardConfigBuilder) WithOddsType(oddsType OddsType) *dashboardConfigBuilder {
	b.config.fields.Set("ctl00$ContentPlaceHolder2$otype", strconv.Itoa(int(oddsType)))
	return b
}

func (b *dashboardConfigBuilder) WithRanges(ranges map[DashboardRange]Range) *dashboardConfigBuilder {
	for rangeKey, prefix := range rangeFieldPrefix {
		r, ok := ranges[rangeKey]
		if !ok {
			r = Range{Min: 0, Max: 100}
		}

		placeHolder := fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", prefix)
		b.config.fields.Set(fmt.Sprintf("%s1", placeHolder), strconv.Itoa(r.Min))
		b.config.fields.Set(fmt.Sprintf("%s2", placeHolder), strconv.Itoa(r.Max))
	}
	return b
}

func (b *dashboardConfigBuilder) WithSummerSeasons(seasons []int) *dashboardConfigBuilder {
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
func (b *dashboardConfigBuilder) WithWinterSeasons(seasons []int) *dashboardConfigBuilder {
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

func (b *dashboardConfigBuilder) WithMonths(months []Month) *dashboardConfigBuilder {
	for _, month := range AllMonths {
		b.config.fields.Del("ctl00$ContentPlaceHolder2$" + string(month))
	}

	if len(months) == 0 {
		months = AllMonths
	}

	for _, month := range months {
		b.config.fields.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", month), monthValues[month])
	}
	return b
}

func (b *dashboardConfigBuilder) WithScorePredictions(scores []ScoreLine) *dashboardConfigBuilder {
	for _, score := range AllScoreLines {
		b.config.fields.Del("ctl00$ContentPlaceHolder2$" + score.FieldKey())
	}

	for _, score := range scores {
		b.config.fields.Set(fmt.Sprintf("ctl00$ContentPlaceHolder2$%s", score.FieldKey()), string(score))
	}

	if len(scores) > 0 {
		b.config.fields.Set("ctl00$ContentPlaceHolder2$sp", "on")
	}

	return b
}

func (b *dashboardConfigBuilder) WithMatchesNoPrediction() *dashboardConfigBuilder {
	b.config.fields.Set("ctl00$ContentPlaceHolder2$np", "on")
	return b
}

func (b *dashboardConfigBuilder) Build() *dashboardConfig {
	return &b.config
}

func (c *dashboardConfig) ExportToExcel(ctx context.Context) ([]byte, error) {
	if err := c.client.ensureReady(); err != nil {
		return nil, err
	}

	slog.Info("Exporting dashboard")

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
		return nil, fmt.Errorf("post dashboard form: %w", err)
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

func isDigits(value string) bool {
	for _, char := range value {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
