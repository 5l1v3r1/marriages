package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/unixpickle/essentials"
)

const StartYear = 2008

func main() {
	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"app1", "app2", "date", "id"})
	writer.Flush()
	date := time.Now()
	for date.Year() >= StartYear {
		dateStr := fmt.Sprintf("%02d/%02d/%04d", date.Month(), date.Day(), date.Year())
		results, err := MarriagesAtDate(dateStr)
		essentials.Must(err)
		for _, result := range results {
			writer.Write([]string{result.Applicant1, result.Applicant2, result.Date,
				result.LicenseID})
		}
		writer.Flush()
		date = date.AddDate(0, 0, -1)
	}
}
