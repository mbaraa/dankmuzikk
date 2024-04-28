package youtube

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var (
	apiKeyMatcher            = regexp.MustCompile(`innertubeApiKey":"([^"]*)`)
	parserScriptMatcher      = regexp.MustCompile(`ytInitialData[^{]*(.*?);\s*<\/script>`)
	otherParserScriptMatcher = regexp.MustCompile(`(?s)/ytInitialData"[^{]*(.*);\s*window\["ytInitialPlayerResponse"\]`)
)

func ScrapeSearch(query string) (results []SearchResult, err error) {
	url := "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)

	c := colly.NewCollector()
	// c.WithTransport(&http.Transport{
	// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// })

	c.Limit(&colly.LimitRule{DomainGlob: "*"})

	// token = ""

	c.OnHTML("html", func(h *colly.HTMLElement) {
		fmt.Println("items", h.ChildText("#items"))
		fmt.Println(h.Text[strings.Index(h.Text, "ytInitialPlayerResponse"):])
		fmt.Println(parserScriptMatcher.FindAllString(h.Text, 2))
		fmt.Println(otherParserScriptMatcher.FindAllString(h.Text, -1))
		fmt.Println(apiKeyMatcher.FindAllString(h.Text, -1)[1])
	})

	c.OnHTML("#items", func(h *colly.HTMLElement) {
		fmt.Println("yello")
	})
	err = c.Visit(url)
	c.Wait()

	return
}
