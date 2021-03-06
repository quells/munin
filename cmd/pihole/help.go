// Copyright 2021 Kai Wells
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

const info = "This graph shows information about DNS queries submitted to this Pi-Hole over a rolling 24-hour period (at the time of retrieval)."

var labels = map[string]string{
	"domains_being_blocked": "Block list count",
	"ads_blocked_today":     "Ads blocked",
	"unique_domains":        "Unique domains",
	"queries_forwarded":     "Queries forwarded",
	"queries_cached":        "Queries cached",
	"clients_ever_seen":     "Clients seen",
	"unique_clients":        "Unique clients",
	"dns_queries_today":     "DNS queries",
	"dns_queries_all_types": "Total queries",
	"reply_NODATA":          "Reply NODATA",
	"reply_NXDOMAIN":        "Reply NXDOMAIN",
	"reply_CNAME":           "Reply CNAME",
	"reply_IP":              "Reply IP",
	"privacy_level":         "Privacy level",
	"status":                "Status",
}

var infos = map[string]string{
	"domains_being_blocked": "Domains in ad block lists",
	"unique_domains":        "Unique domains resolved",
	"queries_forwarded":     "Queries forwarded to upstream resolver",
	"queries_cached":        "Queries served from cache",
	"dns_queries_today":     "Total queries served",
	"dns_queries_all_types": "Total queries served",
	"reply_NODATA":          "Queries resolved with NODATA",
	"reply_NXDOMAIN":        "Queries resolved with NXDOMAIN",
	"reply_CNAME":           "Queries resolved with CNAME",
	"reply_IP":              "Queries resolved with IP",
	"status":                "1 for enabled, 0 for disabled",
}

const help = `Pi-Hole stats Munin plugin.

Must set env.host in configuration for Pi-Hole web admin interface, including scheme e.g. http://pi.hole

Can optionally set env.except to comma separated list of values to skip reporting. Valid entries are:
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
`
