package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gocolly/colly"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	TYPES "CRipper/model"
)

func ScrapeImage(v TYPES.Visit, cl chan string, ctx context.Context) {
	fmt.Println("HIT", v)

	c := colly.NewCollector()

	c.OnHTML("div#readerarea img", func(e *colly.HTMLElement) {
		srcUrl := e.Attr("src")
		ss := strings.Split(srcUrl, "/")
		fmt.Println("start download")
		r := DownloadFile(srcUrl, "./images/"+ss[len(ss)-1])
		fmt.Println(<-r)
		cuf := UploadFile(ctx, "images/"+ss[len(ss)-1], v.Comic+"/"+strconv.FormatInt(v.Chapter, 10)+"/"+ss[len(ss)-1])
		fmt.Println(<-cuf)
		err := os.Remove("images/" + ss[len(ss)-1])
		if err != nil {
			log.Fatal("failed while removing img: ", e)
		}
		defer close(r)
	})
	c.Visit(v.Url)
	c.Wait()
}

func ScrapeSiteForReleases(urls []string, targets map[string]int64) []TYPES.Visit {
	fList := []TYPES.Visit{}
	tc := make(chan TYPES.Visit, len(targets))

	for _, url := range urls {
		c := colly.NewCollector()
		switch url {
		case "https://asura.gg/":
			scrapeForAsura(c, targets, tc)

		case "https://manga4life.com/":
			fmt.Println("NYI")
		}
		err := c.Visit(url)
		c.Wait()

		if err != nil {
			fmt.Println("Error: ", err)
		}

		fList = append(fList, <-tc)
		fmt.Println(fList)
	}
	fmt.Println("returning from scrape")
	return fList
}

func scrapeForAsura(c *colly.Collector, targets map[string]int64, tc chan TYPES.Visit) {

	c.OnHTML("div.luf", func(e *colly.HTMLElement) {
		temps := strings.Split(e.Text, "\n")

		output := temps[2]
		output = strings.TrimFunc(output, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r) && output != ""
		})

		if slices.Contains(maps.Keys(targets), output) {
			fmt.Println("FOUND", output)
			cns := strings.Split(e.DOM.Children().Get(1).FirstChild.FirstChild.FirstChild.Data, " ")[1]
			chapNum, err := strconv.ParseInt(cns, 10, 0)
			if err != nil {
				fmt.Println("Error: while parsing chapter Number", err)
			}
			if chapNum > targets[output] {
				fmt.Println("New Version -", chapNum)
				href := e.DOM.Children().Get(1).FirstChild.FirstChild.Attr[0].Val
				v := TYPES.Visit{
					Chapter: chapNum,
					Url:     href,
					Comic:   output,
				}
				tc <- v
			}
		}
	})
}
