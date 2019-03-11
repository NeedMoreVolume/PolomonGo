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
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/mongo/options"
  "github.com/gocql/gocql"
  "github.com/NeedMoreVolume/PolomonGo/Structs"
  "github.com/NeedMoreVolume/PolomonGo/Calculations"
)

var myClient = &http.Client{Timeout: 10 * time.Second}
const baseurl = "https://poloniex.com/public?command=returnChartData&currencyPair="

func CheckDate(date *float64) bool {
  return (int64(*date) + 86400 < time.Now().Unix())
}

func GetCandlestickData(client *mongo.Client, startTime int64, market *string) () {
  startTime = startTime/1000000000
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.D{}
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
    var element structs.Candlestick
    err := cur.Decode(&element)
    if err != nil { log.Fatal(err) }
    lastDate = int(element.Date)
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
  if *lastDate == 0 {
    url = baseurl + pair + "&start=0&end=9999999999&period=86400"
  } else {
    url = baseurl + pair + "&start=" + strconv.Itoa(*lastDate + 86400) + "&end=9999999999&period=86400"
  }
  res, err := http.Get(url)
  if err != nil { log.Fatal(err) }
  body, err := ioutil.ReadAll(res.Body)
  if err != nil { log.Fatal(err) }
  return body
}

func PutDataInMongo(client *mongo.Client, startTime int64, market *string) () {
  startTime = startTime/1000000000
  collection := client.Database("poloniex").Collection(*market)
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
    var element structs.Candlestick
    err := cur.Decode(&element)
    if err != nil { log.Fatal(err) }
    lastDate = int(element.Date)
  }
  var data []structs.Candlestick
  body := GetCandlestickData(&lastDate, market)
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

func PutDataInCassandra(market *string) {
	// connect to cassandra.
	cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "cluster1"
    cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	if err != nil {log.Fatal(err)}
	var timestamp float64
	var lastDate int = 0
	defer session.Close()
	// check for the last date of stick data in cassandra.
	var iter *gocql.Iter
	query := `SELECT timestamp FROM poloniex_` + *market
	iter = session.Query(query).Iter()
	for iter.Scan(&timestamp) {
		if lastDate < int(timestamp) {
			lastDate = int(timestamp)
		}
	}
	if err := iter.Close(); err != nil {log.Fatal(err)}
	// get data to put in cassandra.
	var data []structs.Candlestick
	body := GetCandlestickData(&lastDate, market)
	json.Unmarshal(body, &data)
	for _, stick := range data {
		if CheckDate(&stick.Date) {
			//INSERT ONE STICK WEW
			query := `INSERT INTO poloniex_` + *market + ` (id, timestamp, high, low, open, close, volume, quotevolume, weightedaverage) VALUES (uuid(), ?, ?, ?, ?, ?, ?, ?, ?)`
			if err := session.Query(query, stick.Date, stick.High, stick.Low, stick.Open, stick.Close, stick.Volume, stick.QuoteVolume, stick.WeightedAverage).Exec(); err != nil {log.Fatal(err)}
		}
	}
}

func GetSMAandBB(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.D{}
  count, countErr := collection.CountDocuments(context.Background(), filter)
  if countErr != nil { log.Fatal(countErr) }
  cur, findErr := collection.Find(context.Background(), filter)
  if findErr != nil { log.Fatal(findErr) }
  elements := make([]structs.Candlestick, count)
  i := 0
  for cur.Next(context.Background()) {
    err := cur.Decode(&elements[i])
    if err != nil { log.Fatal(err) }
    i++
  }
  marker := 20
  for marker <= i {
    frame := elements[marker-20: marker]
    calculations.Smabb(frame)
    marker++
  }
  return
}

func ListMongoCandles(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.D{}
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  defer cur.Close(context.Background())
  for cur.Next(context.Background()) {
    var element structs.Candlestick
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

func GetIchimokuCloud(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.D{}
  count, err := collection.CountDocuments(context.Background(), filter)
  cur, err := collection.Find(context.Background(), filter)
  if err != nil { log.Fatal(err) }
  elements := make([]structs.Candlestick, count)
  i := 0
  for cur.Next(context.Background()) {
    err := cur.Decode(&elements[i])
    if err != nil { log.Fatal(err) }
    i++
  }
  marker := 9
  for marker <= i {
    fmt.Println()
    if marker-52 >= 0 {
      frame := elements[marker-52: marker]
      cloud := calculations.CalculateIchimokuCloud(frame)
      fmt.Printf("Tenkan-sen : %.8f\nKijun-sen : %.8f\nSenkou Span A : %.8f\nSenkou Span B : %.8f\n",cloud[0], cloud[1], cloud[2], cloud[3])
    } else if marker-26 >= 0 {
      frame := elements[marker-26: marker]
      cloud := calculations.CalculateIchimokuCloud(frame)
      fmt.Printf("Tenkan-sen : %.8f\nKijun-sen : %.8f\nSenkou Span A : %.8f\nSenkou Span B : %.8f\n",cloud[0], cloud[1], cloud[2], cloud[3])
    } else if marker-9 >= 0 {
      frame := elements[marker-9: marker]
      cloud := calculations.CalculateIchimokuCloud(frame)
      fmt.Printf("Tenkan-sen : %.8f\nKijun-sen : %.8f\nSenkou Span A : %.8f\nSenkou Span B : %.8f\n",cloud[0], cloud[1], cloud[2], cloud[3])
    }
    if (marker+26) < i {
      fmt.Printf("Chikou Span : %.8f\n", elements[marker+26].Close)
    }
    fmt.Println()
    marker++
  }
  return
}

func GetRsi(client *mongo.Client, market *string) {
  collection := client.Database("poloniex").Collection(*market)
  filter := bson.D{}
  count, countErr := collection.CountDocuments(context.Background(), filter)
  if countErr != nil { log.Fatal(countErr) }
  cur, findErr := collection.Find(context.Background(), filter)
  if findErr != nil { log.Fatal(findErr) }
  elements := make([]structs.Candlestick, count)
  i := 0
  for cur.Next(context.Background()) {
    err := cur.Decode(&elements[i])
    if err != nil { log.Fatal(err) }
    i++
  }
  marker := 14
  var avgLoss, avgGain float64
  avgLoss, avgGain = 0, 0
  for marker <= i {
    frame := elements[marker-14: marker]
    calculations.Rsi(frame, &avgLoss, &avgGain)
    marker++
  }
  return
}
