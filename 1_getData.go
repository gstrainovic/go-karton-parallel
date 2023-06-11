package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func getData(links []string) []Item {
	var returnArray []Item

	collector := colly.NewCollector(
		// colly.Debugger(&debug.LogDebugger{}),
	)

	collector.OnHTML("h1", func(e *colly.HTMLElement) {
		title := e.Text

		var valuesArray []Value

		tableRows := e.DOM.ParentsUntil("~").Find("table tr")

		tableRows.Each(func(index int, element *goquery.Selection) {
			columns := element.Find("td")

			linkText := columns.Eq(0).Text()
			value := strings.Split(columns.Eq(1).Text(), " ")[0]

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
				Title:  title,
				Values: valuesArray,
			})
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for _, link := range links {
		collector.Visit(link)
	}
	return returnArray
}
