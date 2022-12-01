/*
Copyright (c) 2022 Cisco and/or its affiliates.
This software is licensed to you under the terms of the Cisco Sample
Code License, Version 1.1 (the "License"). You may obtain a copy of the
License at

	https://developer.cisco.com/docs/licenses

All use of the material herein must be in accordance with the terms of
the License. All rights not expressly granted by the License are
reserved. Unless required by applicable law or agreed to separately in
writing, software distributed under the License is distributed on an "AS
IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied.
*/
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gve-sw/gve_devnet_c9800_netconf_ap_provisioning/models"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/opoptions"
	ncopt "github.com/scrapli/scrapligo/driver/options"
)

var WLC_user string
var WLC_pass string
var Config models.Configuration

func main() {
	// Open configuration file
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	// Read in config data
	configData, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}

	err = json.Unmarshal(configData, &Config)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Look up WLC User & Password environmental variables
	var ok bool
	WLC_user, ok = os.LookupEnv("WLC_USER")
	if !ok {
		log.Fatal("Please set environment variable: WLC_USER")
	}
	WLC_pass, ok = os.LookupEnv("WLC_PASSWORD")
	if !ok {
		log.Fatal("Please set environment variable: WLC_PASSWORD")
	}

	// Start MQTT Subscription
	err = subscribeMQTT()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for termination signal
	log.Println("Listener running. Press Ctrl-C to quit.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Quitting...")
}

// Get MQTT Messages
func subscribeMQTT() error {
	// MQTT client configuration
	var broker = Config.MQTTConfig.Broker
	var port = Config.MQTTConfig.Port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(Config.MQTTConfig.ClientId)

	// Connect to MQTT broker
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if subscription := client.Subscribe(Config.MQTTConfig.Topic, 0, provisionAP); subscription.Wait() && subscription.Error() != nil {
		log.Printf("Failed to subscribe to %v", broker)
	}
	log.Printf(">> Subscribed to MQTT at %v\n", broker)
	return nil

}

// Take in WLC/MAC map from MQTT & send AP tags to WLC via NETCONF
func provisionAP(client mqtt.Client, msg mqtt.Message) {
	// Parse incoming MQTT message
	var message models.ApMQTTMessage
	json.Unmarshal(msg.Payload(), &message)

	log.Printf(">> Received new AP MAC: %s on WLC at %s", message.MAC, message.WLC)

	// Check AP MAC has mapping defined
	if _, ok := Config.APTagMaps[message.MAC]; !ok {
		log.Printf("No config map for AP MAC: %s\n", message.MAC)
		return
	}

	// Assemble NETCONF config payload
	newTags := models.ApTag{
		ApMAC:     message.MAC,
		PolicyTag: Config.APTagMaps[message.MAC].PolicyTag,
		SiteTag:   Config.APTagMaps[message.MAC].SiteTag,
		RFTag:     Config.APTagMaps[message.MAC].RFTag,
	}

	newConfig := models.ApConfig{}
	newConfig.ApCFGData.Xmlns = "http://cisco.com/ns/yang/Cisco-IOS-XE-wireless-ap-cfg"
	newConfig.ApCFGData.ApTags.ApTag = newTags
	// Build XML payload
	configXML, _ := xml.Marshal(newConfig)

	// Get WLC config
	var netconfAddress string
	var netconfPort int
	for _, wlc := range Config.WirelessControllers {
		if wlc.Name == message.WLC {
			netconfAddress = wlc.Name
			netconfPort = wlc.Port
			break
		}
	}
	if netconfAddress == "" {
		log.Printf("Error: WLC at %s not found in config. AP at MAC %s will not be provisioned\n", message.WLC, message.MAC)
		return
	}

	// Create NETCONF session
	ncSession := createSession(netconfAddress, netconfPort)
	defer log.Printf(">> Session to %s closed.", message.WLC)
	defer ncSession.Close()

	// Send new config payload
	response, err := ncSession.EditConfig("running", string(configXML))
	if err != nil {
		log.Printf("Error when provisioning for MAC %s: %v", message.MAC, err.Error())
	}
	if response.Failed != nil {
		log.Println(response.Failed)
	}
	// Check RPC Reply
	if strings.Contains(response.Result, "ok") {
		log.Printf(">> Provisioning request sent for MAC: %s\n", message.MAC)
	} else {
		log.Printf("Error with provisioning request for MAC %s: %v\n", message.MAC, response.Result)
	}

	// Commit WLC config changes
	log.Printf("Saving WLC config...\n")
	response, err = ncSession.RPC(opoptions.WithFilter("<cisco-ia:save-config xmlns:cisco-ia=\"http://cisco.com/yang/cisco-ia\"/> "))
	if err != nil {
		log.Printf("Error trying to save config on WLC %s: %v", message.WLC, err)
	}
	saveResult := models.SaveConfig{}
	xml.Unmarshal([]byte(response.Result), &saveResult)
	log.Printf("Response from %s: %v\n", message.WLC, saveResult.Result)
}

// Creates NETCONF session to target device
func createSession(address string, port int) *netconf.Driver {
	// Create new NETCONF driver with target device info
	ncSession, _ := netconf.NewDriver(
		address,
		ncopt.WithAuthNoStrictKey(),
		ncopt.WithAuthUsername(WLC_user),
		ncopt.WithAuthPassword(WLC_pass),
		ncopt.WithPort(830),
	)
	// Open NETCONF session
	err := ncSession.Open()
	if err != nil {
		log.Printf("Failed to connect to WLC at %v: %v\n", address, err.Error())
		return nil
	}
	log.Printf(">> Established NETCONF session to %s\n", address)

	return ncSession
}
