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
	Payload         GoHookPayload
	ResultingSha    string
	PreparedData    []byte
	SignatureHeader string
	Secure          bool
}

type GoHookPayload struct {
	Resource string      `json:"resource"`
	Data     interface{} `json:"data"`
}

// Create creates a webhook to be sent to another system, with a SHA 256 signature based on its contents.
func (hook *GoHook) Create(data interface{}, resource, secret string) {
	hook.Payload.Resource = resource
	hook.Payload.Data = data

	preparedHookData, err := json.Marshal(hook.Payload.Data)
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
		hook.SignatureHeader = DefaultSignatureHeader
	}

	if !hook.Secure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	client := &http.Client{}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, receiverURL, bytes.NewBuffer(hook.PreparedData))

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
