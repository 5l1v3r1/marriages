package main

import (
	"github.com/unixpickle/essentials"
)

// MarriagesAtDate scrapes all of the marriages at a given
// MM/DD/YYYY string.
func MarriagesAtDate(date string) (m []*Marriage, err error) {
	defer essentials.AddCtxTo("MarriagesAtDate", &err)
	page, err := Search("%", "%", date, date)
	if err != nil {
		return nil, err
	}
	visited := map[string]bool{"1": true}
	seenIDs := map[string]bool{}
PageLoop:
	for {
		for _, marriage := range page.Marriages {
			// We may see a marriage more than once, since
			// clicking ... brings us to a page, and then
			// we may visit that page again.
			if !seenIDs[marriage.LicenseID] {
				seenIDs[marriage.LicenseID] = true
				m = append(m, marriage)
			}
		}
		for _, link := range continueLinks(page) {
			if !visited[link.Title] {
				if link.Title != "..." {
					visited[link.Title] = true
				}
				page, err = page.GetLink(link)
				if err != nil {
					return nil, err
				}
				continue PageLoop
			}
		}
		break
	}
	return
}

func continueLinks(page *Page) []*PageLink {
	res := page.Pages
	if len(res) > 0 && res[0].Title == "..." {
		return res[1:]
	}
	return res
}
