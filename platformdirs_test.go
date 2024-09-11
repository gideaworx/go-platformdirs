package platformdirs_test

import (
	"fmt"
	"testing"

	"go.gideaworx.io/platformdirs"
)

func TestWindows(t *testing.T) {
	p := platformdirs.New("me", "foo", "0.0.0")
	s, _ := p.UserConfigDir()
	fmt.Println(s)
}
