package wlogs

import (
	"log"

	"github.com/RezaArani/mymon/config"
)
 

func LogMessage(sender interface{},force bool){
	if config.GetConfig().Debug||force {
		log.Println(sender)
	}
}