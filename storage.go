package main

import "time"

type pingedHosts struct {
	Host         string
	Time         time.Time
	ResponseTime time.Duration
}

var storage = make(map[string]pingedHosts)

func addHost(host string, responseTime time.Duration) {
	currentTime := time.Now()

	storage[host] = pingedHosts{
		Host:         host,
		Time:         currentTime,
		ResponseTime: responseTime,
	}

	time.AfterFunc(time.Duration(*checkingInterval)*time.Second, func() { removeHost(host) })
}

func isRecentlyPinged(host string) bool {
	_, found := storage[host]
	return found
}

func removeHost(host string) {
	delete(storage, host)
}
