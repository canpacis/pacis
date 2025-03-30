package components_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/canpacis/pacis-ui/components"
)

func TestComponents(t *testing.T) {
	head := components.CreateHead("/public/")

	buf := bytes.NewBuffer([]byte{})
	err := head.Render(context.Background(), buf)

	if err != nil {
		t.Errorf("icon test failed: %s", err.Error())
	}
}
