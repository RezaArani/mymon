package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Headers []string

// performance Metrics type
type QuickTestConfig struct {
	Url                string `json:"url" cql:"url"`
	Interval           int    `json:"interval" cql:"interval"`
	AlarmFailTestCount int    `json:"alarmfailtestcount" cql:"alarmfailtestcount"`
}

type Config struct {
	HTTPBinding string             `json:"HTTPBinding"`
	WebhookHttpMethod      string            `json:"WebhookHttpMethod"`
	PostBody        string            `json:"PostBody"`
	FollowRedirects bool              `json:"FollowRedirects"`
	HttpHeaders     Headers           `json:"HttpHeaders"`
	ClientCertFile  string            `json:"ClientCertFile"`
	FourOnly        bool              `json:"FourOnly"`
	SixOnly         bool              `json:"SixOnly"`
	MaxRedirects    int               `json:"MaxRedirects"`
	Debug           bool              `json:"Debug"`
	WebhookURL      string            `json:"WebhookURL"`
	Websites        []QuickTestConfig `json:"Websites"`
	folderSep       string
}

var curConfig Config

func (h Headers) String() string {
	var o []string
	for _, v := range h {
		o = append(o, "-H "+v)
	}
	return strings.Join(o, " ")
}

func (h *Headers) Set(v string) error {
	*h = append(*h, v)
	return nil
}

func (h Headers) Len() int      { return len(h) }
func (h Headers) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h Headers) Less(i, j int) bool {
	a, b := h[i], h[j]
	if a == "Server" {
		return true
	}
	if b == "Server" {
		return false
	}
	endtoend := func(n string) bool {
		switch n {
		case "Connection",
			"Keep-Alive",
			"Proxy-Authenticate",
			"Proxy-Authorization",
			"TE",
			"Trailers",
			"Transfer-Encoding",
			"Upgrade":
			return false
		default:
			return true
		}
	}

	x, y := endtoend(a), endtoend(b)
	if x == y {
		// both are of the same class
		return a < b
	}
	return x
}

func (appConfig *Config) InitConfig() {
	ConfigFilePath := ""

	if len(os.Args) > 1 {
		ConfigFilePath = os.Args[1]
	} else {
		if runtime.GOOS != "windows" {
			appConfig.folderSep = "/"
		} else {
			appConfig.folderSep = "\\"
		}
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		ConfigFilePath = exPath + appConfig.folderSep
	}

	configFile, err := os.Open(ConfigFilePath + "config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
	if appConfig.WebhookHttpMethod == "" {
		appConfig.WebhookHttpMethod = "GET"
	}
	if appConfig.MaxRedirects == 0 {
		appConfig.MaxRedirects = 10
	}
	curConfig = *appConfig
}

func GetConfig() Config {
	return curConfig
}
