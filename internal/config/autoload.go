package config

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

func init() {
	if isInTests() {
		// 测试环境
		RunMode = TestMode
		initTPL()
		return
	}

	workDirFlag := flag.String("d", "./", "work directory of flowerss")
	configFile := flag.String("c", "", "config file of flowerss")
	printVersionFlag := flag.Bool("v", false, "prints flowerss-bot version")

	testTpl := flag.Bool("testtpl", false, "test template")

	testing.Init()
	flag.Parse()

	if *printVersionFlag {
		// print version
		fmt.Printf(AppVersionInfo())
		os.Exit(0)
	}

	workDir := filepath.Clean(*workDirFlag)

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		viper.SetConfigFile(filepath.Join(workDir, "config.yml"))
	}

	fmt.Println(logo)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	initTPL()
	if *testTpl {
		validateTPL()
		os.Exit(0)
	}

	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")
	UserAgent = viper.GetString("user_agent")

	if viper.IsSet("telegraph_token") {
		EnableTelegraph = true
		TelegraphToken = viper.GetStringSlice("telegraph_token")
	}

	if viper.IsSet("telegraph_account") {
		EnableTelegraph = true
		TelegraphAccountName = viper.GetString("telegraph_account")

		if viper.IsSet("telegraph_author_name") {
			TelegraphAuthorName = viper.GetString("telegraph_author_name")
		}

		if viper.IsSet("telegraph_author_url") {
			TelegraphAuthorURL = viper.GetString("telegraph_author_url")
		}
	}

	if viper.IsSet("preview_text") {
		PreviewText = viper.GetInt("preview_text")
	}

	if viper.IsSet("allowed_users") {
		intAllowUsers := viper.GetStringSlice("allowed_users")
		for _, useIDStr := range intAllowUsers {
			userID, err := strconv.ParseInt(useIDStr, 10, 64)
			if err != nil {
				panic(fmt.Errorf("Fatal error config file: %s", err))
			}
			AllowUsers = append(AllowUsers, userID)
		}
	}

	if viper.IsSet("disable_web_page_preview") {
		DisableWebPagePreview = viper.GetBool("disable_web_page_preview")
	}

	if viper.IsSet("telegram.endpoint") {
		TelegramEndpoint = viper.GetString("telegram.endpoint")
	}

	if viper.IsSet("error_threshold") {
		ErrorThreshold = uint(viper.GetInt("error_threshold"))
	}

	if viper.IsSet("update_interval") {
		UpdateInterval = viper.GetInt("update_interval")
	}

	if viper.IsSet("mysql.host") {
		EnableMysql = true
		Mysql = DBConfig{
			Host:     viper.GetString("mysql.host"),
			Port:     viper.GetInt("mysql.port"),
			User:     viper.GetString("mysql.user"),
			Password: viper.GetString("mysql.password"),
			DB:       viper.GetString("mysql.database"),
		}
	}

	if !EnableMysql {
		if viper.IsSet("postgreSQL.host") {
			EnablePostgreSQL = true
			PostgreSQL = DBConfig{
				Host:     viper.GetString("postgreSQL.host"),
				Port:     viper.GetInt("postgreSQL.port"),
				User:     viper.GetString("postgreSQL.user"),
				Password: viper.GetString("postgreSQL.password"),
				DB:       viper.GetString("postgreSQL.database"),
			}
		}
	}

	if !EnableMysql && !EnablePostgreSQL {
		if viper.IsSet("sqlite.path") {
			SQLitePath = viper.GetString("sqlite.path")
		} else {
			SQLitePath = filepath.Join(workDir, "data.db")
		}
		log.Println("DB Path:", SQLitePath)
		// 判断并创建SQLite目录
		dir := path.Dir(SQLitePath)
		_, err := os.Stat(dir)
		if err != nil {
			err := os.MkdirAll(dir, os.ModeDir)
			if err != nil {
				log.Printf("mkdir failed![%v]\n", err)
			}
		}
	}

	if viper.IsSet("log.db_log") {
		DBLogMode = viper.GetBool("log.db_log")
	}

	if viper.IsSet("transmission.host") {
		Transmission = TransmissionConfig{
			Host:      viper.GetString("transmission.host"),
			Port:      9091,
			Https:     false,
			User:      viper.GetString("transmission.user"),
			Pass:      viper.GetString("transmission.pass"),
			AutoStart: false,
		}

		https := viper.GetBool("transmission.https")
		if https {
			Transmission.Https = true
		}
		port := viper.GetInt32("transmission.port")
		if port != 0 {
			Transmission.Port = port
		}

		autoStart := viper.GetBool("transmission.auto_start")
		if autoStart {
			Transmission.AutoStart = true
		}
	}

	if viper.IsSet("pt_downlaod_dir") {
		PTDownloadDir = PTDownloadDirConfig{
			Default: viper.GetString("pt_downlaod_dir.default"),
			Movie:   viper.GetString("pt_downlaod_dir.default"),
			TV:      viper.GetString("pt_downlaod_dir.default"),
		}
		if viper.IsSet("pt_downlaod_dir.movie") {
			PTDownloadDir.Movie = viper.GetString("pt_downlaod_dir.movie")
		}

		if viper.IsSet("pt_downlaod_dir.tv") {
			PTDownloadDir.TV = viper.GetString("pt_downlaod_dir.tv")
		}
	}

	if viper.IsSet("pt_sites") {
		for site, passkey := range viper.GetStringMapString("pt_sites") {
			PTSites[site] = passkey
		}
	}
}

