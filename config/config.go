package config

import (
	"time"

	"github.com/spf13/viper"
)

var Vr_gen_nama string = ""
var Vr_gen_devicetipe string = ""

var Log_show bool = false
var Log_dev bool = false

var Rdr_ip string = ""
var Rdr_port string = ""
var Rdr_timeout time.Duration
var Api_server string = ""
var IsUdp bool = false
var IsListen bool = false

var Vr_mqtt_server string = ""
var Vr_mqtt_port int = 1883
var Vr_mqtt_clientID string = ""
var Vr_mqtt_username string = ""
var Vr_mqtt_password string = ""

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	viper.ReadInConfig()

	Vr_gen_nama = viper.GetString("general.name")
	Vr_gen_devicetipe = viper.GetString("general.device_tipe")

	Log_show = viper.GetBool("log.verbose")
	Log_dev = viper.GetBool("log.dev")

	Vr_mqtt_server = viper.GetString("mqtt.server")
	Vr_mqtt_port = viper.GetInt("mqtt.port")
	Vr_mqtt_clientID = "jtmqtt_" + Vr_gen_nama
	Vr_mqtt_username = viper.GetString("mqtt.user")
	Vr_mqtt_password = viper.GetString("mqtt.pass")

	Rdr_ip = viper.GetString("target_reader.ip")
	Rdr_port = viper.GetString("target_reader.port")
	Rdr_timeout = viper.GetDuration("target_reader.timeout") * time.Second
	IsUdp = viper.GetBool("target_reader.is_udp")
	//var sourcePort string
	IsListen = false

	Api_server = viper.GetString("api.server")

}
