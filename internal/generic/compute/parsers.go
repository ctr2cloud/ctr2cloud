package compute

import (
	"net"
	"regexp"
	"strings"

	public "github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
)

func ParseIPAddrOutput(output string) []public.Address {
	var ipAddresses []public.Address
	lines := strings.Split(output, "\n")

	reIPv4 := regexp.MustCompile(`\s+inet (\d+\.\d+\.\d+\.\d+)/(\d+)`)
	reIPv6 := regexp.MustCompile(`\s+inet6 ([a-fA-F0-9:]+)/(\d+)`)

	for _, line := range lines {
		// Match IPv4 address
		if match := reIPv4.FindStringSubmatch(line); match != nil {
			ip := match[1]
			netmask := match[2]
			isRoutable := isRoutableIPv4(ip)

			ipAddresses = append(ipAddresses, public.Address{
				Address:            ip,
				Type:               "IPv4",
				Netmask:            netmask,
				IsPubliclyRoutable: isRoutable,
			})
		}

		// Match IPv6 address
		if match := reIPv6.FindStringSubmatch(line); match != nil {
			ip := match[1]
			netmask := match[2]
			isRoutable := isRoutableIPv6(ip)

			ipAddresses = append(ipAddresses, public.Address{
				Address:            ip,
				Type:               "IPv6",
				Netmask:            netmask,
				IsPubliclyRoutable: isRoutable,
			})
		}
	}

	return ipAddresses
}

func isRoutableIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check for private ranges (RFC 1918)
	privateIPv4Ranges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateIPv4Ranges {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet.Contains(parsedIP) {
			return false
		}
	}

	return true
}

func isRoutableIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check for private IPv6 range (unique local address, fc00::/7)
	_, ulaSubnet, _ := net.ParseCIDR("fc00::/7")
	return !ulaSubnet.Contains(parsedIP)
}
