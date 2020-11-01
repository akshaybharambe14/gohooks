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

	_, _ = hook.Send(RequestBinURL)

	hook.PreferredMethod = http.MethodPatch
	_, _ = hook.Send(RequestBinURL)

	hook.PreferredMethod = http.MethodPut
	_, _ = hook.Send(RequestBinURL)

	hook.PreferredMethod = http.MethodDelete
	_, _ = hook.Send(RequestBinURL)

	hook.PreferredMethod = "invalid"
	_, _ = hook.Send(RequestBinURL)

	_, _ = hook.Send("")
}
