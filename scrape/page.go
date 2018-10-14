package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/unixpickle/essentials"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const PageURL = "http://secureprod.phila.gov/wills/marriagesearch.aspx"

var client = http.Client{}

func init() {
	client.Jar, _ = cookiejar.New(nil)
}

type Marriage struct {
	Applicant1 string
	Applicant2 string
	Date       string
	LicenseID  string
}

type PageLink struct {
	Title    string
	Target   string
	Argument string
}

type Page struct {
	FormValues url.Values
	Marriages  []*Marriage
	Pages      []*PageLink
}

// Search gets the first page for the given date range and
// groom name.
func Search(firstName, lastName, dateStart, dateEnd string) (p *Page, err error) {
	defer essentials.AddCtxTo("Search", &err)
	page, err := fetchPage(nil)
	if err != nil {
		return nil, err
	}
	page.FormValues.Set("txtGROOM_FIRST", firstName)
	page.FormValues.Set("txtGROOM_LAST", lastName)
	page.FormValues.Set("txtMARRIAGE_DATEFrom", dateStart)
	page.FormValues.Set("txtMARRIAGE_DATETo", dateEnd)
	page.FormValues.Set("btnOpenBrowse", "Search")
	return fetchPage(page.FormValues)
}

// GetLink gets the page that results from clicking the
// link.
func (p *Page) GetLink(link *PageLink) (*Page, error) {
	newForm := url.Values{}
	for k, v := range p.FormValues {
		newForm[k] = append([]string{}, v...)
	}
	newForm.Set("__EVENTTARGET", strings.Replace(link.Target, "$", ":", -1))
	newForm.Set("__EVENTARGUMENT", link.Argument)
	page, err := fetchPage(newForm)
	return page, essentials.AddCtx("GetLink", err)
}

func fetchPage(formValues url.Values) (*Page, error) {
	var resp *http.Response
	var err error
	if formValues == nil {
		resp, err = client.Get(PageURL)
	} else {
		postBody := bytes.NewReader([]byte(formValues.Encode()))
		resp, err = client.Post(PageURL, "application/x-www-form-urlencoded", postBody)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	parsed, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	newForm, err := findFormValues(parsed)
	if err != nil {
		return nil, err
	}

	marriages, err := findMarriages(parsed)
	if err != nil {
		return nil, err
	}

	pages, err := findPages(parsed)
	if err != nil {
		return nil, err
	}

	return &Page{
		FormValues: newForm,
		Marriages:  marriages,
		Pages:      pages,
	}, nil
}

func findFormValues(parsed *html.Node) (url.Values, error) {
	form, ok := scrape.Find(parsed, scrape.ById("frmMarriageSearch"))
	if !ok {
		return nil, errors.New("no search form found")
	}
	newForm := url.Values{}
	for _, field := range scrape.FindAll(form, scrape.ByTag(atom.Input)) {
		if scrape.Attr(field, "name") == "image" {
			continue
		}
		name := scrape.Attr(field, "name")
		if !strings.HasPrefix(name, "btn") {
			newForm.Add(name, scrape.Attr(field, "value"))
		}
	}
	return newForm, nil
}

func findMarriages(parsed *html.Node) ([]*Marriage, error) {
	var marriages []*Marriage
	for _, className := range []string{"rowtext", "rowtextb"} {
		for _, row := range scrape.FindAll(parsed, scrape.ByClass(className)) {
			columns := scrape.FindAll(row, scrape.ByTag(atom.Td))
			if len(columns) < 5 {
				continue
			}
			marriages = append(marriages, &Marriage{
				Applicant1: strings.TrimSpace(scrape.Text(columns[0])),
				Applicant2: strings.TrimSpace(scrape.Text(columns[1])),
				Date:       strings.TrimSpace(scrape.Text(columns[2])),
				LicenseID:  strings.TrimSpace(scrape.Text(columns[3])),
			})
		}
	}
	return marriages, nil
}

func findPages(parsed *html.Node) ([]*PageLink, error) {
	pages, ok := scrape.Find(parsed, scrape.ByClass("pagers"))
	if !ok {
		return []*PageLink{}, nil
	}
	pageRegexp := regexp.MustCompile("^javascript:__doPostBack\\('(.*)',''\\)$")
	var links []*PageLink
	for _, pageLink := range scrape.FindAll(pages, scrape.ByTag(atom.A)) {
		href := scrape.Attr(pageLink, "href")
		matches := pageRegexp.FindStringSubmatch(href)
		if matches == nil {
			return nil, errors.New("invalid page link")
		}
		links = append(links, &PageLink{
			Title:  strings.TrimSpace(scrape.Text(pageLink)),
			Target: matches[1],
		})
	}
	return links, nil
}
