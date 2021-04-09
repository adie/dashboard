package main

import (
	"fmt"
	"log"
	"time"

	owm "github.com/briandowns/openweathermap"

	"github.com/pterm/pterm"
)

const (
	owmKey = "9b5b2bd8ccff0d3c3647368f7ab4dcac"
)

var forecast *owm.Forecast5WeatherData
var weather string
var needs_clear bool

func main() {
	weather = ""
	go startForecastTicker()
	clear()

	height := pterm.GetTerminalHeight()
	offset := (height - 7*3) / 2

	for {
		moveCursor(offset)
		now := time.Now()
		ts := bigFont(now.Format(" 15:04:05 "))
		pterm.DefaultCenter.Println(ts)

		ds := bigFont(now.Format(" Mon, Jan 2 "))
		pterm.DefaultCenter.Println(ds)

		ws := bigFont(fmt.Sprintf(" %s ", weather))
		pterm.DefaultCenter.Println(ws)

		time.Sleep(200 * time.Millisecond)
	}
}


func getForecast() {
	w, err := owm.NewForecast("5", "C", "RU", owmKey)
	if err != nil {
		log.Println(err)
	}
	w.DailyByName("St. Petersburg, Russia", 5)
	forecast = w.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	weather = fmt.Sprintf("%.f °C", forecast.List[0].Main.Temp)
}

func startForecastTicker() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				getForecast()
				needs_clear = true
			}
		}
	}()
	getForecast()
}

func bigFont(str string) string {
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(str)).Srender()
	return s
}

func moveCursor(y int) {
	pterm.Printf("\033[%dH", y)
}

func clear() {
	pterm.Print("\033[H\033[2J")
}
