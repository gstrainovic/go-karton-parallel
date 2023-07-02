package main

import (
	"strings"

	"github.com/gocolly/colly"
)

func getLinks(url string) []string {
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

func getDomain(url string) string {
	splitted := strings.Split(url, "/")
	return splitted[0] + "//" + splitted[2] + "/"
}