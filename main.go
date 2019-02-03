package main

import (
  "context"
  "log"
  "time"
  "fmt"
  "github.com/NeedMoreVolume/PolomonGo/Poloniex"
  "github.com/NeedMoreVolume/PolomonGo/Menu"
  "github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
  fmt.Println(`
    Welcome!`)
  client, err := mongo.Connect(context.Background(), "mongodb://localhost:27017")
  err = client.Ping(context.Background(), nil)
  if err != nil {
    log.Fatal(err)
  }
  // connected to MongoDB
  var choice string = ""
  var subchoice string = ""
  for choice != "e" {
    menu.Menu(&choice)
    switch(choice) {
    case "1":
      menu.SubMenu(&subchoice)
      poloniex.GetCandlestickData(client, time.Now().Unix(), &subchoice)
      break
    case "2":
      menu.SubMenu(&subchoice)
      poloniex.ListCandles(client, &subchoice)
      break
    case "3":
      menu.SubMenu(&subchoice)
      poloniex.GetSMAandBB(client, &subchoice)
      break
    case "4":
      menu.SubMenu(&subchoice)
      poloniex.GetIchimokuCloud(client, &subchoice)
      break
    case "e":
      fmt.Println(`
  Goodbye!
        `)
      return
    }
  }
}
