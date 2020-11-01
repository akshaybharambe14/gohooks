package gohooks_test

import (
	"net/http"
	"testing"

	"github.com/averageflow/gohooks/gohooks"
)

func TestSend(t *testing.T) {
	secret := "2014716e-392c-4120-609e-555e295faff5" //nolint:gosec
	data := []int{1, 2, 3}

	hook := &gohooks.GoHook{}
	hook.Create(data, "int-resource", secret)

	resp, _ := hook.Send(RequestBinURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	hook.PreferredMethod = http.MethodPatch

	resp, _ = hook.Send(RequestBinURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	hook.PreferredMethod = http.MethodPut

	resp, _ = hook.Send(RequestBinURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	hook.PreferredMethod = http.MethodDelete

	resp, _ = hook.Send(RequestBinURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	hook.PreferredMethod = "invalid"

	resp, _ = hook.Send(RequestBinURL)
	if resp != nil {
		defer resp.Body.Close()
	}

	resp, _ = hook.Send("")
	if resp != nil {
		defer resp.Body.Close()
	}
}
