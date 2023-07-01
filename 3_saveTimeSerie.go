package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func saveTimeSerie(
    returnArray []Item,
    start int,
    end int,
) {

    fmt.Println("Save range from", start, "to", end)

    conf := getConfig()
    secrets := getSecrets()

    client := influxdb2.NewClient(conf.InfluexUrl, secrets.Token)

    org := conf.InfluexOrg
    bucket := conf.InfluexBucket
    writeAPI := client.WriteAPIBlocking(org, bucket)

    for _, item := range returnArray {
        for _, value := range item.Values {
            tags := map[string]string{
                "sku":   item.Sku,
                "title": item.Title,
            }
			preis, err := strconv.ParseFloat(fmt.Sprintf("%v", value.Value), 64)
            if err != nil {
                log.Fatal(err)
            }
            fields := map[string]interface{}{
                "piecesPerPalette": item.PiecesPerPalette,
                "preis":            preis,
                "anzahl":           value.LinkText,
            }
            point := write.NewPoint(conf.InfluexBucket, tags, fields, time.Now())
            if err := writeAPI.WritePoint(context.Background(), point); err != nil {
                log.Fatal(err)
            }
        }
    }

    fmt.Println("Done saving range from", start, "to", end)
}