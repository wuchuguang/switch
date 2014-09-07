// Copyright 2013 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package archive

import (
	"flag"
	"net"
	"net/http"
	"time"

	"github.com/gpmgo/switch/modules/log"
)

var (
	dialTimeout    = flag.Duration("dial_timeout", 10*time.Second, "Timeout for dialing an HTTP connection.")
	requestTimeout = flag.Duration("request_timeout", 20*time.Second, "Time out for roundtripping an HTTP request.")
)

func timeoutDial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, *dialTimeout)
}

type transport struct {
	t http.Transport
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	timer := time.AfterFunc(*requestTimeout, func() {
		t.t.CancelRequest(req)
		log.Warn("Canceled request for %s, please interrupt the program.", req.URL)
	})
	defer timer.Stop()
	resp, err := t.t.RoundTrip(req)
	return resp, err
}

var (
	httpTransport = &transport{
		t: http.Transport{
			Dial: timeoutDial,
			ResponseHeaderTimeout: *requestTimeout / 2,
		},
	}
	HttpClient = &http.Client{Transport: httpTransport}
)