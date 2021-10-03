package bot

import (
	"testing"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestGetTorrentDownloadUrl(t *testing.T) {
	config.PTSites["kp.m-team.cc"] = "abcd"
	url := "https://kp.m-team.cc/details.php?id=514703"
	actURL := getTorrentDownloadUrl(url)
	assert.Equal(t, "https://kp.m-team.cc/download.php?id=514703&passkey=abcd", actURL)

	url = "https://kp.m-team.cc/download.php?id=514703&passkey=1234"
	actURL = getTorrentDownloadUrl(url)
	assert.Equal(t, "https://kp.m-team.cc/download.php?id=514703&passkey=1234", actURL)

	url = "https://kp.m-team.cc/details.php"
	actURL = getTorrentDownloadUrl(url)
	assert.Equal(t, "https://kp.m-team.cc/details.php", actURL)
}

func TestShouldParseAsPTSite(t *testing.T) {
	config.PTSites["kp.m-team.cc"] = "abcd"
	config.PTDownloadDir.Default = "download"
	urlStr := "https://kp.m-team.cc/details.php"
	actRes := shouldParseAsPTSite(urlStr)
	assert.False(t, actRes)

	config.EnableTransmission = true
	actRes = shouldParseAsPTSite(urlStr)
	assert.True(t, actRes)
}

func TestBuildReplyMarkupForPTSite(t *testing.T) {
	urlStr := "https://kp.m-team.cc/download.php?id=514703&passkey=1234"
	exp := &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{
				{
					Unique: "download_to_movies",
					Text:   "Movies",
					Data:   urlStr,
				},
				{
					Unique: "download_to_tvs",
					Text:   "TVs",
					Data:   urlStr,
				},
				{
					Unique: "download",
					Text:   "Download",
					Data:   urlStr,
				},
			},
		},
	}

	act := buildReplyMarkupForPTSite(urlStr)

	assert.Equal(t, exp, act)
}
