# FBD

Go library for downloading [Football Bet Data](https://www.football-bet-data.com/) exports programmatically.

## Prerequisites

- Go 1.22+
- A valid [Football Bet Data](https://www.football-bet-data.com/) account

## Usage

Install:

```bash
go get github.com/LonecastSystems/fbd-go
```

Import the package:

```go
import fbd "github.com/LonecastSystems/fbd-go"
```

Basic flow (with `os` imported):

```go
client, err := fbd.NewClient()
if err != nil {
    // handle error
}

if err := client.SignIn(username, password); err != nil {
    // handle error
}
defer client.SignOut()

dashboardBytes, err := client.NewDashboardConfigBuilder().
    WithMatchesNoPrediction().
    WithLeagues(map[fbd.Country][]string{fbd.ENGLAND: {"1", "2"}}).
    WithSummerSeasons([]int{2025}).
    WithWinterSeasons([]int{2025}).
    Build().
    ExportToExcel()
if err != nil {
    // handle error
}

if err := os.WriteFile("dashboard.xlsx", dashboardBytes, 0644); err != nil {
    // handle error
}
```

## Notes

- Keep credentials out of source code. Read them from environment variables or your secrets manager.
