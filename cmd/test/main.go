package main

import (
	"context"
	"destr4ct/summer/internal/crawler/rss"
	"fmt"
)

func main() {
	handler := rss.GetHandler()
	articles, err := handler.ParseSource(context.Background(), "https://rssexport.rbc.ru/rbcnews/news/20/full.rss")
	if err != nil {
		panic(err)
	}

	fmt.Println(articles[0])

}
