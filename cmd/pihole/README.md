# Pi-Hole Munin Plugin

Gather statistics from the [Pi-Hole](https://pi-hole.net/) web admin API such as number of queries, blocked ads, etc.

## Installation

Build this Go module and place the binary in the plugins directory (usually `/etc/munin/plugins`). Some minimum configuration is required.

## Configuration

The `host` where the Pi-Hole web admin interface can be found must be specified, including scheme. The plugin reads from `$host/admin/api.php?summary` to get the stats.

Values from this response can be optionally omitted using the `except` environment variable.

Example for `/etc/munin/plugin-conf.d/pihole`:

```
[pihole]
 env.host http://pi.hole
 env.except privacy_level,status
```

Valid values to omit:

- domains_being_blocked
- ads_blocked_today
- unique_domains
- queries_forwarded
- queries_cached
- clients_ever_seen
- unique_clients
- dns_queries_today
- dns_queries_all_types
- reply_NODATA
- reply_NXDOMAIN
- reply_CNAME
- reply_IP
- privacy_level
- status
