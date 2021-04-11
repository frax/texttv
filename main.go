/*
	Retrives and presents a Swedish text-tv page in the terminal

	todo:
	* Fix background color handling
	* Handle DH somehow
	* blank span "graphics" nees fixing

	by frax 2020-04-11
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

var colors map[string]color.Attribute

type Page struct {
	Num               string
	Title             string
	Content           []string
	Next_Page         string
	Prev_Page         string
	Date_Updated_Unix int
	Permalink         string
	// id                string
}

func main() {
	defer color.Unset()
	colors = initColorMap()

	pageNum := getCurrentPageNum()
	src := getHtml(pageNum)

	parseHtml(src)
}

func getCurrentPageNum() int {
	var p int
	flag.IntVar(&p, "page", 100, "page to display")
	flag.Parse()

	if len(os.Args) >= 2 && !strings.HasPrefix(os.Args[1], "--") {
		p, _ = strconv.Atoi(os.Args[1])
	}
	return p
}

func getHtml(pageNum int) string {
	url := fmt.Sprintf("http://api.texttv.nu/api/get/%d?app=gotextv", pageNum)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching page.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// body := []byte(`[{"num":"300","title":"1,2 milj\/dag","content":["<div class=\"root\"><span class=\"toprow\"> 300 SVT Text         M\u00e5ndag 05 apr 2021\n <\/span><span class=\"B bgB\"> <\/span><span class=\"B bgB\">                 <\/span><span class=\"W bgB\">       <\/span><span class=\"Y bgB\"> <\/span><span class=\"W bgB\">1,2 milj\/dag <\/span>\n <span class=\"B bgB\"> <\/span><span class=\"B bgB\">                 <\/span><span class=\"W bgB\">           <\/span><span class=\"W bgB\">          <\/span>\n <span class=\"B bgB\"> <\/span><span class=\"B bgB\">                 <\/span><span class=\"W bgB\">           <\/span><span class=\"W bgB\">          <\/span>\n <span class=\"Y\">                       <\/span><span class=\"W\">                <\/span>\n <span class=\"B bgB\"> <\/span><span class=\"B bgB\"> <\/span><span class=\"W bgB\">* = efter kl 18   Resultatb\u00f6rsen <a href=\"\/330\">330<\/a> <\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Fr\u00f6lundas f\u00f6rsta steg mot l\u00e5ng v\u00e5r <a href=\"\/303\">303<\/a><\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Stolpen r\u00e4ddade FBK i Malm\u00f6........<a href=\"\/304\">304<\/a><\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"W\"> <\/span><span class=\"G\">19.30 SVT24: Basketsemi, Lule\u00e5-A3 Ume\u00e5<\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Uppgifter: Zlatan n\u00e4ra f\u00f6rl\u00e4ngning <a href=\"\/305\">305<\/a><\/span>\n <span class=\"W\"> <\/span><span class=\"W\">                                      <\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Bryn\u00e4s klubbdirekt\u00f6r: \"Katastrof\"..<a href=\"\/306\">306<\/a><\/span>\n <span class=\"W\"> <\/span><span class=\"W\">...och laget l\u00e5ser in sig i kvalet <a href=\"\/307\">307<\/a><\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Polack tog rankingkliv efter titel <a href=\"\/308\">308<\/a><\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"W\"> <\/span><span class=\"W\">Fr\u00e5gan om VAR delar svensk fotboll <a href=\"\/309\">309<\/a><\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"W\"> <\/span><span class=\"W\">M\u00e5l av Forsling - Florida i topp...<a href=\"\/310\">310<\/a><\/span>\n <span class=\"W\">                                       <\/span>\n <span class=\"Y bgY\"> <\/span><span class=\"Y bgY\"> <\/span><span class=\"B bgY\">M\u00e5lservice <a href=\"\/376\">376<\/a>- * Fler rubriker <a href=\"\/301\">301<\/a>  <\/span>\n <span class=\"B bgB\"> <\/span><span class=\"B bgB\"> <\/span><span class=\"W bgB\">   Sport i SVT n\u00e4rmaste tiden <a href=\"\/399\">399<\/a>    <\/span>\n<\/div>"],"next_page":"301","prev_page":"299","date_updated_unix":1617638505,"permalink":"https:\/\/texttv.nu\/300\/resultatbors-sportnyheter-29971743","id":"29971743"}]`)

	var page []Page
	json.Unmarshal(body, &page)

	// fmt.Println(page[0].Title)
	return strings.ReplaceAll(page[0].Content[0], "\n", "")
}

func parseHtml(src string) {
	dom := html.NewTokenizer(strings.NewReader(src))

	// Ouput counter for \n
	counter := 0
	for {
		tt := dom.Next()
		if tt == html.ErrorToken {
			// log.Fatal("error")
			break
		}

		// if div or span and class set text color
		if tt == html.StartTagToken {
			tag, attr := dom.TagName()
			if attr && (string(tag) == "div" || string(tag) == "span") {
				for {
					attr, value, moreAttr := dom.TagAttr()

					if string(attr) == "class" {
						fg, _ := mapColor(string(value)) // discard background for now
						color.Set(fg)                    // https://github.com/fatih/color/issues/135
					}
					if !moreAttr {
						break
					}
				}
			}
		}

		// div or span endtag, disable colors
		if tt == html.EndTagToken {
			tag, _ := dom.TagName()
			if string(tag) == "div" || string(tag) == "span" {
				color.Unset()
			}
		}

		if tt == html.TextToken {
			for _, c := range string(dom.Text()) {
				fmt.Printf("%c", rune(c))
				counter++
				if counter == 40 {
					counter = 0
					fmt.Printf("\n")
				}
			}
		}
	}
}

func mapColor(class string) (color.Attribute, color.Attribute) {
	var fg, bg = color.FgWhite, color.BgBlack // Default
	for _, c := range strings.Split(class, " ") {
		if strings.HasPrefix(c, "bg") {
			bg = colors[c]
		} else {
			fg = colors[c]
		}
		if fg == 0 || bg == 0 {
			log.Fatal("No color mapped for ", c)
		}
	}
	return fg, bg
}

func initColorMap() map[string]color.Attribute {
	return map[string]color.Attribute{
		"root":       color.FgWhite,
		"toprow":     color.FgYellow,
		"added-line": color.FgYellow,
		"DH":         color.FgHiRed, // No, that's not right...
		"B":          color.FgBlue,
		"C":          color.FgCyan,
		"W":          color.FgWhite,
		"Y":          color.FgYellow,
		"R":          color.FgRed,
		"G":          color.FgGreen,
		"bgB":        color.BgBlue,
		"bgW":        color.BgWhite, // Hmm
		"bgR":        color.BgRed,
		"bgC":        color.BgCyan,
		"bgY":        color.BgYellow,
	}
}
