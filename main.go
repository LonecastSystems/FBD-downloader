package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-rod/rod"
)

var Countries = map[string]string{
	"ARGENTINA":    "AR",
	"AUSTRIA":      "AU",
	"BELGIUM":      "B",
	"BOLIVIA":      "B0",
	"BRAZIL":       "BR",
	"BULGARIA":     "BU",
	"CHILE":        "CH",
	"CHINA":        "CN",
	"COLOMBIA":     "CO",
	"CROATIA":      "CR",
	"CZECH":        "CZ",
	"GERMANY":      "D",
	"DENMARK":      "DE",
	"ENGLAND":      "E",
	"FINLAND":      "FI",
	"FRANCE":       "FR",
	"GREECE":       "G",
	"HUNGARY":      "HU",
	"ICELAND":      "IC",
	"IRELAND":      "IR",
	"ITALY":        "IT",
	"JAPAN":        "J",
	"MEXICO":       "MX",
	"NETHERLANDS":  "N",
	"NORWAY":       "NO",
	"PARAGUAY":     "PA",
	"POLAND":       "PL",
	"PORTUGAL":     "PT",
	"ROMANIA":      "RO",
	"RUSSIA":       "RU",
	"SOUTH AFRICA": "SA",
	"SCOTLAND":     "SC",
	"SINGAPORE":    "SI",
	"SLOVAKIA":     "SL",
	"SPAIN":        "SP",
	"SWITZERLAND":  "SU",
	"SLOVENIA":     "SV",
	"SWEDEN":       "SW",
	"TURKEY":       "T",
	"UKRAINE":      "UK",
	"USA":          "US",
	"WALES":        "WA",
}

var Leagues = map[string][]int32{
	"AR": {1},
	"AU": {1, 2},
	"B":  {1, 2},
	"B0": {1},
	"BR": {1},
	"BU": {1},
	"CH": {1},
	"CN": {1},
	"CO": {1},
	"CR": {1},
	"CZ": {1, 2},
	"D":  {1, 2, 3},
	"DE": {1, 2},
	"E":  {0, 1, 2, 3},
	"FI": {1, 2},
	"FR": {1, 2, 3},
	"G":  {1, 2},
	"HU": {1},
	"IC": {1},
	"IR": {1},
	"IT": {1, 2, 3},
	"J":  {1, 2},
	"MX": {1},
	"N":  {1, 2},
	"NO": {1, 2},
	"PA": {1},
	"PL": {1, 2},
	"PT": {1, 2},
	"RO": {1},
	"RU": {1},
	"SA": {1},
	"SC": {0, 1, 2, 3},
	"SI": {1},
	"SL": {1, 2},
	"SP": {1, 2},
	"SU": {1, 2},
	"SV": {1},
	"SW": {1, 2},
	"T":  {1, 2},
	"UK": {1},
	"US": {1},
	"WA": {1},
}

func GetEnvVariables() (email string, password string, yearsInt int, err error) {
	email = os.Getenv("email")
	if email == "" {
		return "", "", 0, errors.New("no email")
	}

	password = os.Getenv("password")
	if password == "" {
		return "", "", 0, errors.New("no password")
	}

	years := os.Getenv("years")
	if years == "" {
		log.Panic("No years threshold set")
		return "", "", 0, errors.New("no years threshold set")
	}

	yearsInt, err = strconv.Atoi(years)
	if err != nil {
		return "", "", 0, errors.New("failed parsing year")
	}

	return email, password, yearsInt, err
}

func main() {
	email, password, yearsInt, err := GetEnvVariables()
	if err != nil {
		log.Panic(err.Error())
	}

	page := rod.New().MustConnect().MustPage("https://www.football-bet-data.com/signin/").MustWaitStable()

	page.MustElement("#ContentPlaceHolder2_unameTextBox").MustInput(email)
	page.MustElement("#ContentPlaceHolder2_pwordTextBox").MustInput(password)
	page.MustElement("#ContentPlaceHolder2_submitButton").MustClick()
	page.MustWaitStable()

	_, err = page.Elements("#ContentPlaceHolder2_Label3")
	if err != nil {
		log.Panic("Invalid credentials")
	}

	//Deselect years
	page.MustNavigate("https://www.football-bet-data.com/dashboard/").MustWaitStable()

	currentYear := time.Now().Year()

	page.MustElement("#ContentPlaceHolder2_seasonSA").MustClick()
	page.MustElement("#ContentPlaceHolder2_summerSA").MustClick()

	for i := 0; i < yearsInt; i++ {
		year := currentYear - i
		page.MustElement(fmt.Sprintf("#ContentPlaceHolder2_%v", year)).MustClick()

		yearShort := year - 2000
		page.MustElement(fmt.Sprintf("#ContentPlaceHolder2_%v", fmt.Sprintf("%v-%v", yearShort, yearShort+1))).MustClick()
	}

	for country, code := range Countries {
		leagueCodes, exists := Leagues[code]
		if !exists {
			continue
		}

		for _, leaguecode := range leagueCodes {
			page.MustElement(fmt.Sprintf("#ContentPlaceHolder2_%v%v", code, leaguecode)).MustClick()
		}
		fmt.Print(country)
		page.MustElement("#ContentPlaceHolder2_ButtonEX2").MustClick().MustWaitInteractable()

		//Safer than MustDoubleClick!
		page.MustElement("#ContentPlaceHolder2_leagueSA").MustClick()
		page.MustElement("#ContentPlaceHolder2_leagueSA").MustClick()
	}

	os.Exit(1)
}