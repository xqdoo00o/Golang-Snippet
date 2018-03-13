package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/wcharczuk/go-chart"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var language_names = []string{
	"JavaScript",
	"Python",
	"Java",
	"Ruby",
	"PHP",
	"Cpp",
	"Csharp",
	"Go",
	"C",
	"TypeScript",
	"Swift",
	"Scala",
	"Objective-C",
	"R",
	"Perl",
}
var llen = len(language_names)

type repo struct {
	name  string
	stars int64
	xlink string
}

func getNo() {
	var i int
	fmt.Print("Please Input Language NO:")
	fmt.Scanf("%d", &i)
	if i > 0 && i <= llen {
		r := getResponse(language_names[i-1])
		s := proItems(r)
		render(s, language_names[i-1])
	} else {
		fmt.Println("error,not the right No")
		getNo()
	}

}
func toK(v interface{}) string {
	b := v.(float64)
	if b < 1000 {
		return strconv.Itoa(int(b))
	} else {
		n := float64(b / 1000.0)
		return strconv.FormatFloat(n, 'f', 1, 64) + "k"
	}
}
func render(items []repo, lName string) {
	var bars []chart.Value
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	for _, v := range items {
		bar := chart.Value{Value: float64(v.stars), Label: v.name}
		bars = append(bars, bar)
	}
	sbc := chart.BarChart{
		Title:      "Most Star " + lName + " Github projects - " + strconv.Itoa(year) + "." + strconv.Itoa(month),
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   600,
		BarWidth: 80,
		XAxis: chart.Style{
			Show:                true,
			TextHorizontalAlign: chart.TextHorizontalAlign(1),
			TextWrap:            chart.TextWrap(1),
			TextRotationDegrees: float64(45),
		},
		YAxis: chart.YAxis{
			ValueFormatter: toK,
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: bars,
	}
	outputFile, _ := os.OpenFile(lName+"_github_repos.svg", os.O_WRONLY|os.O_CREATE, 0777)
	outputWriter := bufio.NewWriter(outputFile)
	err := sbc.Render(chart.SVG, outputWriter)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
	}
	fmt.Println("success!")
	outputWriter.Flush()
	outputFile.Close()
}
func getResponse(name string) *simplejson.Json {
	fmt.Println(name)
	resp, err := http.Get("https://api.github.com/search/repositories?q=language:" + name + "&sort=stars")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json, err := simplejson.NewJson([]byte(string(body)))
	items := json.Get("items")
	return items
}
func proItems(items *simplejson.Json) []repo {
	im, err := items.Array()
	if err != nil {
		panic(err)
	}
	var repos []repo
	for _, v := range im {
		l := v.(map[string]interface{})
		n, _ := l["stargazers_count"].(json.Number).Int64()
		rp := repo{l["name"].(string), n, l["html_url"].(string)}
		repos = append(repos, rp)
	}
	return repos
}
func main() {
	fmt.Println("NO. Language")
	for i := 0; i < llen; i++ {
		fmt.Printf("%d. %s\n", i+1, language_names[i])
	}
	getNo()
}
