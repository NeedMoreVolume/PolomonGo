package menu

import (
  "os"
  "fmt"
  "bufio"
)

func Menu(choice *string) {
  fmt.Println(`
Please pick an option from the list below.
  1. Get candlestick data.
  2. List candlestick data.
  3. Calculate 20 day SMA and Bollinger Bands.
  4. Calculate standard Ichimoku Cloud.
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
