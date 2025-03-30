package icons_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/canpacis/pacis-ui/icons"
)

func TestIcon(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	icon := icons.Search()
	err := icon.Render(context.Background(), buf)
	if err != nil {
		t.Errorf("icon test failed: %s", err.Error())
	}
}
