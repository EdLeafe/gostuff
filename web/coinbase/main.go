package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"` //float64
}

type Warning struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

type cbResp struct {
	Data     Data
	Warnings []Warning
}

/*
{"data":{
    "base":"BTC",
    "currency":"USD",
    "amount":"15745.77"
   },
 "warnings":[
    {"id":"missing_version",
     "message":"Please supply API version (YYYY-MM-DD) as CB-VERSION header",
     "url":"https://developers.coinbase.com/api#versioning"
    }
  ]
}
*/
func main() {
	args := os.Args[1:]
	cryptType := "BTC"
	if len(args) > 0 {
		cryptType = strings.ToUpper(args[0])
	}
	version := "2016-03-08"
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-USD/sell",
		cryptType)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("CB-VERSION", version)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//	fmt.Println(string(contentBytes))

	// (TESTING!)    contentBytes := `{"data":{"base":"BTC","currency":"USD","amount":"15688.15"},"warnings":[{"id":"missing_version","message":"Please supply API version (YYYY-MM-DD) as CB-VERSION header","url":"https://developers.coinbase.com/api#versioning"}]}`

	var cbresp cbResp
	err = json.Unmarshal([]byte(contentBytes), &cbresp)
	if err != nil {
		fmt.Println("ERROR", err)
	}
    fmt.Printf(cbresp.Data.Amount)
}
