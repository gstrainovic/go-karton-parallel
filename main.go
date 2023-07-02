package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	Title            string
	Sku              string
	PiecesPerPalette int
	Values           []Value
}

type Value struct {
	LinkText int
	Value    float64
}

func main() {
	startTime := time.Now()
	conf := getConfig()
	fmt.Println("Start time:", startTime.Format("2006-01-02 15:04:05.000"))
	fmt.Println("URL:", conf.URL)
	fmt.Println("Domain:", conf.Domain)
	fmt.Println("Links pro Durchlauf:", conf.LinksProDurchlauf)

	allLinks := getLinks(conf)
	fmt.Println("Anzahl Links:", len(allLinks))

	var wg sync.WaitGroup

	for i := 0; i < len(allLinks); i += conf.LinksProDurchlauf {
		links := allLinks[i : i+conf.LinksProDurchlauf]
		start := i
		end := i + conf.LinksProDurchlauf
		wg.Add(1)
		go func(links []string, start int, end int) {
			fmt.Println("Starting range from", start, "to", end)
			data := getData(links)
			// saveData(data, start, end)
			saveTimeSerie(data, start, end)
			wg.Done()
		}(links, start, end)
	}

	wg.Wait()

	exportAll()
	exportPriceDifferences()

	fmt.Println("Finished after", time.Since(startTime))
}
