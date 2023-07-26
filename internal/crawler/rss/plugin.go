package rss

import "github.com/mmcdole/gofeed"

type bodyPlugin func(item *gofeed.Item) string

var plugins = []bodyPlugin{
	rbcPlugin,
}

func rbcPlugin(item *gofeed.Item) string {
	if rbcContainer, found := item.Extensions["rbc_news"]; found {
		return rbcContainer["full-text"][0].Value
	}
	return ""
}
