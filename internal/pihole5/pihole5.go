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

package pihole5

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/quells/munin/internal/set"
	"github.com/quells/munin/pkg/munin"
)

type Client struct {
	host string
	skip set.Strings
}

func NewClient(host string, skip set.Strings) *Client {
	c := new(Client)
	c.host = host
	c.skip = skip
	return c
}

func (c *Client) Load() (values munin.Values, precision munin.Precision, err error) {
	if c == nil {
		err = fmt.Errorf("nil pihole5 config")
		return
	}

	u := fmt.Sprintf("%s/admin/api.php?summary", c.host)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	var respData []byte
	respData, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	respBody := make(map[string]interface{})
	if err = json.Unmarshal(respData, &respBody); err != nil {
		return
	}

	values, precision = c.filter(respBody)
	return
}

func (c *Client) filter(raw map[string]interface{}) (values munin.Values, precision munin.Precision) {
	values = make(munin.Values)

	for k, vint := range raw {
		if _, skip := c.skip[k]; skip {
			continue
		}

		if vstr, ok := vint.(string); ok {
			if k == "status" {
				if vstr == "enabled" {
					values[k] = 1
				} else {
					values[k] = 0
				}
				continue
			}

			if x, err := strconv.Atoi(strings.ReplaceAll(vstr, ",", "")); err == nil {
				values[k] = float64(x)
			}
		}
	}

	return
}