func (t TplData) Render(mode tb.ParseMode) (string, error) {
	var buf []byte
	wb := bytes.NewBuffer(buf)

	if mode == tb.ModeMarkdown {
		mkd := regexp.MustCompile("(\\[|\\*|\\`|\\_)")
		t.SourceTitle = mkd.ReplaceAllString(t.SourceTitle, "\\$1")
		t.ContentTitle = mkd.ReplaceAllString(t.ContentTitle, "\\$1")
		t.PreviewText = mkd.ReplaceAllString(t.PreviewText, "\\$1")
	} else if mode == tb.ModeHTML {
		t.SourceTitle = t.replaceHTMLTags(t.SourceTitle)
		t.ContentTitle = t.replaceHTMLTags(t.ContentTitle)
		t.PreviewText = t.replaceHTMLTags(t.PreviewText)
	}

	if err := MessageTpl.Execute(wb, t); err != nil {
		return "", err
	}

	return strings.TrimSpace(string(wb.Bytes())), nil
}

func (t TplData) replaceHTMLTags(s string) string {

	rStr := strings.ReplaceAll(s, "&", "&amp;")
	rStr = strings.ReplaceAll(rStr, "\"", "&quot;")
	rStr = strings.ReplaceAll(rStr, "<", "&lt;")
	rStr = strings.ReplaceAll(rStr, ">", "&gt;")

	return rStr
}

func validateTPL() {
	testData := []TplData{
		TplData{
			"RSS 源标识 - 无预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"",
			"",
			"",
			false,
		},
		TplData{
			"RSS源标识 - 有预览无telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁[1](123)",
			"",
			"#标签",
			false,
		},
		TplData{
			"RSS源标识 - 有预览有telegraph的消息",
			"这是标题",
			"https://www.github.com/",
			"这里是很长很长很长的消息预览字数补丁紫薯补丁紫薯补丁紫薯补丁紫薯补丁",
			"https://telegra.ph/markdown-07-07",
			"#标签1 #标签2",
			true,
		},
	}

	for _, d := range testData {
		fmt.Println("\n////////////////////////////////////////////")
		fmt.Println(d.Render(MessageMode))
	}
	fmt.Println("\n////////////////////////////////////////////")
}

func initTPL() {
	var tplMsg string
	if viper.IsSet("message_tpl") {
		tplMsg = viper.GetString("message_tpl")
	} else {
		tplMsg = defaultMessageTpl
	}
	MessageTpl = template.Must(template.New("message").Parse(tplMsg))

	if viper.IsSet("message_mode") {
		switch strings.ToLower(viper.GetString("message_mode")) {
		case "md", "markdown":
			MessageMode = tb.ModeMarkdown
		case "html":
			MessageMode = tb.ModeHTML
		default:
			MessageMode = tb.ModeDefault
		}
	} else {
		MessageMode = defaultMessageTplMode
	}
}

func getInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func (m *DBConfig) GetMySQLConnectingString() string {
	usr := m.User
	pwd := m.Password
	host := m.Host
	port := m.Port
	db := m.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", usr, pwd, host, port, db)
}

func (m *DBConfig) GetPostgreSQLConnectingString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		m.Host,
		m.User,
		m.Password,
		m.DB,
		m.Port)
}

func isInTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			if arg == "-testtpl" {
				continue
			}
			return true
		}
	}
	return false
}
