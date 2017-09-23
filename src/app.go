package main

import (
    "io"
    "net/http"
    "log"
    "os"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "errors"
)

const LISTEN_ADDRESS = ":9204"
const API_URL = "https://api.coinmarketcap.com/v1"

var testMode string;

type CoinMarketCapStatistics []struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Symbol string `json:"symbol"`
    Rank string `json:"rank"`
    PriceUSD string `json:"price_usd"`
    PriceBTC string `json:"price_btc"`
    PriceEur string `json:"price_eur"`
    VolumeUsd24h string `json:"24h_volume_usd"`
    VolumeEur24h string `json:"24h_volume_eur"`
    MarketCapUsd string `json:"market_cap_usd"`
    MarketCapEur string `json:"market_cap_eur"`
    AvailableSupply string `json:"available_supply"`
    TotalSupply string `json:"total_supply"`
    PercentChange1h string `json:"percent_change_1h"`
    PercentChange24h string `json:"percent_change_24h"`
    PercentChange7d string `json:"percent_change_7d"`
    LastUpdated string `json:"last_updated"`
}

func integerToString(value int) string {
    return strconv.Itoa(value)
}

func formatValue(key string, meta string, value string) string {
    result := key;
    if (meta != "") {
        result += "{" + meta + "}";
    }
    result += " "
    result += value
    result += "\n"
    return result
}

func queryData() (string, error) {
    var err error

    // Build URL
    url := API_URL + "/ticker/?convert=EUR&limit=20"

    // Perform HTTP request
    resp, err := http.Get(url);
    if err != nil {
        return "", err;
    }

    // Parse response
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return "", errors.New("HTTP returned code " + integerToString(resp.StatusCode))
    }
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    bodyString := string(bodyBytes)
    if err != nil {
        return "", err;
    }

    return bodyString, nil;
}

func getTestData() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    body, err := ioutil.ReadFile(dir + "/test.json")
    if err != nil {
        log.Fatal(err)
    }
    return string(body), nil
}

func metrics(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /metrics")

    var up = 1
    var jsonString string
    var err error

    if (testMode == "1") {
        jsonString, err = getTestData()
    } else {
        jsonString, err = queryData()
    }
    if err != nil {
        log.Print(err)
        up = 0
    }

    // Parse JSON
    jsonData := CoinMarketCapStatistics{}
    json.Unmarshal([]byte(jsonString), &jsonData)

    // Output
    io.WriteString(w, formatValue("coinmarketcap_up", "", integerToString(up)))

    for _, Coin := range jsonData {
        // Output
        io.WriteString(w, formatValue("coinmarketcap_rank", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.Rank))
        io.WriteString(w, formatValue("coinmarketcap_price_usd", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PriceUSD))
        io.WriteString(w, formatValue("coinmarketcap_price_btc", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PriceBTC))
        io.WriteString(w, formatValue("coinmarketcap_price_eur", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PriceEur))
        io.WriteString(w, formatValue("coinmarketcap_24h_volume_usd", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.VolumeUsd24h))
        io.WriteString(w, formatValue("coinmarketcap_24h_volume_eur", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.VolumeEur24h))
        io.WriteString(w, formatValue("coinmarketcap_market_cap_usd", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.MarketCapUsd))
        io.WriteString(w, formatValue("coinmarketcap_market_cap_eur", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.MarketCapEur))
        io.WriteString(w, formatValue("coinmarketcap_available_supply", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.AvailableSupply))
        io.WriteString(w, formatValue("coinmarketcap_total_supply", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.TotalSupply))
        io.WriteString(w, formatValue("coinmarketcap_percent_change_1h", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PercentChange1h))
        io.WriteString(w, formatValue("coinmarketcap_percent_change_24h", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PercentChange24h))
        io.WriteString(w, formatValue("coinmarketcap_percent_change_7d", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.PercentChange7d))
        io.WriteString(w, formatValue("coinmarketcap_last_updated", "id=\"" + Coin.ID + "\",name=\"" + Coin.Name + "\",symbol=\"" + Coin.Symbol + "\"", Coin.LastUpdated))
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    log.Print("Serving /index")
    html := `<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <title>CoinMarketCap Exporter</title>
    </head>
    <body>
        <h1>CoinMarketCap Exporter</h1>
        <p><a href="/metrics">Metrics</a></p>
    </body>
</html>`
    io.WriteString(w, html)
}

func main() {
    testMode = os.Getenv("TEST_MODE")
    if (testMode == "1") {
        log.Print("Test mode is enabled")
    }

    log.Print("CoinMarketCap exporter listening on " + LISTEN_ADDRESS)
    http.HandleFunc("/", index)
    http.HandleFunc("/metrics", metrics)
    http.ListenAndServe(LISTEN_ADDRESS, nil)
}