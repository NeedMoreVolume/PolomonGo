package main

import (
  "net/http"
  "strconv"
  "context"
  "log"
  "fmt"
  "bufio"
  "os"
  "encoding/json"
  "time"
  "io/ioutil"
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo"
  //"github.com/mongodb/mongo-go-driver/mongo/options"
)

// base variables
var myClient = &http.Client{Timeout: 10 * time.Second}
var baseurl = "https://poloniex.com/public?command=returnChartData&currencyPair="

type Candlestick struct {
  Date float64
  High float64
  Low float64
  Open float64
  Close float64
  Volume float64
  QuoteVolume float64
  WeightedAverage float64
}

func checkDate(date float64) bool {
  return (int64(date) + 86400 < time.Now().Unix())
}

func getBTC_ETH(client *mongo.Client, startTime int64) () {
  startTime = startTime/1000000000
  collection := client.Database("poloniex").Collection("btc_eth")
  // check the last date in the collections
  filter := bson.M(nil)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  defer cur.Close(context.Background())
  var lastDate int = 0
  for cur.Next(context.Background()) {
    var element Candlestick
    err := cur.Decode(&element)
    if err != nil { log.Fatal(err) }
    lastDate = int(element.Date)
  }
  if err := cur.Err(); err != nil {
		log.Fatal(err)
  }
  var url string
  if lastDate == 0 {
    url = baseurl + "BTC_ETH&start=0&end=9999999999&period=86400"
  } else {
    url = baseurl + "BTC_ETH&start=" + strconv.Itoa(lastDate + 86400) + "&end=9999999999&period=86400"
  }
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Fatal(err)
  }
  var data []Candlestick
  json.Unmarshal(body, &data)
  for _, stick := range data {
    if checkDate(stick.Date) {
      _, err := collection.InsertOne(context.Background(), stick)
      if err != nil {
        log.Fatal(err)
      }
    }
  }
}

func listBTC_ETH(client *mongo.Client) {
  collection := client.Database("poloniex").Collection("btc_eth")
  filter := bson.M(nil)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  defer cur.Close(context.Background())
  for cur.Next(context.Background()) {
    var element Candlestick
    err := cur.Decode(&element)
    if err != nil { log.Fatal(err) }
    high := fmt.Sprintf("%9.8f", element.High)
    low := fmt.Sprintf("%9.8f", element.Low)
    open := fmt.Sprintf("%9.8f", element.Open)
    close := fmt.Sprintf("%9.8f", element.Close)
    qv := fmt.Sprintf("%8.8f", element.QuoteVolume)
    wa := fmt.Sprintf("%9.8f", element.WeightedAverage)
    fmt.Println("--------------------------------")
    fmt.Println("  Date:             " + strconv.Itoa(int(element.Date)))
    fmt.Println("| High:             " + high + " |")
    fmt.Println("| Low:              " + low + " |")
    fmt.Println("| Open:             " + open + " |")
    fmt.Println("| Close:            " + close + " |")
    fmt.Println("| WeightedAverage:  " + wa + " |")
    fmt.Println("  QuoteVolume:  " + qv)
    fmt.Println("--------------------------------")
  }
}

func doMenu() (string) {
  var choice string
  fmt.Println("Welcome!")
  fmt.Println("Please pick an option from the list below.")
  fmt.Println("1. Get ETH/BTC chart/candlestick data.")
  fmt.Println("2. Get ETH/BTC Orderbook information.")
  fmt.Println("3. List ETH/BTC chart/candlestick data.")
  fmt.Println("4. Exit.")
  fmt.Println()
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  choice = scanner.Text()
  fmt.Println()
  return choice
}

func main() {
  client, err := mongo.Connect(context.Background(), "mongodb://localhost:27017")
  err = client.Ping(context.Background(), nil)
  if err != nil {
    log.Fatal(err)
  }
  // connected to MongoDB
  var choice string = "0"
  for choice != "4" {
    choice = doMenu()
    switch(choice) {
    case "1":
      getBTC_ETH(client, time.Now().Unix())
      break
    case "2":
      fmt.Println("This is not complete yet.")
      break
    case "3":
      listBTC_ETH(client)
      break
    }
  }
}
