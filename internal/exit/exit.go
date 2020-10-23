package exit

import (
	"fmt"
	"os"
	"strings"
)

func Exit(status int) {
	os.Exit(status)
}

func WithMessage(status int, t string, args ...interface{}) {
	t = strings.TrimSpace(t) + "\n"
	out := os.Stdout
	if status != 0 {
		out = os.Stderr
	}
	if len(args) > 0 {
		fmt.Fprintf(out, t, args...)
	} else {
		fmt.Fprint(out, t)
	}
	Exit(status)
}
