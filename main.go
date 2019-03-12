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
	dbFlagPtr := flag.String("db", "cassandra", "DB to operate on.")
	actionFlagPtr := flag.String("action", "menu", "Action for CLI to take (update DB or menu)")
	flag.Parse()
	client, err := mongo.Connect(context.Background(), "mongodb://localhost:27017")
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	switch *actionFlagPtr {
	case "update":
		switch *dbFlagPtr {
		case "cassandra":
			for _, market := range markets {
				poloniex.PutDataInCassandra(&market)
			}
			fmt.Println("Update complete!")
			return
		case "mongo":
			for _, market := range markets {
				poloniex.PutDataInMongo(client, time.Now().Unix(), &market)
			}
			fmt.Println("Update complete!")
			return
		case "":
			fmt.Println("Please enter a db to update.")
		}
	case "menu":
		menu.MenuHandler(client)
	}
	return
}
