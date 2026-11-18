package main

import (
	"errors"
	"strings"
	"testing"
)

// mustGlass / mustFrame: Test-Helfer, brechen den Test bei Fehlern ab.
func mustGlass(t *testing.T, name string, c int) *Glass {
	t.Helper()
	g, err := NewGlass(name, c)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func mustFrame(t *testing.T, name string, c int) *Frame {
	t.Helper()
	f, err := NewFrame(name, c)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestPrice(t *testing.T) {
	sub, _ := NewWindow("sub", mustFrame(t, "r", 100), mustGlass(t, "g", 50))
	tests := []struct {
		name string
		item Item
		want int
	}{
		{"Glas", mustGlass(t, "g", 1234), 1234},
		{"Rahmen", mustFrame(t, "r", 0), 0},
		{"flaches Fenster", sub, 150},
	}
	nested, _ := NewWindow("haupt", mustFrame(t, "r", 1000), sub)
	tests = append(tests, struct {
		name string
		item Item
		want int
	}{"verschachtelt", nested, 1150})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Price(); got != tt.want {
				t.Errorf("Price() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestExampleWindow(t *testing.T) {
	w, err := buildExampleWindow()
	if err != nil {
		t.Fatal(err)
	}
	// 80 + 3*30 + (25 + 4*10) = 235,00 €
	if got := w.Price(); got != 23500 {
		t.Errorf("Price() = %d, want 23500", got)
	}
	out := w.Display()
	for _, want := range []string{
		"Fenster: Hauptfenster (gesamt 235,00 €)",
		"  Fenster: Unterfenster unten rechts (gesamt 65,00 €)",
		"    Glas: klein d (10,00 €)", // doppelt eingerückt = verschachtelt
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Display() fehlt %q\n%s", want, out)
		}
	}
}

func TestErrors(t *testing.T) {
	g := mustGlass(t, "g", 1)
	var nilGlass *Glass

	t.Run("negativer Preis", func(t *testing.T) {
		if _, err := NewGlass("x", -1); !errors.Is(err, ErrNegativePrice) {
			t.Errorf("got %v", err)
		}
		if _, err := NewFrame("x", -1); !errors.Is(err, ErrNegativePrice) {
			t.Errorf("got %v", err)
		}
	})
	t.Run("keine Items", func(t *testing.T) {
		if _, err := NewWindow("leer"); !errors.Is(err, ErrNoItems) {
			t.Errorf("got %v", err)
		}
	})
	t.Run("nil Item", func(t *testing.T) {
		if _, err := NewWindow("w", nil); !errors.Is(err, ErrNilItem) {
			t.Errorf("got %v", err)
		}
		if _, err := NewWindow("w", nilGlass); !errors.Is(err, ErrNilItem) {
			t.Errorf("typisiertes nil: got %v", err)
		}
	})
	t.Run("Zyklus", func(t *testing.T) {
		a, _ := NewWindow("a", g)
		if err := a.Add(a); !errors.Is(err, ErrCycle) {
			t.Errorf("selbst: got %v", err)
		}
		b, _ := NewWindow("b", a)
		if err := a.Add(b); !errors.Is(err, ErrCycle) {
			t.Errorf("indirekt: got %v", err)
		}
	})
	t.Run("Items liefert Kopie", func(t *testing.T) {
		w, _ := NewWindow("w", g)
		items := w.Items()
		items[0] = nil
		if w.Price() != 1 {
			t.Error("interne Liste wurde verändert")
		}
	})
}

func TestFormatEuro(t *testing.T) {
	for in, want := range map[int]string{0: "0,00 €", 5: "0,05 €", 1234: "12,34 €", -150: "-1,50 €"} {
		if got := formatEuro(in); got != want {
			t.Errorf("formatEuro(%d) = %q, want %q", in, got, want)
		}
	}
}
