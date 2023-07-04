package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type SitemapIndex struct {
    Locations []string `xml:"sitemap>loc"`
}

type Sitemap struct {
    Locations []string `xml:"url>loc"`
}

func getLinks(url string) []string {

	// if the url is xml, parse it
	if strings.Contains(url, ".xml") {
		return XMLparseURLs(url)
	}

	domain := getDomain(url)
	var links []string

	collector := colly.NewCollector()
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// the link must be a link without a domain and contain 'x'
		if strings.Contains(link, "x") && !strings.Contains(link, "http") {
			fullPathURL := domain + link
			links = append(links, fullPathURL)
		}
	})
	err := collector.Visit(url)
	if err != nil {
		panic(err)
	}

	return links
}

func XMLparseURLs(url string) []string {
	resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    bytes, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var s Sitemap
    xml.Unmarshal(bytes, &s)

    return s.Locations

}

func getDomain(url string) string {
	splitted := strings.Split(url, "/")
	return splitted[0] + "//" + splitted[2] + "/"
}