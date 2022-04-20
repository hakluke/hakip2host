# hakip2host
hakip2host takes a list of IP addresses via stdin, then does a series of checks to return associated domain names.

Current supported checks are:

- DNS PTR lookups
- Subject Alternative Names (SANs) on SSL certificates
- Common Names (CNs) on SSL certificates

## Installation

Install golang, then:

```
go install github.com/hakluke/hakip2host@latest
```

## Help

```
./hakip2host --help
```

## Example usage

```
hakluke$ prips 173.0.84.0/24 | ./hakip2host
[DNS-PTR] 173.0.84.23 new-creditcenter.paypal.com.
[DNS-PTR] 173.0.84.11 slc-a-origin-www-1.paypal.com.
[DNS-PTR] 173.0.84.10 admin.paypal.com.
[DNS-PTR] 173.0.84.30 ss-www.paypal.com.
[DNS-PTR] 173.0.84.5 www.gejscript-paypal.com.
[DNS-PTR] 173.0.84.24 slc-a-origin-demo.paypal.com.
[DNS-PTR] 173.0.84.20 origin-merchantweb.paypal.com.
[SSL-SAN] 173.0.84.67 uptycspay.paypal.com
[SSL-SAN] 173.0.84.67 a.paypal.com
[SSL-CN] 173.0.84.67 api.paypal.com
[SSL-SAN] 173.0.84.76 svcs.paypal.com
[SSL-SAN] 173.0.84.76 uptycshon.paypal.com
[SSL-SAN] 173.0.84.76 uptycshap.paypal.com
[SSL-SAN] 173.0.84.76 uptycsven.paypal.com
[SSL-SAN] 173.0.84.76 uptycsize.paypal.com
[SSL-SAN] 173.0.84.76 uptycspay.paypal.com
[SSL-CN] 173.0.84.76 svcs.paypal.com
```
