package scheduler

import (
	"log"
	"time"

	"github.com/RezaArani/mymon/config"
	"github.com/RezaArani/mymon/webhook"
	"github.com/RezaArani/mymon/webmonitor"
	"github.com/RezaArani/mymon/wlogs"
)


func AddToSchedule(qt *webmonitor.QuickTest) {
	if qt.Interval == 0 {
		qt.Interval = 5
	}

	if qt.AlarmFailTestCount == 0 {
		qt.Interval = 12
	}

	log.Println(qt.Url + ", Monitoring started.")
	qt.IgnoreCertificate = true
	for range time.Tick(time.Second * time.Duration(qt.Interval)) {
		go func(qt *webmonitor.QuickTest) {
			WebhookShouldCalled := false
			wlogs.LogMessage("Testing:"+qt.Url, false)
			lastWebhookState := qt.LastWebhookState
			url := webmonitor.ParseURL(qt.Url)
			qt.Visit(url, nil)

			wlogs.LogMessage(qt.StatusCode, false)
			 
			if !qt.IsHealthy {
				qt.LastError = time.Now()
				if qt.ActiveFailures <= qt.AlarmFailTestCount {
					//prevent int overflow
					qt.ActiveFailures++
				}
			} else {
				if qt.ActiveFailures >= qt.AlarmFailTestCount {
					//restore webhook
					qt.LastWebhookState = 1
					qt.ErrorText="Website has been restored."
					WebhookShouldCalled = true
					wlogs.LogMessage("website restored", false)

				}
				qt.ActiveFailures = 0
			}
			qt.TestCount++
			if !qt.IsHealthy {
				if qt.TestCount > 1  && qt.ActiveFailures == qt.AlarmFailTestCount {
					//failure webhook
					qt.LastWebhookState = -1
					WebhookShouldCalled = true
					wlogs.LogMessage("website load failed", false)

				}
			}
			if WebhookShouldCalled{
				if lastWebhookState!=qt.LastWebhookState{
					var wh webhook.Webhook
					wh.Url = config.GetConfig().WebhookURL
					wh.QT = *qt
					wh.Method = config.GetConfig().WebhookHttpMethod
					wh.CallWebhook()
				}
			}
		}(qt)
	}
}
