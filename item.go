// Package main implementiert das Composite Pattern am Beispiel eines Fensters.
//
// UML (Tafelbild):
//
//	       ┌───────────┐ 1..*
//	       │   Item    │◄────────┐
//	       │ Price()   │         │
//	       │ Display() │         │
//	       └─────▲─────┘         │
//	   ┌─────────┼─────────┐     │
//	Window     Glass     Frame   │
//	   └─────────────────────────┘ (Window enthält Items)
//
// Glass und Frame sind Blätter (Leafs), Window ist das Kompositum.
package main

import (
	"errors"
	"fmt"
	"strings"
)

// Item ist die gemeinsame Schnittstelle für Blätter und Komposita.
// Der Client behandelt einzelne Teile und ganze Fenster identisch.
type Item interface {
	// Price liefert den Preis in Cent (int vermeidet Rundungsfehler von float).
	Price() int
	// Display liefert eine (ggf. mehrzeilige) Textdarstellung des Items.
	Display() string
}

// Fehlerwerte – exportiert, damit Aufrufer/Tests sie mit errors.Is prüfen können.
var (
	ErrNegativePrice = errors.New("preis darf nicht negativ sein")
	ErrNilItem       = errors.New("item darf nicht nil sein")
	ErrNoItems       = errors.New("ein Fenster braucht mindestens ein Item (1..*)")
	ErrCycle         = errors.New("zyklus: fenster würde sich selbst enthalten")
)

// formatEuro wandelt Cent in eine deutsche Preisangabe um, z. B. 1234 -> "12,34 €".
func formatEuro(cents int) string {
	sign := ""
	if cents < 0 { // nur zur Sicherheit, Konstruktoren verhindern negative Preise
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%d,%02d €", sign, cents/100, cents%100)
}

// ---------------------------------------------------------------------------
// Leafs
// ---------------------------------------------------------------------------

// part bündelt die gemeinsamen Felder der Blätter (Name + Preis).
type part struct {
	name  string
	price int // Cent
}

// newPart prüft die Eingaben zentral für alle Blatt-Typen.
func newPart(name string, priceCents int) (part, error) {
	if priceCents < 0 {
		return part{}, fmt.Errorf("%q: %w", name, ErrNegativePrice)
	}
	return part{name: name, price: priceCents}, nil
}

// Glass ist eine einzelne Glasscheibe (Leaf).
type Glass struct{ part }

// NewGlass erstellt eine Glasscheibe; negativer Preis -> Fehler.
func NewGlass(name string, priceCents int) (*Glass, error) {
	p, err := newPart(name, priceCents)
	if err != nil {
		return nil, err
	}
	return &Glass{p}, nil
}

// Price liefert den Preis der Scheibe in Cent.
func (g *Glass) Price() int { return g.price }

// Display liefert eine einzeilige Beschreibung.
func (g *Glass) Display() string {
	return fmt.Sprintf("Glas: %s (%s)", g.name, formatEuro(g.price))
}

// Frame ist ein Fensterrahmen (Leaf).
type Frame struct{ part }

// NewFrame erstellt einen Rahmen; negativer Preis -> Fehler.
func NewFrame(name string, priceCents int) (*Frame, error) {
	p, err := newPart(name, priceCents)
	if err != nil {
		return nil, err
	}
	return &Frame{p}, nil
}

// Price liefert den Preis des Rahmens in Cent.
func (f *Frame) Price() int { return f.price }

// Display liefert eine einzeilige Beschreibung.
func (f *Frame) Display() string {
	return fmt.Sprintf("Rahmen: %s (%s)", f.name, formatEuro(f.price))
}

// ---------------------------------------------------------------------------
// Composite
// ---------------------------------------------------------------------------

// Window ist das Kompositum: es enthält 1..* Items (Scheiben, Rahmen
// oder weitere Fenster) und delegiert Price/Display an seine Kinder.
type Window struct {
	name  string
	items []Item
}

// NewWindow erstellt ein Fenster mit mindestens einem Item (Multiplizität 1..*).
func NewWindow(name string, items ...Item) (*Window, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("%q: %w", name, ErrNoItems)
	}
	w := &Window{name: name}
	for _, it := range items {
		if err := w.Add(it); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// Add fügt ein Item hinzu. Verhindert nil-Items und Zyklen
// (ein Zyklus würde Price/Display in eine Endlosrekursion schicken).
func (w *Window) Add(it Item) error {
	if isNil(it) {
		return fmt.Errorf("%q: %w", w.name, ErrNilItem)
	}
	if contains(it, w) {
		return fmt.Errorf("%q: %w", w.name, ErrCycle)
	}
	w.items = append(w.items, it)
	return nil
}

// Items liefert eine Kopie der Kinder, damit Aufrufer die interne
// Liste nicht unbemerkt verändern (z. B. 1..* unterlaufen) können.
func (w *Window) Items() []Item {
	out := make([]Item, len(w.items))
	copy(out, w.items)
	return out
}

// Price summiert rekursiv die Preise aller enthaltenen Items.
func (w *Window) Price() int {
	sum := 0
	for _, it := range w.items {
		sum += it.Price()
	}
	return sum
}

// Display liefert den Baum: Kopfzeile mit Gesamtpreis, darunter die Kinder
// eingerückt. Mehrzeilige Kinder (Unterfenster) werden zeilenweise eingerückt.
func (w *Window) Display() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Fenster: %s (gesamt %s)", w.name, formatEuro(w.Price()))
	for _, it := range w.items {
		for _, line := range strings.Split(it.Display(), "\n") {
			b.WriteString("\n  ")
			b.WriteString(line)
		}
	}
	return b.String()
}

// contains prüft, ob target in item (rekursiv) enthalten ist bzw. item selbst ist.
func contains(item Item, target *Window) bool {
	win, ok := item.(*Window)
	if !ok || win == nil {
		return false // Blätter können nichts enthalten
	}
	if win == target {
		return true
	}
	for _, child := range win.items {
		if contains(child, target) {
			return true
		}
	}
	return false
}

// isNil erkennt auch "typisierte nil-Zeiger" in einem Interface
// (z. B. ein (*Glass)(nil)), die ein einfaches it == nil nicht erwischt.
func isNil(it Item) bool {
	switch v := it.(type) {
	case nil:
		return true
	case *Glass:
		return v == nil
	case *Frame:
		return v == nil
	case *Window:
		return v == nil
	}
	return false
}
