package main

import (
	"log"
	api_mqtt "mantap2/api/mqtt"
	"mantap2/config"
	uhfrf_tcp "mantap2/nfcreader"
)

func wrLog(isdevonly bool, msg string) {
	if config.Log_show {
		if isdevonly {
			if config.Log_dev {
				log.Print(msg)
			}
		} else {
			log.Print(msg)
		}
	}
}

func main() {

	//arguments := os.Args

	api_mqtt.JT_mqtt_run()

	uhfrf_tcp.DoJobs()
}
