package main

import (
	"fmt"
	"os"
)

// buildExampleWindow baut das Fenster von der Tafel nach:
//
//	┌─────┬─────┐
//	│  1  │  2  │
//	├─────┼──┬──┤
//	│     │a │b │
//	│  3  ├──┼──┤   unten rechts: Unterfenster mit 4 Scheiben
//	│     │c │d │
//	└─────┴──┴──┘
func buildExampleWindow() (*Window, error) {
	// Kleine Hilfsfunktion, damit die Fehlerprüfung nicht jede Zeile aufbläht.
	var firstErr error
	glass := func(name string, cents int) Item {
		g, err := NewGlass(name, cents)
		if err != nil && firstErr == nil {
			firstErr = err
		}
		return g
	}

	frameSub, err := NewFrame("Sprossenrahmen", 2500)
	if err != nil {
		return nil, err
	}
	sub, err := NewWindow("Unterfenster unten rechts",
		frameSub,
		glass("klein a", 1000), glass("klein b", 1000),
		glass("klein c", 1000), glass("klein d", 1000),
	)
	if firstErr != nil {
		return nil, firstErr
	}
	if err != nil {
		return nil, err
	}

	frameMain, err := NewFrame("Hauptrahmen", 8000)
	if err != nil {
		return nil, err
	}
	root, err := NewWindow("Hauptfenster",
		frameMain,
		glass("oben links", 3000),
		glass("oben rechts", 3000),
		glass("unten links", 3000),
		sub, // Composite im Composite
	)
	if firstErr != nil {
		return nil, firstErr
	}
	return root, err
}

func main() {
	w, err := buildExampleWindow()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Fehler:", err)
		os.Exit(1)
	}
	fmt.Println(w.Display())
	fmt.Println()
	fmt.Println("Gesamtpreis:", formatEuro(w.Price()))
}
