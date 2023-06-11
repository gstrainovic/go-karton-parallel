package main

import (
    "fmt"
    "os"
    "sync"
    "time"

    "github.com/BurntSushi/toml"
)

type Item struct {
    Title  string
    Values []Value
}

type Value struct {
    LinkText int
    Value    float64
}

type Config struct {
    URL              string
    Domain           string
    LinksProDurchlauf int
}

func main() {
    b, err := os.ReadFile("./config.toml")
    if err != nil {
        panic(err)
    }

    var conf Config
    _, err = toml.Decode(string(b), &conf)
    if err != nil {
        panic(err)
    }

    startTime := time.Now()
    fmt.Println("Start time:", startTime.Format("2006-01-02 15:04:05.000"))
    fmt.Println("URL:", conf.URL)
    fmt.Println("Domain:", conf.Domain)
    fmt.Println("Links pro Durchlauf:", conf.LinksProDurchlauf)

    allLinks := getLinks(conf)
    fmt.Println("Anzahl Links:", len(allLinks))

    // Create a wait group to wait for all goroutines to finish
    var wg sync.WaitGroup

    // run parallel
    for i := 0; i < len(allLinks); i += conf.LinksProDurchlauf {
        links := allLinks[i : i+conf.LinksProDurchlauf]
        start := i
        end := i + conf.LinksProDurchlauf
        wg.Add(1)
        go func(links []string, start int, end int) {
            fmt.Println("Starting range from", start, "to", end)
            data := getData(links)
            saveData(data, start, end)
            wg.Done()
        }(links, start, end)
    }

    // Wait for all goroutines to finish
    wg.Wait()

    fmt.Println("Finished after", time.Since(startTime))
}