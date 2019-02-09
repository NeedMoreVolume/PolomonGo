package calculations

import (
  "math"
  "fmt"
  "github.com/NeedMoreVolume/PolomonGo/Structs"
)

// Ichimoku Cloud Calucations.

func calculateTenkanSen(elements []structs.Candlestick) float64 {
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

func calculateKijunSen(elements []structs.Candlestick) float64 {
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

func calculateSenkouSpanB(elements []structs.Candlestick) float64 {
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

func CalculateIchimokuCloud(elements []structs.Candlestick) ([]float64) {
  switch {
  case len(elements)==9:
    cloud := make([]float64, 4)
    cloud[0] = calculateTenkanSen(elements)
    cloud[1] = 0
    cloud[2] = 0
    cloud[3] = 0
    return cloud
  case len(elements)==26:
    cloud := make([]float64, 4)
    cloud[0] = calculateTenkanSen(elements[16:25])
    cloud[1] = calculateKijunSen(elements)
    cloud[2] = (cloud[0] + cloud[1])/2
    cloud[3] = 0
    return cloud
  case len(elements)==52:
    cloud := make([]float64, 4)
    cloud[0] = calculateTenkanSen(elements[42:51])
    cloud[1] = calculateKijunSen(elements[25:51])
    cloud[2] = (cloud[0] + cloud[1])/2
    cloud[3] = calculateSenkouSpanB(elements)
    return cloud
  }
  return make([]float64, 4)
}

// BB and SMA calcuations

func Smabb(elements []structs.Candlestick) {
  var sum,sma,sd float64
  length := len(elements)
  for i:=0; i<length; i++ {
    sum += elements[i].Close
  }
  sma = sum/float64(length)
  for i:=0; i<length; i++ {
   sd += math.Pow(elements[i].Close - sma, 2)
  }
  sd = math.Sqrt(sd/float64(length))
  upperband := sma + (sd * 2)
  lowerband := sma - (sd * 2)
  fmt.Printf("--------------------------------\n")
  fmt.Printf("    Date    :  %.0f\n", elements[length-1].Date)
  fmt.Printf("     SD     :  %.8f\n", sd)
  fmt.Printf("  Upper BB  :  %.8f\n", upperband)
  fmt.Printf("     SMA    :  %.8f\n", sma)
  fmt.Printf("  Lower BB  :  %.8f\n", lowerband)
  fmt.Printf("--------------------------------\n")
  return
}

// RSI Calculations (14 days)
func Rsi(elements []structs.Candlestick, avgLoss *float64, avgGain *float64) {
  var gains, losses, last float64
  gains, losses, last = 0, 0, 0
  length := len(elements)
  if *avgLoss == 0 || *avgGain == 0 {
    for i:=0; i<length; i++ {
      if elements[i].Close > last {
        gains += elements[i].Close - last
      } else {
        losses += last - elements[i].Close
      }
      last = elements[i].Close
    }
    *avgGain = gains/float64(length)
    *avgLoss = losses/float64(length)
  } else {
    if elements[length-2].Close > elements[length-1].Close {
      losses = elements[length-2].Close - elements[length-1].Close
      gains = 0
    } else {
      gains = elements[length-1].Close - elements[length-2].Close
      losses = 0
    }
    *avgGain = ((float64(length-1) * *avgGain) + gains) / float64(length)
    *avgLoss = ((float64(length-1) * *avgLoss) + losses) / float64(length)
  }
  rs := *avgGain / *avgLoss
  rsi := 100 - (100/(1+rs))
  fmt.Printf("--------------------------------\n Date  :  %.0f\n", elements[length-1].Date)
  fmt.Printf(" RS    :  %f\n", rs)
  fmt.Printf(" RSI   :  %f\n--------------------------------\n", rsi)
  return
}
