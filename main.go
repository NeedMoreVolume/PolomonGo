package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/NeedMoreVolume/PolomonGo/Menu"
	"github.com/NeedMoreVolume/PolomonGo/Poloniex"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"time"
)

func main() {
	markets := [6]string{"btc_etc", "btc_eth", "btc_ltc", "btc_rep", "btc_xmr", "btc_xrp"}
	actionFlagPtr := flag.String("action", "menu", "Available actions: menu, update")
	flag.Parse()
	client, err := mongo.Connect(context.Background(), "mongodb://localhost:27017")
	err = client.Ping(context.Background(), nil)
	if err != nil {log.Fatal(err)}
  if *actionFlagPtr == "menu" {
    menu.MenuHandler(client)
  } else if *actionFlagPtr == "update" {
    for _, market := range markets {
      poloniex.PutDataInMongo(client, time.Now().Unix(), &market)
      poloniex.PutDataInCassandra(&market)
    }
    fmt.Println("All DB updates complete!")
  } else {
    fmt.Println("Please use -h for help.")
  }
	return
}
