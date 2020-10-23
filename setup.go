package synthdomain

import (
	"fmt"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"net"
)

const pluginName = "synthdomain" // name to be used in Corefile

var log = clog.NewWithPlugin(pluginName)

func init() {
	plugin.Register(pluginName, setup)
}

func setup(controller *caddy.Controller) error {
	s, err := fromCoreFileConfiguration(controller)

	if err != nil {
		return plugin.Error(pluginName, err)
	}

	SynthDomainPlugInChain := func(next plugin.Handler) plugin.Handler {
		s.Next = next
		return s
	}

	dnsServerConfiguration := dnsserver.GetConfig(controller)
	dnsServerConfiguration.AddPlugin(SynthDomainPlugInChain)

	controller.Next()
	if controller.NextArg() {
		return plugin.Error(pluginName, controller.ArgErr())
	}

	return nil
}

func fromCoreFileConfiguration(c *caddy.Controller) (*SynthDomain, error) {
	s := new(SynthDomain)

	for c.Next() {
		for c.NextBlock() {
			switch v := c.Val(); v {

			// Configuration for forward lookup zones for which to do resolution
			case "allow":
				args := c.RemainingArgs()
				for _, arg := range args {
					_, cidr, err := net.ParseCIDR(arg)
					if err == nil {
						fmt.Println(cidr)
					} else {
						clog.Errorf("'%s' is not a valid CIDR", arg)
					}
				}

			// Configuration for reverse lookup zones for the forward lookup zone name
			case "forward":
				args := c.RemainingArgs()
				for _, arg := range args {
					s.reverseLookupToForwardLookupZone = arg
				}

			default:
				return nil, c.Errf("unknown property '%s'", v)
			}
		}
	}

	return s, nil
}
