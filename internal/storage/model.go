package storage

import "time"

type SArticle struct {
	ID int

	// Source пока только в rss формате
	Source  string
	Content string

	Summary     string
	DateCreated time.Time
	HasSummary  bool
}

type User struct {
	ID int

	Username string
	TGID     string

	PreferredKeywords []string
	Sources           []string

	DateCreated time.Time
}

func NewUser(username, tgid string) *User {
	return &User{
		Username:          username,
		TGID:              tgid,
		PreferredKeywords: make([]string, 0),
		Sources:           make([]string, 0),
		DateCreated:       time.Now(),
	}
}

func NewArticle(link, content string) *SArticle {
	return &SArticle{
		Source:      link,
		Content:     content,
		DateCreated: time.Now(),
		HasSummary:  false,
	}
}
