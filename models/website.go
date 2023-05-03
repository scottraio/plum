package models

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
)

type Website struct {
	Sitemap   string
	PageLinks []string
}

type Webpage struct {
	URL         string
	Body        string
	Description string
	Title       string
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

func (w *Website) FetchAndParseHTMLIntoMarkdown(url string) Webpage {
	article, err := readability.FromURL(url, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		log.Fatal(err)
	}

	return Webpage{
		URL:         url,
		Body:        markdown,
		Description: article.Excerpt,
		Title:       article.Title,
	}
}

func (page *Webpage) ToString(chunk string) string {
	// Extract the text content from the article
	var sb strings.Builder
	sb.WriteString("Page Title: " + page.Title + "\n")
	sb.WriteString("Page URL: " + page.URL + "\n")
	sb.WriteString("Page Description: " + page.Description + "\n")
	sb.WriteString("Page Contents: \n" + chunk + "\n")
	return sb.String()
}
