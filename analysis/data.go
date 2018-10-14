package main

import (
	"encoding/csv"
	"errors"
	"os"
	"time"

	"github.com/unixpickle/essentials"
)

type Marriage struct {
	Names [2]string
	Date  time.Time
	ID    string
}

func ReadData(csvPath string) (m []*Marriage, err error) {
	defer essentials.AddCtxTo("ReadData", &err)
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	contents, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	} else if len(contents) < 2 {
		return nil, errors.New("empty data file")
	} else if len(contents[0]) != 4 {
		return nil, errors.New("incorrect number of columns")
	}
	for _, entry := range contents[1:] {
		date, err := time.Parse("01/02/2006", entry[2])
		if err != nil {
			return nil, err
		}
		m = append(m, &Marriage{
			Names: [2]string{entry[0], entry[1]},
			Date:  date,
			ID:    entry[2],
		})
	}
	return
}
