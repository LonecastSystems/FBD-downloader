package main

func main() {
	seasons := []int{2021, 2022, 2023, 2024, 2025}

	config := DashboardConfig{
		Email:         "",
		Password:      "",
		Path:          "",
		SummerSeasons: seasons,
		WinterSeasons: seasons,
	}

	Download(config)
}
