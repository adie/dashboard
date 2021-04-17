package main

import (
	"io"
	"log"
	"time"

	"net/http"

	"github.com/pterm/pterm"
)

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

		//ws := bigFont(fmt.Sprintf(" %s ", weather))
		pterm.DefaultCenter.Println(weather)

		time.Sleep(200 * time.Millisecond)
	}
}

func getForecast() {
	resp, err := http.Get("https://wttr.in/?0A")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	weather = string(body)
}

func startForecastTicker() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				getForecast()
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
