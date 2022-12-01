package models

type Configuration struct {
	WirelessControllers []WirelessController `json:"wireless-controllers"`
	APTagMaps           APEthernetMAC        `json:"ap-tag-map"`
	MQTTConfig          MQTTConfig           `json:"mqtt"`
}

type WirelessController struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type APEthernetMAC map[string]APTagMap

type APTagMap struct {
	SiteTag   string `json:"site-tag"`
	PolicyTag string `json:"policy-tag"`
	RFTag     string `json:"rf-tag"`
}

type MQTTConfig struct {
	Broker   string `json:"broker"`
	Port     int32  `json:"port"`
	ClientId string `json:"client-id"`
	Topic    string `json:"topic"`
}
