# synthdomain

## Overview

*synthdomain* is a [CoreDNS](http://coredns.io) plugin to synthetically handle DNS records with IP addresses embedded.
Named after DNSMASQ's "synth-domain" [option](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html).

`synthdomain` aims to provide an easy mechanism for alignment between forward and reverse lookups.
This is a common DNS operational and configuration error as noted in [RFC1912](https://tools.ietf.org/html/rfc1912#section-2.1).

This plugin supports works nicely with the file plugin such that records present in the file will take precedence over this plugin.  

### Forward Lookups

Forward Lookups are hostname -> IP address. 
`synthdomain` supports IPs "embedded" in the DNS hostname. 
For IP addresses embedded in DNS hostnames the general model is `ip-<address>.example.com`
(where "address" can be either IPv4 or IPv6, and "example.com" is a domain of your choosing).
In IPv4 the dots are converted to hyphins; In IPv6 the colons are converted to hyphins.

The following are all considered valid for A or AAAA queries.

 * `ip-192-0-2-0.example.com`
 * `ip-2001-0db8-0000-0000-0000-0000-0000-0001.example.com`
 * `ip-2001-db8--1.example.com`

### Reverse Lookups

Reverse Lookups are IP -> hostname, and are known as pointer records (PTR).
`synthdomain` will respond to a PTR query and return a result that is also supported by the forward lookup mechanism.
Reverse lookups for IPv6 addresses will return a fully compressed IPv6 address (per [RFC5952](https://tools.ietf.org/html/rfc5952#section-2.2)).

## Corefile Configuration Examples

Reverse Lookup Example

~~~
2001:db8:abcd::/48 {
    synthdomain {
        forward example.com
    }
    file zones/a.b.c.d.8.d.b.0.1.0.0.2.ip6.arpa
}
~~~

Forward Lookup Example 

~~~
example.com {
    synthdomain {
        net 2001:db8:abcd::/48
    }
    file db.example.com
}
~~~


## Compiling into CoreDNS

To compile this with CoreDNS you can follow the [normal procedure](https://coredns.io/manual/plugins/#plugins) for external plugins. 
This plugin can be used by adding the following to `plugin.cfg`:
```
synthdomain:github.com/DonavanFritz/SynthDomain
```

## FAQ

### Why not use templates?

1- It appears that the `template` plugin is the recommended pattern for providing the resolution pattern we're after here.
However, it's not possible to have the `file` plugin provide the primary source of data and use a `template` at the same time. 
See [this](https://github.com/coredns/coredns/issues/2977#issuecomment-555938144) GitHub comment. 
Thus, it's not possible to have a PTR response from a file take priority over a template. 

2-  Using regex in a template for IPv4 and IPv6 addresses is very challanging with CIDR notation.
This plugin provides an easier experience by just providing an IP prefix in CIDR notation.  
