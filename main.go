package main

import "github.com/LonecastSystems/FBD-downloader/fbd"

func main() {
	seasons := []int{2021, 2022, 2023, 2024, 2025}

	dashboard := fbd.NewDashboard(fbd.DashboardConfig{
		SummerSeasons: seasons,
		WinterSeasons: seasons,
	})

	dashboard.Email = ""
	dashboard.Password = ""
	dashboard.Path = ""

	dashboard.Download()
}
