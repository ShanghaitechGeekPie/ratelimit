package main

import (
	"net"
	"time"
)

const bufferSize = 8192

func Init(rawConfig *RawServerConfig) *ServerConfig {
	config := &ServerConfig{
		Listen: rawConfig.Listen,
		Policy: []RateLimitPolicy{},
	}
	for _, v := range rawConfig.Policy {
		_, IPNet, err := net.ParseCIDR(v.Source)
		if err != nil {
			panic("unresolved CIDR " + v.Source)
		}
		config.Policy = append(config.Policy, RateLimitPolicy{
			Bandwidth:   int64(v.Bandwidth),
			Burst:       int64(v.Burst),
			Source:      IPNet,
			Destination: v.Destination,
		})
	}
	initBucket(config)
	return config
}

func initBucket(config *ServerConfig) {
	for _, v := range config.Policy {
		BucketList = append(BucketList, NewBucketWithQuantum(time.Second, v.Bandwidth*Kilobyte.GetByte(), v.Burst*Kilobyte.GetByte()))
	}
}
