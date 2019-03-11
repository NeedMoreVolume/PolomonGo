package menu

import (
  "os"
  "fmt"
  "bufio"
)

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
