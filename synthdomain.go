package synthdomain

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"net"

	"github.com/miekg/dns"
)

const name = "synthdomain"

type SynthDomain struct {
	Next                             plugin.Handler
	forwardLookupSyntheticNetworks   []*net.IPNet
	reverseLookupToForwardLookupZone string
}

func (synth SynthDomain) Name() string { return name }

func (synth SynthDomain) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// run the next plugin.
	//		if the next plugin chain does not respond with an NXDOMAIN, we go with that.
	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	plugin.NextOrFailure(synth.Name(), synth.Next, ctx, rec, r)
	if rec.Rcode != dns.RcodeNameError {
		return plugin.NextOrFailure(synth.Name(), synth.Next, ctx, w, r)
	}

	var rr dns.RR

	// handle PTR requests
	if state.QType() == dns.TypePTR {
		ip := arpaToIp(state.QName())
		if ip != nil {
			rr = &dns.PTR{
				Hdr: header(state),
				Ptr: ipToName(ip, synth.ReverseLookupToForwardLookupZone()),
			}
		}
	}

	// Handle A request
	if state.QType() == dns.TypeA {
		ip := nameToIp(state.QName())
		if ip.To4() != nil {
			rr = &dns.A{Hdr: header(state), A: ip}
		} else {
			log.Debug(state.QName(), "\tCannot reply to A query with IPv6 address. No RR Added.")
		}
	}

	// Handle AAAA request
	if state.QType() == dns.TypeAAAA {
		ip := nameToIp(state.QName())
		if ip.To4() == nil {
			rr = &dns.AAAA{Hdr: header(state), AAAA: ip}
		} else {
			log.Debug(state.QName(), "\tCannot reply to AAAA query with IPv6 address. No RR Added.")
		}
	}

	if rr == nil {
		return plugin.NextOrFailure(synth.Name(), synth.Next, ctx, w, r)
	}

	// todo
	// 		- need to find some way to add the NS records in the authority section
	//			the recorded message has the SOA in the authority section, not the NS records.

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true
	a.Answer = []dns.RR{rr}
	w.WriteMsg(a)

	return dns.RcodeSuccess, nil

}

// The Reverse Lookup to Forward Lookup Zone Name
// This is used for the zone name to be used when responding to a reverse lookup
func (synth SynthDomain) ReverseLookupToForwardLookupZone() string {
	if synth.reverseLookupToForwardLookupZone == "" {
		return "local"
	}
	return synth.reverseLookupToForwardLookupZone
}

func header(state request.Request) dns.RR_Header {
	return dns.RR_Header{Name: state.QName(), Rrtype: state.QType(), Class: state.QClass(), Ttl: 3600}
}
