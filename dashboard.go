package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
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

var Leagues = map[Country][]int32{
	ARGENTINA: {1},
	AUSTRIA:   {1, 2},
	BELGIUM:   {1, 2},
	//BOLIVIA: {1},
	BRAZIL:       {1},
	BULGARIA:     {1},
	CHILE:        {1},
	CHINA:        {1},
	COLOMBIA:     {1},
	CROATIA:      {1},
	CZECH:        {1, 2},
	GERMANY:      {1, 2, 3},
	DENMARK:      {1, 2},
	ENGLAND:      {0, 1, 2, 3},
	FINLAND:      {1, 2},
	FRANCE:       {1, 2, 3},
	GREECE:       {1, 2},
	HUNGARY:      {1},
	ICELAND:      {1},
	IRELAND:      {1},
	ITALY:        {1, 2, 3},
	JAPAN:        {1, 2},
	MEXICO:       {1},
	NETHERLANDS:  {1, 2},
	NORWAY:       {1, 2},
	PARAGUAY:     {1},
	POLAND:       {1, 2},
	PORTUGAL:     {1, 2},
	ROMANIA:      {1},
	RUSSIA:       {1},
	SOUTH_AFRICA: {1},
	SCOTLAND:     {0, 1, 2, 3},
	SINGAPORE:    {1},
	SLOVAKIA:     {1, 2},
	SPAIN:        {1, 2},
	SWITZERLAND:  {1, 2},
	SLOVENIA:     {1},
	SWEDEN:       {1, 2},
	TURKEY:       {1, 2},
	UKRAINE:      {1},
	USA:          {1},
	WALES:        {1},
}

type Dashboard struct {
	Email    string
	Password string
	Path     string
	Config   DashboardConfig
}

type DashboardConfig struct {
	Leagues       map[Country][]int32
	SummerSeasons []int
	WinterSeasons []int
}

func NewDashboard(config DashboardConfig) *Dashboard {
	return &Dashboard{Config: config}
}

func (dashboard Dashboard) Download() {
	config := dashboard.Config

	email, password := dashboard.Email, dashboard.Password
	if email == "" || password == "" {
		log.Fatal("Invalid credentials!")
	}

	summerSeasons, winterSeasons := config.SummerSeasons, config.WinterSeasons
	if len(summerSeasons) == 0 && len(winterSeasons) == 0 {
		log.Fatal("No seasons selected!")
	}

	path := dashboard.Path
	if path == "" {
		log.Fatal("Invalid path!")
	}

	page := rod.New().MustConnect().MustPage("https://www.football-bet-data.com/signin/").MustWaitStable()

	//Login
	page.MustElement("#ContentPlaceHolder2_unameTextBox").MustInput(email)
	page.MustElement("#ContentPlaceHolder2_pwordTextBox").MustInput(password)
	page.MustElement("#ContentPlaceHolder2_submitButton").MustClick()
	page.MustWaitStable()

	//Login fail?
	if errors, _ := page.Elements("#ContentPlaceHolder2_Label3"); !errors.Empty() {
		log.Panicf("Invalid credentials: %v", errors.First().MustText())
	}

	//Go to dashboard
	page.MustNavigate("https://www.football-bet-data.com/dashboard/").MustWaitStable()

	currentYear := time.Now().Year()

	//Deselect summer years
	page.MustElement("#ContentPlaceHolder2_summerSA").MustClick()

	//Select summer seasons
	for _, year := range summerSeasons {
		if year > currentYear {
			log.Printf("%v is an invalid summer season", year)
			continue
		}

		page.MustElement(fmt.Sprintf("#ContentPlaceHolder2_%v", year)).MustClick()
	}

	//Deselect winter years
	page.MustElement("#ContentPlaceHolder2_seasonSA").MustClick()

	//Select winter seasons
	for _, year := range winterSeasons {
		if year > currentYear {
			log.Printf("%v is an invalid winter season", year)
			continue
		}

		yearShort := year - 2000
		if winterSeason, err := page.Timeout(10 * time.Second).Element(fmt.Sprintf("#ContentPlaceHolder2_%v", fmt.Sprintf("%v-%v", yearShort, yearShort+1))); winterSeason != nil && err == nil {
			winterSeason.MustClick()
		}
	}

	//Uncheck all leagues
	page.MustElement("#ContentPlaceHolder2_leagueSA").MustClick()

	//Select leagues
	leagues := config.Leagues
	if len(leagues) == 0 {
		leagues = Leagues
	}

	//Select all leagues
	for code, leagueCodes := range leagues {
		for _, leagueCode := range leagueCodes {
			page.MustElement(fmt.Sprintf("#ContentPlaceHolder2_%v%v", code, leagueCode)).MustClick()
		}

		//Download the excel sheet
		clicked := page.MustElement("#ContentPlaceHolder2_ButtonEX2").MustClick()
		bytes := clicked.MustFrame().Browser().MustWaitDownload()()
		os.WriteFile(fmt.Sprintf("%v\\FBDResults_%v.xlsx", path, code), bytes, 0644)

		//Uncheck all leagues
		//Safer than MustDoubleClick!
		page.MustElement("#ContentPlaceHolder2_leagueSA").MustClick()
		page.MustElement("#ContentPlaceHolder2_leagueSA").MustClick()
	}
}
