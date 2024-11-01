package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"time"

	"github.com/pkoehlers/mqtt-wakeonlan/config"
	"github.com/pkoehlers/mqtt-wakeonlan/wakeonlan"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	stateTopic string
	nameTopic  string
	macTopic   string
)

var messageNameTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	mac, err := config.LookupMacAddressByName(string(msg.Payload()))
	checkAndHandleErrorWithMqtt(err, client)
	wakeonlan.SendDefaultWOL(mac)
}

var messageMacTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	wakeonlan.SendDefaultWOL(string(msg.Payload()))
}
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	publishState(client, "idle")

	subscriptions := []struct {
		Topic   string
		Handler mqtt.MessageHandler
	}{
		{nameTopic, messageNameTopicHandler},
		{macTopic, messageMacTopicHandler},
	}
	for _, sub := range subscriptions {
		token := client.Subscribe(sub.Topic, 1, sub.Handler)
		token.Wait()
		log.Printf("Subscribed to topic: %s", sub.Topic)
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func initMqttTopics() {
	topicPrefix := config.MqttTopicPrefix()
	stateTopic = topicPrefix + "/status"
	nameTopic = topicPrefix + "/name"
	macTopic = topicPrefix + "/mac"
}

func main() {
	var broker = config.MqttHost()
	var port = config.MqttPort()
	var protocol = "tcp"
	initMqttTopics()
	opts := mqtt.NewClientOptions()
	if config.MqttTLSEnabled() {
		protocol = "ssl"
		if len(config.MqttTLSCA()) > 0 {
			opts.TLSConfig = new(tls.Config)
			opts.TLSConfig.InsecureSkipVerify = false
			opts.TLSConfig.RootCAs = x509.NewCertPool()
			opts.TLSConfig.RootCAs.AppendCertsFromPEM([]byte(config.MqttTLSCA()))
		}
	}
	url := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	log.Printf("Connecting to %s", url)
	opts.AddBroker(url)

	opts.SetClientID(config.MqttIdentifier())
	opts.SetUsername(config.MqttUsername())
	opts.SetPassword(config.MqttPassword())
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	opts.WillTopic = stateTopic
	opts.WillEnabled = true
	opts.WillPayload = []byte("offline")
	opts.WillRetained = true

	client := mqtt.NewClient(opts)
	retry := time.NewTicker(5 * time.Second)
	for range retry.C {
		token := client.Connect()
		token.Wait()
		error := token.Error()
		if error != nil {
			log.Printf("MQTT connection failed: %v\n", error)
		} else {
			retry.Stop()
			break
		}
	}
	for {
		time.Sleep(time.Second)
	}
}

func publishState(client mqtt.Client, status string) {
	token := client.Publish(stateTopic, 0, true, status)
	token.Wait()
	time.Sleep(time.Second)
}

func checkAndHandleErrorWithMqtt(err error, client mqtt.Client) {
	if err != nil {
		log.Printf("Error occured: %v", err)
		publishState(client, "error")
	}
}
