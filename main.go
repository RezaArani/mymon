package main

import (
	"github.com/RezaArani/mymon/config"
	"github.com/RezaArani/mymon/scheduler"
	"github.com/RezaArani/mymon/webmonitor"
	"github.com/RezaArani/mymon/webserver"
)

func main() {
	var appConfig config.Config
	appConfig.InitConfig()
	

	for _,webToMonitor:= range appConfig.Websites{
		var qt webmonitor.QuickTest
		qt.Url = webToMonitor.Url
		qt.Interval = webToMonitor.Interval
		qt.AlarmFailTestCount = webToMonitor.AlarmFailTestCount
		go scheduler.AddToSchedule(&qt)
	}
	
	webserver.Init(appConfig.HTTPBinding)
 
}
