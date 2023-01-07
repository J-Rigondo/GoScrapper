package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const BASE_URL string = "https://lolchess.gg/leaderboards"

type summoner struct {
	rank      string
	id        string
	point     string
	winRate   string
	topRate   string
	playCount string
	winCount  string
	top4Count string
	link      string
}

func main() {
	pageCount := getPageCount()
	var summoners []summoner

	for i := 0; i < pageCount; i++ {
		pageSummoners := getPage(i + 1)
		summoners = append(summoners, pageSummoners...)
	}

	writeToCsv(summoners)

}

func writeToCsv(summoners []summoner) {
	file, err := os.Create("summoners.csv")
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"랭크", "아이디", "포인트", "승률", "TOP4비율", "게임수", "승수", "TOP4수"}

	wErr := w.Write(headers)
	if wErr != nil {
		log.Fatal(wErr)
	}

	for _, summoner := range summoners {
		summonerSlice := []string{summoner.rank, summoner.id, summoner.point,
			summoner.winRate, summoner.topRate, summoner.playCount, summoner.winCount, summoner.top4Count}
		swErr := w.Write(summonerSlice)

		if swErr != nil {
			log.Fatal(swErr)
		}
	}

}

func getPageCount() int {
	pageCount := 0
	res, err := http.Get(BASE_URL)
	defer res.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode != 200 {
		log.Fatalln("failed status", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".pagination").Each(func(i int, selection *goquery.Selection) {
		pageCount = selection.Find("li").Length()
	})

	return pageCount
}

func getSummoner(selection *goquery.Selection, channel chan<- summoner) {
	rank := strings.TrimSpace(selection.Find(".summoner > span").Text())
	rank = strings.TrimLeft(rank, "#")
	id := strings.TrimSpace(selection.Find(".summoner > a").Text())
	link, _ := selection.Find(".summoner > a").Attr("href")
	point := strings.TrimSpace(selection.Find(".lp").Text())
	point = strings.TrimRight(point, " LP")
	winRate := strings.TrimSpace(selection.Find(".winrate").Text())
	topRate := strings.TrimSpace(selection.Find(".toprate").Text())
	playCount := strings.TrimSpace(selection.Find(".played").Text())
	winCount := strings.TrimSpace(selection.Find(".wins").Text())
	topCOunt := strings.TrimSpace(selection.Find(".tops").Text())

	fmt.Println(rank, id, link, point, winRate, topRate, playCount, winCount, topCOunt)

	channel <- summoner{rank,
		id, point, winRate,
		topRate, playCount, winCount, topCOunt, link}
}

func getPage(pageNum int) []summoner {
	channel := make(chan summoner)
	pageUrl := BASE_URL + "?mode=ranked&region=kr&page=" + strconv.Itoa(pageNum)
	var summoners []summoner

	res, err := http.Get(pageUrl)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalln("failed status", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	searchSummoners := doc.Find("table > tbody > tr")
	searchSummoners.Each(func(i int, selection *goquery.Selection) {
		go getSummoner(selection, channel)
	})

	for i := 0; i < searchSummoners.Length(); i++ {
		summoner := <-channel
		summoners = append(summoners, summoner)
	}

	return summoners
}
