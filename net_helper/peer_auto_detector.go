package net_helper

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

func detectLocalNetIPs() ([]net.IP, error) {
	cmd := exec.Command("arp-scan", "--interface=en0", "--localnet")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// parse the output
	return parseArpScanOutput(string(out))
}

func parseArpScanOutput(output string) ([]net.IP, error) {
	var ips []string
	scanner := bufio.NewScanner(strings.NewReader(output))

	// 使用正则表达式匹配 IP 地址
	ipRegex := regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}`)

	for scanner.Scan() {
		line := scanner.Text()

		// 查找 IP 地址
		matches := ipRegex.FindString(line)
		if matches != "" {
			ips = append(ips, matches)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 将字符串转换为 net.IP 类型
	var netIPs []net.IP
	for _, ip := range ips {
		netIP := net.ParseIP(ip)
		if netIP == nil {
			return nil, fmt.Errorf("invalid IP address: %s", ip)
		}
		netIPs = append(netIPs, netIP)
	}

	return netIPs, nil
}
