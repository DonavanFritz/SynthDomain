package synthdomain

import (
	"encoding/hex"
	"net"
	"strings"
)

func nameToIp(name string) net.IP {
	if !strings.HasPrefix(name, "ip-") {
		return nil
	}

	name = strings.TrimPrefix(name, "ip-")
	name = strings.Split(name, ".")[0]
	name = strings.ReplaceAll(name, "-", ".")
	ip := net.ParseIP(name)

	if ip == nil {
		name = strings.ReplaceAll(name, ".", ":")
		ip = net.ParseIP(name)
	}
	return ip
}

func ipToName(ip net.IP, zone string) string {
	if ip == nil {
		return ""
	}

	if !strings.HasPrefix(zone, ".") {
		zone = "." + zone
	}

	if !strings.HasSuffix(zone, ".") {
		zone = zone + "."
	}

	sep := ":"
	if ip.To4() != nil {
		sep = "."
	}

	response := strings.Join(strings.Split(ip.String(), sep), "-")
	return "ip-" + response + zone
}

func arpaToIp(name string) net.IP {
	ipv4Suffix := ".in-addr.arpa."
	ipv6Suffix := ".ip6.arpa."

	if idx := strings.Index(name, ipv4Suffix); idx > 6 {
		name = name[:idx]
		parts := strings.Split(name, ".")
		if len(parts) != 4 {
			return nil
		}

		name = parts[3] + "." + parts[2] + "." + parts[1] + "." + parts[0]
		return net.ParseIP(name)
	}

	if len(name) == 73 && name[63:] == ipv6Suffix {
		// we can rely on the fact that v6 reverse hostnames have a fixed length

		// read the characters from the hostname into a buffer in reverse
		v6chars := make([]byte, 32)
		for i, j := 62, 0; i >= 0; i -= 2 {
			v6chars[j] = name[i]
			j++
		}

		// decode the characters in the buffer into 16 bytes and return it
		v6bytes := make([]byte, 16)
		if _, err := hex.Decode(v6bytes, v6chars); err != nil {
			return nil
		}

		return v6bytes
	}

	return nil
}
