package bot

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/transmission"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func downloadToMoviesCtr(c *tb.Callback) {
	downloadBT(c, config.PTDownloadDir.Movie)
}

func downloadToTVsCtr(c *tb.Callback) {
	downloadBT(c, config.PTDownloadDir.Movie)
}

func downloadCtr(c *tb.Callback) {
	downloadBT(c, config.PTDownloadDir.Default)
}

func downloadBT(c *tb.Callback, downloadDir string) {
	torrent, err := transmission.AddTorrent(
		getTorrentDownloadUrl(c.Data),
		downloadDir,
	)
	if err != nil {
		_, _ = B.Send(c.Message.Chat, fmt.Sprintf("Add torrent Error\n\nError Message:\n%#v", err), &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
		return
	}

	_, _ = B.Send(c.Message.Chat, fmt.Sprintf("Add torrent Success\n\n*Torrent Name:*\n%s", *torrent.Name), &tb.SendOptions{
		DisableWebPagePreview: true,
		ParseMode:             tb.ModeMarkdown,
	})
}

func getTorrentDownloadUrl(urlStr string) string {
	// only supported site will get this far, so only need check or set passkey for the url
	u, err := url.Parse(urlStr)
	if err != nil {
		zap.S().Errorw(
			"Parse download url error",
			"error", err,
		)
		return urlStr
	}
	// TODO for now only support download use passkey, will support more in the future
	// check passkey
	if strings.Contains(urlStr, "passkey") {
		return urlStr
	}

	// set passkey
	passkey := config.PTSites[u.Host]
	return urlStr + "&passkey=" + passkey
}

func shouldParseAsPTSite(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err == nil {
		_, ok := config.PTSites[u.Host]
		if ok {
			if config.PTDownloadDir.Default != "" && config.EnableTransmission {
				return true
			}
		}
	}
	return false
}

func buildReplyMarkupForPTSite(urlStr string) *tb.ReplyMarkup {
	return &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{
			{
				{
					Unique: "download_to_movies",
					Text:   "Download Movies",
					Data:   urlStr,
				},
				{
					Unique: "download_to_tvs",
					Text:   "Download TVs",
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
}
