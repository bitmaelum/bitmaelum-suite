// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gen2brain/beeep"
)

type options struct {
	Key      string        `short:"k" long:"key" description:"API key"`
	Config   string        `short:"c" long:"config" description:"Path to your configuration file"`
	Account  string        `short:"a" long:"account" description:"Account"`
	Interval time.Duration `short:"i" long:"interval" description:"Interval between checks" default:"1m"`
}

func main() {
	var opts options
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	addr, err := address.NewAddress(opts.Account)
	if err != nil {
		log.Fatal(err)
	}

	r := container.Instance.GetResolveService()
	info, err := r.ResolveAddress(addr.Hash())
	if err != nil {
		log.Fatal(err)
	}

	host := api.CanonicalHost(info.RoutingInfo.Routing)
	url := fmt.Sprintf("%s/account/%s/box/1", host, addr.Hash().String())

	lastChecked := internal.TimeNow()

	msgTicker := time.NewTicker(opts.Interval)
	for {
		// Wait for ticker
		<-msgTicker.C

		// Poll messages
		n := poll(url, opts.Key, lastChecked)
		lastChecked = internal.TimeNow()
		if n > 0 {
			// send notification if new messages found
			msg := fmt.Sprintf("%d new message(s) received for '%s'", n, addr.String())
			_ = beeep.Notify("New message(s) received", msg, "logo.png")
		}
	}
}

func poll(url string, key string, since time.Time) int {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url = fmt.Sprintf("%s?since=%d", url, since.Unix())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("Cannot create request: ", err)
		return 0
	}
	req.Header.Set("Authorization", "Bearer "+key)
	res, err := client.Do(req)
	if err != nil {
		log.Print("Cannot fetch from client: ", err)
		return 0
	}

	if res.StatusCode != 200 {
		log.Print("Not 200: ", res.StatusCode)
		return 0
	}

	type listType struct {
		MessageIds []string
		Meta       struct {
			Total    int
			Returned int
			Limit    int
			Offset   int
		}
	}

	l := &listType{}
	b, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		log.Print("Cannot read body")
		return 0
	}
	err = json.Unmarshal(b, &l)
	if err != nil {
		log.Print("Cannot unmarshal body")
		return 0
	}

	log.Printf("Returning %d items", l.Meta.Returned)
	return l.Meta.Returned
}
