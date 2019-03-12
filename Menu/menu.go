package menu

import (
  "os"
  "fmt"
  "time"
  "bufio"
  "github.com/NeedMoreVolume/PolomonGo/Poloniex"
  "github.com/mongodb/mongo-go-driver/mongo"
)

func MenuHandler(client *mongo.Client) {
  fmt.Println(`
    Welcome!`)
  var choice string = ""
  var subchoice string = ""
  for choice != "e" {
    Menu(&choice)
    switch(choice) {
    case "1":
      SubMenu(&subchoice)
      poloniex.PutDataInMongo(client, time.Now().Unix(), &subchoice)
      break
    case "2":
      SubMenu(&subchoice)
      poloniex.PutDataInCassandra(&subchoice)
      break
    case "3":
      SubMenu(&subchoice)
      poloniex.ListMongoCandles(client, &subchoice)
      break
    case "4":
      SubMenu(&subchoice)
      //poloniex.ListCassandraCandles()
      break
    case "5":
      SubMenu(&subchoice)
      poloniex.GetSMAandBB(client, &subchoice)
      break
    case "6":
      SubMenu(&subchoice)
      poloniex.GetIchimokuCloud(client, &subchoice)
      break
    case "7":
      SubMenu(&subchoice)
      poloniex.GetRsi(client, &subchoice)
      break
    case "e":
      fmt.Println(`
  Goodbye!
        `)
      return
    }
  }
}

func Menu(choice *string) {
  fmt.Println(`
Please pick an option from the list below.
  1. Get candlestick data for Mongo.
  2. Get candlestick data for Cassandra.
  3. List candlestick data for Mongo.
  4. List candlestick data for Cassandra.
  5. Calculate 20 day SMA and Bollinger Bands.
  6. Calculate Standard Ichimoku Cloud.
  7. Calculate 14 day RSI.
  e. Exit.
  `)
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  *choice = scanner.Text()
  fmt.Println()
  return
}

func SubMenu(choice *string) {
  fmt.Println(`
  Please specify a currency pair from the list below.

    btc_etc
    btc_eth
    btc_ltc
    btc_rep
    btc_xmr
    btc_xrp
    `)
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  *choice = scanner.Text()
  return
}
