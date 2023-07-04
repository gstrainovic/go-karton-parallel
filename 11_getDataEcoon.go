package main

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func getDataEcoon(links []string) []Item {
    log := getLogger()
    var returnArray []Item

    failedLinks := []string{}

    c := colly.NewCollector()

    c.OnHTML("h2.product_title.entry-title", func(e *colly.HTMLElement) {
        title := e.Text

        var sku string
        e.DOM.ParentsUntil("div").Find(".sku").Each(func(_ int, s *goquery.Selection) {
            sku = s.Text()
        })

        var piecesPerPalette int
        e.DOM.ParentsUntil("div").Find("td.woocommerce-product-attributes-item__value p").Each(func(_ int, s *goquery.Selection) {
            if strings.Contains(s.Text(), "Stk") {
                piecesPerPaletteString := strings.TrimSuffix(strings.TrimSpace(s.Text()), " Stk")
                piecesPerPalette, _ = strconv.Atoi(piecesPerPaletteString)
            }
        })

        var valuesArray []Value
        e.DOM.ParentsUntil("table").Find("tr").Each(func(_ int, s *goquery.Selection) {
            linkTextNumber := s.Find("td.h1 a").Text()
            linkTextInt, err := strconv.Atoi(linkTextNumber)
            if err != nil {
                return
            }
            valueString := s.Find("td.bulk-price-cell").Eq(1).Text()
            valueFloat, err := strconv.ParseFloat(strings.Replace(valueString, ",", ".", -1), 64)
            if err == nil {
                valuesArray = append(valuesArray, Value{
                    LinkText: linkTextInt,
                    Value:    valueFloat,
                })
            }
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

    c.OnError(func(r *colly.Response, err error) {
        log.Println("Request URL:", r.Request.URL, "failed with response:")
        // save the url to parsing again
        failedLinks = append(failedLinks, r.Request.URL.String())
    })

    c.Visit("https://ecoon.de/produkt/graspapierkarton-550-x-300-x-550-350-300-mm-2-wellig/")

    // for _, link := range links {
        // c.Visit(link)
    // }

    if len(failedLinks) > 0 {
        log.Println("Try again failed links:", failedLinks)
        for _, link := range failedLinks {
            c.Visit(link)
        }
    }

    return returnArray
}