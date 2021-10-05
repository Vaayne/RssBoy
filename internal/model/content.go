package model

import (
	"strings"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/tgraph"

	"github.com/SlyMarbo/rss"
)

// Content fetcher content
type Content struct {
	SourceID     uint   `json:"source_id"`
	HashID       string `gorm:"primary_key" json:"hash_id"`
	RawID        string `json:"raw_id"`
	RawLink      string `json:"raw_link"`
	Title        string `json:"title"`
	Description  string `gorm:"-"` //ignore to db
	TelegraphURL string `json:"telegraph_url"`
	// IsBroaded mean the content has been broad to user or not
	IsBroaded bool `gorm:"default:false" json:"is_broaded"`
	EditTime
}

func getContentByFeedItem(source *Source, item *rss.Item) (Content, error) {
	TelegraphURL := ""

	html := item.Content
	if html == "" {
		html = item.Summary
	}

	html = strings.Replace(html, "<![CDATA[", "", -1)
	html = strings.Replace(html, "]]>", "", -1)

	if config.EnableTelegraph && len([]rune(html)) > config.PreviewText {
		TelegraphURL = PublishItem(source, item, html)
	}

	var c = Content{
		Title:        strings.Trim(item.Title, " "),
		Description:  html, //replace all kinds of <br> tag
		SourceID:     source.ID,
		RawID:        item.ID,
		HashID:       genHashID(source.Link, item.ID),
		TelegraphURL: TelegraphURL,
		RawLink:      item.Link,
		IsBroaded:    false,
	}

	return c, nil
}

// GenContentAndCheckByFeedItem generate content by fetcher item
func GenContentAndCheckByFeedItem(s *Source, item *rss.Item) (*Content, bool, error) {
	var (
		content   Content
		isBroaded bool
	)

	hashID := genHashID(s.Link, item.ID)
	db.Where("hash_id=?", hashID).First(&content)
	if content.HashID == "" {
		isBroaded = false
		content, _ = getContentByFeedItem(s, item)
		db.Save(&content)
	} else {
		isBroaded = true
	}

	return &content, isBroaded, nil
}

// DeleteContentsBySourceID delete contents in the db by sourceID
func DeleteContentsBySourceID(sid uint) {
	db.Delete(Content{}, "source_id = ?", sid)
}

// PublishItem publish item to telegraph
func PublishItem(source *Source, item *rss.Item, html string) string {
	url, _ := tgraph.PublishHtml(source.Title, item.Title, item.Link, html)
	return url
}

func QueryUnBroadcastedContent() []Content {
	var contents []Content
	db.Where(map[string]interface{}{"is_broaded": false}).Order("updated_at ASC").Limit(10).Find(&contents)
	return contents
}

func SetIsBroadTrue(content *Content) {
	content.IsBroaded = true
	db.Save(content)
}
