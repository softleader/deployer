package slack

import (
	"fmt"
	"testing"
)

func TestReplaceLast(t *testing.T) {
	fmt.Println(replaceLast("asd---asd---vzxc", "---", "/"))
}
