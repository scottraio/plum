package models

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Website struct {
	Sitemap   string
	PageLinks []string
}

type Sitemap struct {
	URLSet []URL `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func (w *Website) GetURLsFromSitemap() []string {
	resp, err := http.Get(w.Sitemap)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}
	}

	var sitemap Sitemap
	err = xml.Unmarshal(body, &sitemap)
	if err != nil {
		return []string{}
	}

	urls := make([]string, len(sitemap.URLSet))
	for i, url := range sitemap.URLSet {
		urls[i] = url.Loc
	}

	w.PageLinks = urls
	return w.PageLinks
}

func (w *Website) FetchAndParseHTML(url string) (*goquery.Document, error) {
	// Fetch the HTML content for the current URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML content using GoQuery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (w *Website) Split(str string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(str); i += chunkSize {
		end := i + chunkSize
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[i:end])
	}
	return chunks
}
