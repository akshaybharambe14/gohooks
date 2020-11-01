package gohooks_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

var RequestBinURL string //nolint:gochecknoglobals

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)

	binID := CreateRequestBin()
	if binID == "" {
		os.Exit(1)
	}

	RequestBinURL = fmt.Sprintf("https://postb.in/%s", binID)
	fmt.Println(RequestBinURL)

	code := m.Run()
	os.Exit(code)
}

type NewBinResponse struct {
	BinID    string `json:"binId"`
	Inserted int64  `json:"inserted"`
	Expiry   int64  `json:"expiry"`
}

func CreateRequestBin() string {
	resp, err := http.Post(
		"https://postb.in/api/bin",
		"",
		nil,
	)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if resp == nil {
		fmt.Println("Obtained no response from request bin!")
		return ""
	}

	defer resp.Body.Close()

	var newBinResponse NewBinResponse

	err = json.NewDecoder(resp.Body).Decode(&newBinResponse)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return newBinResponse.BinID
}
