package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"

	"github.com/jordanorelli/belt-mud/internal/exit"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		exit.WithMessage(1, "unable to initialize screen: %v", err)
	}
	fmt.Println(screen)
}
