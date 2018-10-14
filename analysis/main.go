package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/unixpickle/essentials"
)

const BarWidth = 20

func main() {
	var dataFile string
	flag.StringVar(&dataFile, "data", "../data/phila_10172008_10142018.csv", "marriage data")
	flag.Parse()

	marriages, err := ReadData(dataFile)
	essentials.Must(err)

	yearAnalysis(marriages)
	monthAnalysis(marriages)
	weekdayAnalysis(marriages)
}

func yearAnalysis(m []*Marriage) {
	printHeading("Yearly marriage tally")
	histogram(m, func(m *Marriage) int {
		return m.Date.Year()
	}, func(year int) string {
		return strconv.Itoa(year)
	})
}

func monthAnalysis(m []*Marriage) {
	printHeading("Monthly marriage tally")
	histogram(m, func(m *Marriage) int {
		return int(m.Date.Month())
	}, func(month int) string {
		return time.Month(month).String()
	})
}

func weekdayAnalysis(m []*Marriage) {
	printHeading("Week-day marriage tally")
	histogram(m, func(m *Marriage) int {
		return int(m.Date.Weekday())
	}, func(day int) string {
		return time.Weekday(day).String()
	})
}

func histogram(m []*Marriage, k func(*Marriage) int, s func(int) string) {
	counts := map[int]int{}
	keys := []int{}
	for _, marriage := range m {
		key := k(marriage)
		if counts[key] == 0 {
			keys = append(keys, key)
		}
		counts[key]++
	}
	sort.Ints(keys)

	keyStrs := []string{}
	countStrs := []string{}
	maxCount := 0
	for _, key := range keys {
		keyStrs = append(keyStrs, s(key))
		countStrs = append(countStrs, strconv.Itoa(counts[key]))
		maxCount = essentials.MaxInt(maxCount, counts[key])
	}

	padColumn(keyStrs, 0)
	padColumn(countStrs, 1)

	for i, key := range keys {
		frac := float64(counts[key]) / float64(maxCount)
		bar := ""
		for i := 0; i < int(frac*BarWidth); i++ {
			bar += "+"
		}
		fmt.Printf("%s  %s  %s\n", keyStrs[i], countStrs[i], bar)
	}
}

func padColumn(values []string, colIdx int) {
	maxLen := 0
	for _, value := range values {
		maxLen = essentials.MaxInt(maxLen, len(value))
	}
	for i, val := range values {
		for len(val) < maxLen {
			if colIdx == 0 {
				val = " " + val
			} else {
				val += " "
			}
		}
		values[i] = val
	}
}

func printHeading(name string) {
	fmt.Println("---", name, "---")
}
