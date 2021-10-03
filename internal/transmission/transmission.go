package transmission

import (
	"github.com/hekmon/transmissionrpc"
	"github.com/indes/flowerss-bot/internal/config"
	"go.uber.org/zap"
)

var client transmissionrpc.Client

func init() {
	transmissionbt, err := transmissionrpc.New(
		config.Transmission.Host,
		config.Transmission.User,
		config.Transmission.Pass,
		&transmissionrpc.AdvancedConfig{
			HTTPS: config.Transmission.Https,
			Port:  uint16(config.Transmission.Port),
		})
	if err != nil {
		zap.S().Error(
			"Init transmission Client error",
			"host", config.Transmission.Host,
			"port", config.Transmission.Port,
			"user", config.Transmission.User,
			"https", config.Transmission.Https,
			"error", err,
		)
	} else {
		zap.S().Infow(
			"Init transmission Client success",
			"host", config.Transmission.Host,
			"port", config.Transmission.Port,
			"user", config.Transmission.User,
			"https", config.Transmission.Https,
		)
		client = *transmissionbt
		checkVersion()
	}
}

func checkVersion() {
	ok, serverVersion, serverMinimumVersion, err := client.RPCVersion()
	if err != nil {
		zap.S().Warnw(
			"Init transmission error, set enable EnableTransmission to false",
			"error", err,
		)
		return
	}
	if !ok {
		zap.S().Errorf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
	}
	zap.S().Infof("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
		serverVersion, transmissionrpc.RPCVersion)
	config.EnableTransmission = true
}

// AddTorrent AddTorrent to remote transmission rpc
func AddTorrent(url, downloadDir string) (*transmissionrpc.Torrent, error) {
	paused := true
	torrent, err := client.TorrentAdd(
		&transmissionrpc.TorrentAddPayload{
			DownloadDir: &downloadDir,
			Filename:    &url,
			Paused:      &paused,
		},
	)
	if err != nil {
		zap.S().Errorw(
			"AddTorrent error",
			"url", url,
			"downloadDir", downloadDir,
			"error", err,
		)
		return nil, err
	}
	zap.S().Infow(
		"AddTorrent Success",
		"url", url,
		"downloadDir", downloadDir,
		"Name", torrent.Name,
		"ID", torrent.ID,
	)
	return torrent, nil
}

func GetTorrent(fields []string, id int64) (*transmissionrpc.Torrent, error) {
	torrents, err := client.TorrentGet(fields, []int64{id})
	if err != nil {
		return nil, err
	}
	return torrents[0], nil
}
