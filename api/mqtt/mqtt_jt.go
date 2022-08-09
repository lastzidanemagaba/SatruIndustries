package apimqtt

import (
	"fmt"
	"log"
	"mantap2/config"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	wrLog(true, fmt.Sprintf("MQTT Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	wrLog(false, "MQTT Connected to "+config.Vr_mqtt_server+" PORT "+strconv.Itoa(config.Vr_mqtt_port))
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	wrLog(false, fmt.Sprintf("MQTT Connect lost: %v", err))
}

func JT_mqtt_run() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Vr_mqtt_server, config.Vr_mqtt_port))
	opts.SetClientID(config.Vr_mqtt_clientID)
	opts.SetUsername(config.Vr_mqtt_username)
	opts.SetPassword(config.Vr_mqtt_password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
