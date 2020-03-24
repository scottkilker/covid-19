package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
)

type covidResponse struct {
	ID             string      `json:"id"`
	DisplayName    string      `json:"displayName"`
	Areas          []areas     `json:"areas"`
	TotalConfirmed int64       `json:"totalConfirmed"`
	TotalDeaths    int64       `json:"totalDeaths"`
	TotalRecovered int64       `json:"totalRecovered"`
	LastUpdated    time.Time   `json:"lastUpdated"`
	Lat            int         `json:"lat"`
	Long           int         `json:"long"`
	Country        string      `json:"country"`
	ParentID       interface{} `json:"parentId"`
}

type areas struct {
	ID             string    `json:"id"`
	DisplayName    string    `json:"displayName"`
	Areas          []areas   `json:"areas"`
	TotalConfirmed int64     `json:"totalConfirmed"`
	TotalDeaths    int64     `json:"totalDeaths"`
	TotalRecovered int64     `json:"totalRecovered"`
	LastUpdated    time.Time `json:"lastUpdated"`
	Lat            float64   `json:"lat"`
	Long           float64   `json:"long"`
	Country        string    `json:"country"`
	ParentID       string    `json:"parentId"`
}

func main() {

	boolPtr := flag.Bool("detail", false, "show detailed information")
	flag.Parse()

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.bing.com/covid/data", nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	data := covidResponse{}

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Fatal(err)
	}

	if *boolPtr == true {
		parse(data.Areas, "")
	}

	displayInformation(data)
}

func displayInformation(response covidResponse) {
	fmt.Println("TOTAL CONFIRMED CASES (", humanize.Comma(response.TotalConfirmed), ")")
	fmt.Println("Active:\t\t", humanize.Comma(response.TotalConfirmed-response.TotalRecovered-response.TotalDeaths))
	fmt.Println("Recovered:\t", humanize.Comma(response.TotalRecovered))
	fmt.Println("Fatal:\t\t", humanize.Comma(response.TotalDeaths))
	fmt.Println("Updated:\t", humanize.Time(response.LastUpdated))
}

func parse(root []areas, indent string) {

	for index, child := range root {
		add := "│\t"

		if index == len(root)-1 {
			fmt.Printf(indent + "└── ")
			add = "   "
		} else {
			fmt.Printf(indent + "├── ")
		}

		fmt.Println(child.DisplayName, ":\t", humanize.Comma(child.TotalConfirmed))

		parse(child.Areas, indent+add)
	}
}
