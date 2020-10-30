package gohooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

// isGoHookValid checks the sha256 of the data matches the one given on the signature.
func isGoHookValid(data interface{}, signature, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	preparedData, _ := json.Marshal(data)

	_, err := h.Write(preparedData)
	if err != nil {
		log.Println(err.Error())
	}

	sha := hex.EncodeToString(h.Sum(nil))

	return sha == signature
}