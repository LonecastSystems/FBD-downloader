package main

import (
	"log"

	"github.com/LonecastSystems/FBD-downloader/fbd"
)

func main() {
	email := ""
	password := ""
	outputPath := "E1_2025.xlsx"

	client, err := fbd.NewClient()
	if err != nil {
		log.Fatalf("client setup failed: %v", err)
	}

	err = client.SignIn(email, password)
	if err != nil {
		log.Fatalf("sign-in failed: %v", err)
	}
	defer func() {
		if err := client.SignOut(); err != nil {
			log.Printf("sign-out failed: %v", err)
		}
	}()

	if err := client.NewDashboardConfigBuilder().
		WithMatchesNoPrediction().
		WithLeagues(map[fbd.Country][]string{fbd.ENGLAND: {"1", "2"}}).
		WithSummerSeasons([]int{2025}).
		WithWinterSeasons([]int{2025}).
		WithMonths([]fbd.Month{
			fbd.MonthJan, fbd.MonthFeb, fbd.MonthMar, fbd.MonthApr,
			fbd.MonthMay, fbd.MonthJun, fbd.MonthJul, fbd.MonthAug,
			fbd.MonthSep, fbd.MonthOct, fbd.MonthNov, fbd.MonthDec,
		}).
		Build().
		ExportToExcel(outputPath); err != nil {
		log.Fatalf("dashboard export failed: %v", err)
	}

	log.Printf("Dashboard Excel exported to %s", outputPath)
}
