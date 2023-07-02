package main

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func getData(links []string) []Item {
	log := getLogger()
	var returnArray []Item

	failedLinks := []string{}

	collector := colly.NewCollector(
		// colly.Debugger(&debug.LogDebugger{}),
	)

	collector.OnHTML("h1", func(e *colly.HTMLElement) {
		title := e.Text
		sku := e.DOM.ParentsUntil("~").Find("p.product-sku span").Text()
		
		// find the number in cUNNummer id
		piecesPerPaletteStr := e.DOM.ParentsUntil("~").Find("#cUNNummer").Text()
		// remove all non-numeric characters
		piecesPerPaletteStr = strings.Map(func(r rune) rune {
			if r < '0' || r > '9' {
				return -1
			}
			return r
		}, piecesPerPaletteStr)


		piecesPerPalette := 0
		if piecesPerPaletteStr != "" {
			var err error
			piecesPerPalette, err = strconv.Atoi(strings.Split(piecesPerPaletteStr, " ")[0])
			if err != nil {
				log.Println("Error parsing piecesPerPalette with the string:", piecesPerPaletteStr, "\nError:", err)
			}
		}


		var valuesArray []Value 

		tableRows := e.DOM.ParentsUntil("~").Find("table tr")

		tableRows.Each(func(index int, element *goquery.Selection) {
			columns := element.Find("td")

			linkText := columns.Eq(0).Text()
			value := strings.Split(columns.Eq(2).Text(), " ")[0]

			if linkText == "" || value == "" {
				return
			}

			valueFloat, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 64)
			if err != nil {
				log.Println("Error parsing value:", err)
				return
			}

			linkTextNumber, err := strconv.Atoi(linkText)
			if err != nil {
				log.Println("Error parsing linkText:", err)
				return
			}

			valuesArray = append(valuesArray, Value{
				LinkText: linkTextNumber,
				Value:    valueFloat,
			})
		})

		if len(valuesArray) > 0 && title != "" {
			returnArray = append(returnArray, Item{
				Title:            title,
				Sku:              sku,
				PiecesPerPalette: piecesPerPalette,
				Values:           valuesArray,
			})
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:")
		// save the url to parsing again
		failedLinks = append(failedLinks, r.Request.URL.String())
	})

	for _, link := range links {
		collector.Visit(link)
	}

	if len(failedLinks) > 0 {
		log.Println("Try again failed links:", failedLinks)
		for _, link := range failedLinks {
			collector.Visit(link)
		}
	}

	return returnArray
}

