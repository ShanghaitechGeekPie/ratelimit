package main

import "net"

type RawServerConfig struct {
	Listen string               `json:"listen"`
	Policy []RawRateLimitPolicy `json:"policy"`
}

type ServerConfig struct {
	Listen string            `json:"listen"`
	Policy []RateLimitPolicy `json:"policy"`
}

type RawRateLimitPolicy struct {
	Bandwidth   int    `json:"bandwidth"`
	Burst       int    `json:"burst"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type RateLimitPolicy struct {
	Bandwidth   int64
	Burst       int64
	Source      *net.IPNet
	Destination string
}
