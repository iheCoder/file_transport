package net_helper

import "testing"

func TestDetectLocalIPs(t *testing.T) {
	ips, err := detectLocalNetIPs()
	if err != nil {
		t.Fatal(err)
	}

	for _, ip := range ips {
		t.Log(ip)
	}
}
