package main

func main() {
	seasons := []int{2021, 2022, 2023, 2024, 2025}

	dashboard := NewDashboard(DashboardConfig{
		SummerSeasons: seasons,
		WinterSeasons: seasons,
	})

	dashboard.Email = ""
	dashboard.Password = ""
	dashboard.Path = ""

	dashboard.Download()
}
