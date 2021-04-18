package main

import (
	"fmt"
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
		ts := bigFont(now.Format(" 15:04:05 "), pterm.FgLightWhite)
		pterm.DefaultCenter.Println(ts)

		ds := bigFont(now.Format(" Mon, Jan 2 "), pterm.FgWhite)
		pterm.DefaultCenter.Println(ds)

		ws := bigFont(fmt.Sprintf(" %s ", weather), pterm.FgLightGreen)
		pterm.DefaultCenter.Println(ws)

		time.Sleep(200 * time.Millisecond)
	}
}

func getForecast() {
	resp, err := http.Get("https://wttr.in/?format=%t")
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

func bigFont(str string, color pterm.Color) string {
	s, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle(str, pterm.NewStyle(color))).
		Srender()
	return s
}

func moveCursor(y int) {
	pterm.Printf("\033[%dH", y)
}

func clear() {
	pterm.Print("\033[H\033[2J")
}
