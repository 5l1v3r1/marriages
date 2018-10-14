package main

import (
	"fmt"
	"time"

	"github.com/unixpickle/essentials"
)

const NumYears = 10

func main() {
	fmt.Println("app1,app2,date,id")
	date := time.Now()
	for i := 0; i < 365*NumYears; i++ {
		dateStr := fmt.Sprintf("%02d/%02d/%04d", date.Month(), date.Day(), date.Year())
		results, err := MarriagesAtDate(dateStr)
		essentials.Must(err)
		for _, result := range results {
			fmt.Printf("%s,%s,%s,%s\n", result.Applicant1, result.Applicant2, result.Date,
				result.LicenseID)
		}
		date = date.AddDate(0, 0, -1)
	}
}
