package main

import (
	"context"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
    log := getLogger()
    log.Println("Starting")

	conf := getConfig()
    secrets := getSecrets()
	ctx := context.Background()

    client := influxdb2.NewClient(conf.InfluexUrl, secrets.Token)
	_, err := client.Health(ctx)
	if err != nil {
		log.Println("Abbruch: InfluxDB nicht erreichbar, ist es gestartet?")
		panic(err)
	}

	startTime := time.Now()
	log.Println("Start time:", startTime.Format("2006-01-02 15:04:05.000"))
	log.Println("Links pro Durchlauf:", conf.LinksProDurchlauf)
	
	urls := conf.URLS

	for _, url := range urls {
		log.Println("URL:", url)
		allLinks := getLinks(url)
		log.Println("Anzahl Links:", len(allLinks))

		var wg sync.WaitGroup

		for i := 0; i < len(allLinks); i += conf.LinksProDurchlauf {
			links := allLinks[i : i+conf.LinksProDurchlauf]
			start := i
			end := i + conf.LinksProDurchlauf
			wg.Add(1)
			go func(links []string, start int, end int) {
				log.Println("Starting range from", start, "to", end)
				data := getDataEcoon(links)
				if conf.TeilExporte {
					saveData(data, start, end)
				}
				saveTimeSerie(data, start, end)
				wg.Done()
			}(links, start, end)
		}

		wg.Wait()

		if conf.AlleExportieren {
			exportAll()
		}
		exportPriceDifferences()

		log.Println("Finished after", time.Since(startTime))
	}
}
