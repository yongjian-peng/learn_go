package main

import "testing"

// TestProxy StationProxy 代理了 Station 代理类中持有被代理类对象，并且和被代理类对象实现了同一个接口
func TestProxy(t *testing.T) {
	station := &Station{stock: 100}

	station.sell("中心火车站")
	station.sell("中心火车站")
	station.sell("中心火车站")

	stationProxy := &StationProxy{station: station}

	stationProxy.sell("代理火车站")
	stationProxy.sell("代理火车站")
	station.sell("中心火车站")
}
