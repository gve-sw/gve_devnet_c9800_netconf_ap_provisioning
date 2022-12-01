package models

import "encoding/xml"

// Assign Site, Policy, or RF tags to AP by Ethernet MAC
type ApConfig struct {
	XMLName   xml.Name  `xml:"config"`
	ApCFGData ApCFGData `xml:"ap-cfg-data"`
}

type ApCFGData struct {
	ApTags  ApTags   `xml:"ap-tags"`
	Xmlns   string   `xml:"xmlns,attr"`
	XMLName xml.Name `xml:"ap-cfg-data"`
}

type ApTags struct {
	ApTag ApTag `xml:"ap-tag"`
}

type ApTag struct {
	ApMAC     string `xml:"ap-mac"`
	PolicyTag string `xml:"policy-tag"`
	SiteTag   string `xml:"site-tag"`
	RFTag     string `xml:"rf-tag"`
}

// RPC reply when saving config
type SaveConfig struct {
	XMLName   xml.Name `xml:"rpc-reply"`
	Xmlns     string   `xml:"xmlns,attr"`
	MessageID string   `xml:"message-id,attr"`
	Result    string   `xml:"result"`
}
