// Copyright (c) 2020 BitMaelum Authors
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

package dispatcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"text/template"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
	"github.com/sirupsen/logrus"
)

// default Slack Template when none is given by the webhook
var defaultSlackTemplate = `
Event *{{.meta.event}}* for account *{{.meta.account}}*:

` + "```" + `
{{json . -}}
` + "```"

// Work is the main function that will get dispatched as a job. It will do the actual work of sending data
func Work(w webhook.Type, payload interface{}) {
	switch w.Type {
	case webhook.TypeHTTP:
		_ = execHTTP(w, payload)
	case webhook.TypeSlack:
		_ = execSlack(w, payload)

	}
}

func execHTTP(w webhook.Type, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	cfg := &webhook.ConfigHTTP{}
	err = json.Unmarshal([]byte(w.Config), &cfg)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(cfg.URL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("webhook: invalid status code returned from HTTP endpoint")
	}

	return nil
}

func execSlack(w webhook.Type, payload interface{}) error {
	cfg := &webhook.ConfigSlack{}
	err := json.Unmarshal([]byte(w.Config), &cfg)
	if err != nil {
		return err
	}

	slackPayload := map[string]string{}
	if cfg.Channel != "" {
		slackPayload["channel"] = cfg.Channel
	}
	if cfg.Username != "" {
		slackPayload["username"] = cfg.Username
	}
	if cfg.IconEmoji != "" {
		slackPayload["icon_emoji"] = cfg.IconEmoji
	}
	if cfg.IconURL != "" {
		slackPayload["icon_url"] = cfg.IconURL
	}

	// Create text from (default) template and webhook payload
	templateBody := defaultSlackTemplate
	if cfg.Template != "" {
		templateBody = cfg.Template
	}

	logrus.Trace("template: ", templateBody)
	logrus.Trace("payload: ", payload)

	tmpl, err := template.New("slack").Funcs(template.FuncMap{
		"json": func(v interface{}) string {
			b, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return ""
			}

			return string(b)
		},
	}).Parse(templateBody)
	if err != nil {
		logrus.Trace("error while creating template: ", err)
		return err
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, payload)
	if err != nil {
		logrus.Trace("error while executing template: ", err)
		return err
	}
	slackPayload["text"] = sb.String()

	// Post slack payload
	slackPayloadBytes, err := json.Marshal(slackPayload)
	if err != nil {
		logrus.Trace("error while marshalling payload ", err)
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(cfg.WebhookURL, "application/json", bytes.NewReader(slackPayloadBytes))
	if err != nil {
		logrus.Trace("error while posting slack webhook ", err)
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("webhook: invalid status code returned from HTTP endpoint")
	}

	return nil
}
