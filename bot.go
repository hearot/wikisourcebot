/*
   wikisourcebot, a simple Telegram bot which implements wikisource.org APIs
   Copyright (C) 2020 Hearot

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"cgt.name/pkg/go-mwclient"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	defaultLanguage = "en"
	maximum         = 15
	message         = "This bot can help you find and share links to Wikisource articles. It works automatically, no need to add it anywhere. Simply open any of your chats and type @Wikisource_bot + language code (en, es, etc.) + something in the message field. Then tap on a result to send.\n\nFor example, try typing <code>@Wikisource_bot en Divine Comedy</code> here."
	template        = "https://%s.wikisource.org/w/api.php"
	timeout         = 10
	token           = ""
)

var langs = [...]string{"ar", "de", "en", "es", "he", "it", "pl", "ru", "zh"}

var clients map[string]*mwclient.Client

var defaultParameters = map[string]string{
	"action":    "query",
	"format":    "json",
	"generator": "search",
	"gsrlimit":  strconv.Itoa(maximum),
	"gsrprop":   "snippet",
	"inprop":    "url",
	"prop":      "info",
}

func main() {
	clients = make(map[string]*mwclient.Client, len(langs))

	for _, lang := range langs {
		var client *mwclient.Client

		client, err := mwclient.New(fmt.Sprintf(template, lang), "Wikisourcebot")

		if err != nil {
			panic(err)
		}

		clients[lang] = client
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: timeout * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		bot.Send(m.Chat, message, tb.ModeHTML)
	})

	bot.Handle(tb.OnQuery, func(q *tb.Query) {
		splitString := strings.SplitN(q.Text, " ", 2)

		if len(splitString) > 1 {
			var client *mwclient.Client

			found := false
			lowerText := strings.ToLower(splitString[0])
			search := splitString[1]

			for _, lang := range langs {
				if lang == lowerText {
					client = clients[lang]
					found = true
					break
				}
			}

			if !found {
				client = clients[defaultLanguage]
				search = splitString[0] + " " + search
			}

			parameters := defaultParameters
			parameters["gsrsearch"] = search

			resp, err := client.Get(parameters)

			if err != nil {
				log.Println(err)
				return
			}

			pages, err := resp.GetObjectArray("query", "pages")

			if err != nil {
				log.Println(err)
				return
			}

			results := make(tb.Results, len(pages))

			for id, value := range pages {
				title, err := value.GetString("title")

				if err != nil {
					log.Println(err)
					return
				}

				url, err := value.GetString("fullurl")

				if err != nil {
					log.Println(err)
					return
				}

				result := &tb.ArticleResult{
					HideURL: true,
					Text:    url,
					Title:   title,
					URL:     url,
				}

				result.SetResultID(strconv.Itoa(id))

				results[id] = result
			}

			err = bot.Answer(q, &tb.QueryResponse{
				Results:   results,
				CacheTime: 600,
			})

			if err != nil {
				log.Println(err)
			}
		}
	})

	bot.Start()
}
