package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/cli"
)

func main() {
	ggl := cli.NewCLI("ggl", "0.0.0")
	ggl.Args = os.Args[1:]
	ggl.Commands = map[string]cli.CommandFactory{
		"search": func() (cli.Command, error) {
			return &Ggl{}, nil
		},
	}

	status, err := ggl.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(status)
}

type Ggl struct {
}

func (*Ggl) Help() string {
	return "for classic google search results"
}

func (*Ggl) Run(args []string) int {
	q := strings.Join(args, " ")

	var URL *url.URL
	URL, err := url.Parse("https://google.com/search")

	if err != nil {
		log.Println(err)
	}

	params := url.Values{}
	params.Add("q", q)
	URL.RawQuery = params.Encode()

	resp, err := http.Get(URL.String())

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	result, err := googleResultParser(resp)

	if err != nil {
		log.Println(err)
	}

	for _, item := range result {
		fmt.Println(item)
	}
	return 0
}

func (h *Ggl) Synopsis() string {
	return h.Help()
}

type GoogleResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

func googleResultParser(response *http.Response) ([]GoogleResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}
	results := []GoogleResult{}
	sel := doc.Find("div.g")
	rank := 1
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3.r")
		descTag := item.Find("span.st")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")
		if link != "" && link != "#" {
			result := GoogleResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank++
		}
	}
	return results, err
}
