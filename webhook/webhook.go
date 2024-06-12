package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/RezaArani/mymon/config"
	"github.com/RezaArani/mymon/webmonitor"
	"github.com/RezaArani/mymon/wlogs"
)

type Webhook struct {
	Url    string
	Method string
	// ExpectedResult string // future implementation
	QT webmonitor.QuickTest
}

func (w Webhook) CallWebhook() error {
	w.Url = config.GetConfig().WebhookURL
	w.Method = strings.ToUpper(w.Method)
	if w.Method == "" {
		w.Method = "GET"
	} else {
		if w.Method != "POST" && w.Method != "GET" {
			return errors.New("invalid service method call")
		}
	}
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(w.QT)
	testBytes, _ := json.Marshal(w.QT)
	urlToCall := strings.ReplaceAll(w.Url, "~TESTINFO~", url.QueryEscape(string(testBytes)))
	urlToCall = strings.ReplaceAll(urlToCall, "~URL~", url.QueryEscape(w.QT.Url))
	urlToCall = strings.ReplaceAll(urlToCall, "~ERROR~", url.QueryEscape(w.QT.ErrorText))
	req, _ := http.NewRequest(w.Method, urlToCall, &buf)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		wlogs.LogMessage(err, true)
	}
	return err
}
