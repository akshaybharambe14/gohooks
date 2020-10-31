package gohooks

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
)

// GoHook represents the definition of a GoHook.
type GoHook struct {
	// Data to be sent in the GoHook
	Payload GoHookPayload
	// The encrypted SHA resulting with the used salt
	ResultingSha string
	// Prepared JSON marshaled data
	PreparedData []byte
	// Choice of signature header to use on sending a GoHook
	SignatureHeader string
	// Should validate SSL certificate
	IsSecure bool
	// Preferred HTTP method to send the GoHook
	// Please choose only POST, DELETE, PATCH or PUT
	// Any other value will make the send use POST as fallback
	PreferredMethod string
}

// GoHookPayload represents the data that will be sent in the GoHook.
type GoHookPayload struct {
	Resource string      `json:"resource"`
	Data     interface{} `json:"data"`
}

// Create creates a webhook to be sent to another system, with a SHA 256 signature based on its contents.
func (hook *GoHook) Create(data interface{}, resource, secret string) {
	hook.Payload.Resource = resource
	hook.Payload.Data = data

	preparedHookData, err := json.Marshal(hook.Payload)
	if err != nil {
		log.Println(err.Error())
	}

	hook.PreparedData = preparedHookData

	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write(preparedHookData)

	if err != nil {
		log.Println(err.Error())
	}

	// Get result and encode as hexadecimal string
	hook.ResultingSha = hex.EncodeToString(h.Sum(nil))
}

// Send sends a GoHook to the specified URL.
func (hook *GoHook) Send(receiverURL string) (*http.Response, error) {
	if hook.SignatureHeader == "" {
		// Use the DefaultSignatureHeader as default if no custom header is specified
		hook.SignatureHeader = DefaultSignatureHeader
	}

	if !hook.IsSecure {
		// By default do not verify SSL certificate validity
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	if hook.PreferredMethod == "" || (hook.PreferredMethod != http.MethodPost &&
		hook.PreferredMethod != http.MethodPatch &&
		hook.PreferredMethod != http.MethodPut &&
		hook.PreferredMethod != http.MethodDelete) {
		// By default send GoHook using a POST method
		hook.PreferredMethod = http.MethodPost
	}

	client := &http.Client{}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		hook.PreferredMethod,
		receiverURL,
		bytes.NewBuffer(hook.PreparedData),
	)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Charset", "utf-8")
	req.Header.Add(DefaultSignatureHeader, hook.ResultingSha)
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}
