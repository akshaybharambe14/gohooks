package gohooks_test

import (
	"encoding/json"
	"testing"

	"github.com/averageflow/gohooks/gohooks"
)

func TestIsGoHookValid(t *testing.T) {
	secret := "0014716e-392c-4120-609e-555e295faff5" //nolint:gosec
	data := []int{1, 2, 3}

	hook := &gohooks.GoHook{}
	hook.Create(data, "int-resource", secret)

	type WebhookData struct {
		Resource string `json:"resource"`
		Data     []int  `json:"data"`
	}

	var hookData WebhookData

	_ = json.Unmarshal(hook.PreparedData, &hookData)

	isValid := gohooks.IsGoHookValid(hookData, hook.ResultingSha, secret)
	if !isValid {
		t.Errorf("Expected first GoHook to be valid")
	}

	isValid = gohooks.IsGoHookValid(hook.PreparedData, "invalid-signature", secret)
	if isValid {
		t.Errorf("Expected first GoHook to be invalid")
	}
}
