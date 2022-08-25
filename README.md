# Kubernetes NGINX Ingress Whitelist

This is meant to be used as a CronJob that checks the list of IPv4 and IPv6 CIDRs from CloudFlare and adds a global
whitelist to only allow traffic that's been proxied through CloudFlare.