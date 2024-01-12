# SOCKS5-Mixer

SOCKS4/5/5a mixer prototype. This prototype receives incoming SOCKS connections, consolidates them into a single connection to a server, and then separates or "demixes" each stream into its specific connection.

From the perspective of a third party sniffing traffic, they will only see one SOCKS connection between the client and the proxy (Server), instead of the N SOCKS connections that might be launched.

NOT FUNCTIONAL YET!

## TODO
* Make a configuration file for settings.
* Solve networking issue.
