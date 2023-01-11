package main

import (
	"JobScrapper/scrapper"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/home", func(c echo.Context) error {
		return c.File("home.html")
	})
	e.GET("/summoners/csv", func(c echo.Context) error {
		defer os.Remove("summoners.csv")
		scrapper.Scrape()
		return c.Attachment("summoners.csv", "summoner_list.csv")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
