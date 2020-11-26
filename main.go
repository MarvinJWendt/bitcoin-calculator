package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
)

func main() {
	var qs = []*survey.Question{
		{
			Name:     "firstDate",
			Prompt:   &survey.Input{Message: "First date:"},
			Validate: survey.Required,
		},
		{
			Name:     "secondDate",
			Prompt:   &survey.Input{Message: "Second date:"},
			Validate: survey.Required,
		},
		{
			Name:     "amount",
			Prompt:   &survey.Input{Message: "Amount:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		FirstDate  string
		SecondDate string
		Amount     string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	first, _ := GetPrice(answers.FirstDate)
	_, second := GetPrice(answers.SecondDate)

	if strings.HasSuffix(answers.Amount, "$") {
		a, err := strconv.ParseFloat(strings.ReplaceAll(answers.Amount, "$", ""), 64)
		check(err)
		percent := a / first
		secondPercent := percent * second
		pterm.Info.Printf("On %s one bitcoin was worth %s dollars"+
			"\nOn %s one bitcoin was worth %s dollars."+
			"\nYou would have made %s dollars. \n", pterm.LightMagenta(answers.FirstDate), pterm.LightMagenta(first), pterm.LightMagenta(answers.SecondDate), pterm.LightMagenta(second), pterm.LightMagenta(secondPercent-a))

	} else {
		a, err := strconv.ParseFloat(answers.Amount, 64)
		check(err)
		pterm.Info.Printf("On %s one bitcoin was worth %s dollars"+
			"\nOn %s one bitcoin was worth %s dollars."+
			"\nYou would have made %s dollars. \n", pterm.LightMagenta(answers.FirstDate), pterm.LightMagenta(first), pterm.LightMagenta(answers.SecondDate), pterm.LightMagenta(second), pterm.LightMagenta((second-first)*a))
	}
}

func GetPrice(date string) (low, high float64) {
	url := pterm.Sprintf("https://api.coinpaprika.com/v1/coins/btc-bitcoin/ohlcv/historical?start=%s&end=%s", date, date)
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	json, err := ioutil.ReadAll(resp.Body)
	check(err)

	lowString := strings.ReplaceAll(strings.ReplaceAll(gjson.Get(string(json), "#.low").Raw, "[", ""), "]", "")
	highString := strings.ReplaceAll(strings.ReplaceAll(gjson.Get(string(json), "#.high").Raw, "[", ""), "]", "")

	low, _ = strconv.ParseFloat(lowString, 64)
	high, _ = strconv.ParseFloat(highString, 64)

	return
}

func check(err error) {
	if err != nil {
		pterm.Fatal.Println(err)
	}
}
