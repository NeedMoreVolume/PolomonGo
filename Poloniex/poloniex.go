package poloniex

import (
  "net/http"
  "strconv"
  "context"
  "log"
  "fmt"
  "encoding/json"
  "time"
  "io/ioutil"
  "math"
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/mongo/options"
)

// base variables
var myClient = &http.Client{Timeout: 10 * time.Second}
const baseurl = "https://poloniex.com/public?command=returnChartData&currencyPair="


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

func CheckDate(date *float64) bool {
  return (int64(*date) + 86400 < time.Now().Unix())
}

func GetCandlestickData(client *mongo.Client, startTime int64, market *string) () {
  startTime = startTime/1000000000
  collection := client.Database("poloniex").Collection(*market)
  // check the last date in the collections
  filter := bson.M(nil)
  count, countErr := collection.CountDocuments(context.Background(), filter)
  if countErr != nil { log.Fatal(countErr) }
  if count > 0 {
    count -= 1
  } else {
    count = 0
  }
  options := options.FindOptions{}
  options.SetSkip(count)
  cur, err := collection.Find(context.Background(), filter, &options)
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
  var pair string
  switch *market {
  case "btc_etc":
    pair = "BTC_ETC"
  case "btc_eth":
    pair = "BTC_ETH"
  case "btc_ltc":
    pair = "BTC_LTC"
  case "btc_rep":
    pair = "BTC_REP"
  case "btc_xmr":
    pair = "BTC_XMR"
  case "btc_xrp":
    pair = "BTC_XRP"
  }
  if lastDate == 0 {
    url = baseurl + pair + "&start=0&end=9999999999&period=86400"
  } else {
    url = baseurl + pair + "&start=" + strconv.Itoa(lastDate + 86400) + "&end=9999999999&period=86400"
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
    if CheckDate(&stick.Date) {
      _, err := collection.InsertOne(context.Background(), stick)
      if err != nil {
        log.Fatal(err)
      }
    }
  }
  return
}

func Smabb(elements []Candlestick) {
  var sum,sma,sd float64
  for i:=0; i<20; i++ {
    sum += elements[i].Close
  }
  sma = sum/20
  for i:=0; i<20; i++ {
   sd += math.Pow(elements[i].Close - sma, 2)
  }
  sd = math.Sqrt(sd/20)
  upperband := sma + (sd * 2)
  lowerband := sma - (sd * 2)
  fmt.Printf("--------------------------------\n")
  fmt.Printf("    Date    :  %.0f\n", elements[19].Date)
  fmt.Printf("     SD     :  %.8f\n", sd)
  fmt.Printf("  Upper BB  :  %.8f\n", upperband)
  fmt.Printf("     SMA    :  %.8f\n", sma)
  fmt.Printf("  Lower BB  :  %.8f\n", lowerband)
  fmt.Printf("--------------------------------\n")
  return
}

func GetSMAandBB(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.M(nil)
  count, err := collection.Count(context.Background(), filter)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  elements := make([]Candlestick, count)
  i := 0
  // build array of stick data
  for cur.Next(context.Background()) {
    err := cur.Decode(&elements[i])
    if err != nil { log.Fatal(err) }
    i++
  }
  marker := 20
  for marker <= i {
    frame := elements[marker-20: marker]
    Smabb(frame)
    marker++
  }
  return
}

func ListCandles(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.M(nil)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  defer cur.Close(context.Background())
  for cur.Next(context.Background()) {
    var element Candlestick
    err := cur.Decode(&element)
    if err != nil { log.Fatal(err) }
    fmt.Printf("-------------------------------\n")
    fmt.Printf("    Date     :    %.0f\n", element.Date)
    fmt.Printf("    High     :    %9.8f\n", element.High)
    fmt.Printf("    Low      :    %9.8f\n", element.Low)
    fmt.Printf("    Open     :    %9.8f\n", element.Open)
    fmt.Printf("    Close    :    %9.8f\n", element.Close)
    fmt.Printf("   Average   :    %9.8f\n", element.WeightedAverage)
    fmt.Printf("   Volume    : %9.8f\n", element.QuoteVolume)
    fmt.Printf("-------------------------------\n")
  }
  return
}

func CalculateTenkanSen(elements []Candlestick) float64 {
  var high9, low9 float64
  for _, element := range elements {
    if element.High > high9 {
      high9 = element.High
    }
    if element.Close < low9 || low9 == 0 {
      low9 = element.Close
    }
  }
  ts := (high9 + low9)/2
  return ts
}

func CalculateKijunSen(elements []Candlestick) float64 {
  var high26, low26 float64
  for _, element := range elements {
    if element.High > high26 {
      high26 = element.High
    }
    if element.Close < low26 || low26 == 0 {
      low26 = element.Close
    }
  }
  ks := (high26 + low26)/2
  return ks
}

func CalculateSenkouSpanB (elements []Candlestick) float64 {
  var high52, low52 float64
  for _, element := range elements {
    if element.High > high52 {
      high52 = element.High
    }
    if element.Close < low52 || low52 == 0 {
      low52 = element.Close
    }
  }
  ssb := (high52+low52)/2
  return ssb
}

func CalculateIchimokuCloud(elements []Candlestick) {
  var ts float64
  switch {
  case len(elements)==9:
    ts = CalculateTenkanSen(elements)
  case len(elements)==26:
    ks := CalculateKijunSen(elements)
    ts = CalculateTenkanSen(elements[16:25])
    ssa := (ts + ks)/2
    fmt.Printf("Senkou Span A : %.8f\n", ssa)
    fmt.Printf("Kijun-sen : %.8f\n", ks)
  case len(elements)==52:
    ts = CalculateTenkanSen(elements[42:51])
    ks := CalculateKijunSen(elements[25:51])
    ssb := CalculateSenkouSpanB(elements)
    ssa := (ts + ks)/2
    fmt.Printf("Senkou Span B : %.8f\n", ssb)
    fmt.Printf("Senkou Span A : %.8f\n", ssa)
    fmt.Printf("Kijun-sen : %.8f\n", ks)
  }
  fmt.Printf("Tenkan-sen : %.8f\n", ts)
  return
}

func GetIchimokuCloud(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.M(nil)
  count, err := collection.Count(context.Background(), filter)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  elements := make([]Candlestick, count)
  i := 0
  // build array of stick data
  for cur.Next(context.Background()) {
    err := cur.Decode(&elements[i])
    if err != nil { log.Fatal(err) }
    i++
  }
  marker := 9
  for marker <= i {
    fmt.Println()
    if marker-52 >= 0 {
      //we have 52 previous days of data, so we can calculate the full cloud
      frame := elements[marker-52: marker]
      CalculateIchimokuCloud(frame)
    } else if marker-26 >= 0 {
      // we have 26 previous days of data
      frame := elements[marker-26: marker]
      CalculateIchimokuCloud(frame)
    } else if marker-9 >= 0 {
      // we have 9 days of data to start the ichimoku cloud calc.
      frame := elements[marker-9: marker]
      CalculateIchimokuCloud(frame)
    }
    if (marker+26) < i {
      fmt.Printf("Chikou Span : %.8f\n", elements[marker+26].Close)
    }
    fmt.Println()
    marker++
  }
  return
}
