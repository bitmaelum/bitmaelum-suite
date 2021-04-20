package internal

import (
	"bytes"
	"fmt"
	"text/template"
)

const AccountCreatedTemplate = `
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your account. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:

{{ .mnemonic }}

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.

WITHOUT THESE WORDS, ALL ACCESS TO YOUR ACCOUNT IS LOST!
`

const OrganisationCreatedTemplate = `
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your organisation. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:

{{ .mnemonic }}

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.

WITHOUT THESE WORDS, ALL ACCESS TO YOUR ACCOUNT IS LOST!
`

const AccountProofTemplate = `could not find proof in the DNS.

In order to register this reserved address, make sure you add the following information to the DNS:

   _bitmaelum TXT {{ .Fingerprint }}

This entry could be added to any of the following domains: {{ .Domains }}. Once we have found the entry, we can 
register the account onto the keyserver. For more information, please visit https://bitmaelum.com/reserved
`

// Generate from generic template data
func generateFromTemplateData(messageTemplate string, data interface{}) string {
	msg := fmt.Sprintf("%v", data) // when things fail
	tmpl, err := template.New("template").Parse(messageTemplate)
	if err != nil {
		return msg
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return msg
	}

	return buf.String()
}

// Generate fingerprint data
func GenerateFromFingerprintTemplate(messageTemplate string, fingerprint string, domains []string) string {
	type tplData struct {
		Fingerprint string
		Domains     []string
	}

	return generateFromTemplateData(messageTemplate, tplData{
		Fingerprint: fingerprint,
		Domains:     domains,
	})
}

// Generate mnemonic data
func GenerateFromMnemonicTemplate(messageTemplate string, mnemonic string) string {
	type tplData struct {
		Mnemonic string
	}

	return generateFromTemplateData(messageTemplate, tplData{
		Mnemonic: mnemonic,
	})
}
