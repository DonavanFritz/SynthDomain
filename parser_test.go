package synthdomain

import (
	"net"
	"testing"
)

func TestNameToIp(t *testing.T) {
	testCases := []struct {
		name string
		ip   net.IP
	}{
		{"ip-192-0-2-0.example.com", net.ParseIP("192.0.2.0")},
		{"ip-192-0-2-100.a.b.c.e.f.g", net.ParseIP("192.0.2.100")},
		{"ip-192-0-2-255.tld", net.ParseIP("192.0.2.255")},
		{"ip-2001-db8--100.example.com", net.ParseIP("2001:db8::100")},
		{"ip-2001-db8--beef.a.b.c.d.e.f.g", net.ParseIP("2001:db8::beef")},
		{"ip---1.example.com", net.ParseIP("::1")},
		{"ip-127-0-0-1.example.com", net.ParseIP("127.0.0.1")},
		{"foobar", nil},
		{"pi-192-0-2-0.example.com", nil},
		{"192-0-2-0.example.com", nil},
		{"ip-256-0-2-0.example.com", nil},
	}

	for _, tt := range testCases {
		if result := nameToIp(tt.name); !result.Equal(tt.ip) {
			t.Errorf("expected '%s' for '%s' but got '%s'", tt.ip, tt.name, result)
		}
	}
}

func TestIpToName(t *testing.T) {
	testCases := []struct {
		ip   net.IP
		zone string
		name string
	}{
		{net.ParseIP("192.0.2.0"), "example.com", "ip-192-0-2-0.example.com."},
		{net.ParseIP("192.0.2.0"), ".example.com", "ip-192-0-2-0.example.com."},
		{net.ParseIP("192.0.2.0"), "example.com.", "ip-192-0-2-0.example.com."},
		{net.ParseIP("192.0.2.0"), ".example.com.", "ip-192-0-2-0.example.com."},
		{net.ParseIP("2001:db8::1"), "example.com.", "ip-2001-db8--1.example.com."},
		{net.ParseIP("2001:db8::100"), "example.com.", "ip-2001-db8--100.example.com."},
		{net.ParseIP("::1"), "example.com.", "ip---1.example.com."},
		{nil, "example.com.", ""},
	}

	for _, tt := range testCases {
		if result := ipToName(tt.ip, "example.com"); result != tt.name {
			t.Errorf("expected '%s' for '%s' but got '%s'", tt.name, tt.ip, result)
		}
	}
}

func TestInArpa(t *testing.T) {
	testCases := []struct {
		name string
		ip   net.IP
	}{
		{"0.2.0.192.in-addr.arpa.", net.ParseIP("192.0.2.0")},
		{"0.2.0.300.in-addr.arpa.", nil},
		{"0.1.0.1.in-addr.arpa.", net.ParseIP("1.0.1.0")},
		{"foobar.in-addr.arpa.", nil},
		{"foobar", nil},
		{"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.", net.ParseIP("2001:db8::1")},
		{"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.", nil},
		{"foobar.ip6.arpa.", nil},
	}

	for _, tt := range testCases {
		if result := arpaToIp(tt.name); !result.Equal(tt.ip) {
			t.Errorf("expected '%s' for '%s' but got '%s'", tt.ip, tt.name, result)
		}
	}
}

var inArpaBenchmarks = []string{
	"0.2.0.192.in-addr.arpa.",
	"0.2.0.in-addr.arpa.",
	"foobar.in-addr.arpa.",
	"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
	"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
	"foobar.ip6.arpa.",
}

func BenchmarkInArpa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range inArpaBenchmarks {
			arpaToIp(test)
		}
	}
}
