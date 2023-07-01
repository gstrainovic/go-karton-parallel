package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// "os"

	"github.com/tealeg/xlsx"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func exportDiffs() {
	conf := getConfig()
	secrets := getSecrets()
	client := influxdb2.NewClient(conf.InfluexUrl, secrets.Token)
	org := conf.InfluexOrg
	// bucket := conf.InfluexBucket
	queryAPI := client.QueryAPI(org)

	// export all skus which have a price difference
	query := `from(bucket: "karton.eu")
				|> range(start: -10m)
				|> filter(fn: (r) => r._measurement == "karton.eu")
				|> filter(fn: (r) => r._field == "preis")
				|> filter(fn: (r) => r._value != 0)
				|> group(columns: ["sku"])
				|> distinct(column: "sku")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	for results.Next() {
		fmt.Println(results.Record().ValueByKey("sku"))
	}

}

func displaySkuWithPriceChange() {
	conf := getConfig()
	secrets := getSecrets()

	client := influxdb2.NewClient(conf.InfluexUrl, secrets.Token)
	queryAPI := client.QueryAPI(conf.InfluexOrg)

	query := fmt.Sprintf(`from(bucket: "%s")
        |> range(start: -24h)
        |> filter(fn: (r) => r._measurement == "karton.eu")
        |> distinct(column: "sku")`, conf.InfluexBucket)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for result.Next() {
		if result.TableChanged() {
			fmt.Printf("\n%s\n", result.TableMetadata().String())
		}
		fmt.Printf("%v\n", result.Record().Value())
	}

	if result.Err() != nil {
		log.Fatal(result.Err())
	}
}

func exportAll() {
	conf := getConfig()
	secrets := getSecrets()

	client := influxdb2.NewClient(conf.InfluexUrl, secrets.Token)
	queryAPI := client.QueryAPI(conf.InfluexOrg)

	query := fmt.Sprintf(`from(bucket: "%s")
        |> range(start: -24h)
        |> filter(fn: (r) => r._measurement == "karton.eu")
        |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
        |> map(fn: (r) => ({sku: r.sku, anzahl: r.anzahl, preis: r.preis, piecesPerPalette: r.piecesPerPalette, title: r.title}))`, conf.InfluexBucket)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new file
	file := xlsx.NewFile()

	// Create a new sheet
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new row
	row := sheet.AddRow()

	cell := row.AddCell()
	cell.Value = "Artikelnummer"

	cell = row.AddCell()
	cell.Value = "Anzahl"

	cell = row.AddCell()
	cell.Value = "Preis"

	cell = row.AddCell()
	cell.Value = "Stück pro Palette"

	cell = row.AddCell()
	cell.Value = "Titel"

	for result.Next() {
		row = sheet.AddRow()
		
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", result.Record().ValueByKey("sku"))

		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", result.Record().ValueByKey("anzahl"))

		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", result.Record().ValueByKey("preis"))

		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", result.Record().ValueByKey("piecesPerPalette"))

		cell = row.AddCell()
		cell.Value = result.Record().ValueByKey("title").(string)
	}

	// Save the Excel file named with date, time and miliseconds
	filename := fmt.Sprintf("Alle_%s.xlsx", time.Now().Format("2006-01-02_15-04-05.000"))
	err = file.Save(fmt.Sprintf("data/%s", filename))
	if err != nil {
		log.Fatal("Error saving Excel file:", err)
	}

	fmt.Println("Saved all data to Excel file:", filename)


}

