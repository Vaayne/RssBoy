package bot

import (
	"fmt"
	"net/url"

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
		_, _ = B.Send(c.Message.Chat, fmt.Sprintf("Download torrent Error\nError Message:\n%#v", err), &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
		return
	}

	t, err := transmission.GetTorrent([]string{"id", "name", "totalSize", "downloadDir", "status"}, *torrent.ID)

	if err != nil {
		_, _ = B.Send(c.Message.Chat, fmt.Sprintf("Get torrent Error\nError Message:\n%#v", err), &tb.SendOptions{
			DisableWebPagePreview: true,
			ParseMode:             tb.ModeMarkdown,
		})
		return
	}

	msg := "Download torrent Success:\n- *ID:* %d\n- *Name:* %s\n- *DownloadDir:* %s\n- *TotalSize:* %f GB\n- *Status:* %s"
	_, _ = B.Send(
		c.Message.Chat,
		fmt.Sprintf(msg, *t.ID, *t.Name, *t.DownloadDir, t.TotalSize.GB(), t.Status.String()),
		&tb.SendOptions{
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
	// build download url
	id := u.Query().Get("id")

	// if there is not id, it is an invalid url, just return
	if id == "" {
		zap.S().Errorw(
			"Invalid downlaod url",
			"url", urlStr,
		)
		return urlStr
	}

	passkey := u.Query().Get("passkey")

	// if is download.php and contains id and passkey
	// it is a downlaod url, just return
	if u.Path == "/download.php" {
		if passkey != "" {
			zap.S().Infow(
				"Use the download url",
				"url", urlStr,
			)
			return urlStr
		}
	}

	passkey = config.PTSites[u.Host]

	if passkey == "" {
		zap.S().Errorw(
			"There is no valid passkey, can not download",
			"url", urlStr,
		)
		return urlStr
	}

	// build new download url
	params := url.Values{
		"id":      []string{id},
		"passkey": []string{passkey},
	}
	newURL := &url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     "/download.php",
		RawQuery: params.Encode(),
	}
	zap.S().Infow(
		"Parse download url",
		"old_url", urlStr,
		"new_url", newURL.String(),
	)
	return newURL.String()
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
}
