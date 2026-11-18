# Prompt für Claude Code

```text
Implementiere in Go (ein Modul, Standardbibliothek only) das Composite Pattern
nach folgendem UML:

- Interface Item mit Price() und Display().
- Leafs: Glass (Glasscheibe) und Frame (Rahmen), je mit Name und Preis.
- Composite: Window enthält 1..* Items (also auch weitere Windows).
  Price() = Summe der Kind-Preise, Display() = eingerückte Baumdarstellung.

Anforderungen:
- Preise als int in Cent (keine float-Rundungsfehler), Ausgabe als "12,34 €".
- Display() liefert einen string (testbar), main gibt ihn aus.
- Fehlerbehandlung: NewWindow ohne Items, nil-Item, negativer Preis und
  Zyklen (Fenster in sich selbst) liefern error.
- main.go baut dieses Fenster nach: Rahmen + 2 Scheiben oben + 1 Scheibe
  unten links + unten rechts ein Unterfenster (eigener Rahmen + 4 Scheiben).
  Ausgabe: Baum und Gesamtpreis.
- Tabellengetriebene Tests (go test) für Preise, Display und Fehlerfälle.
- Code gut kommentiert (deutsch), go vet sauber, kurze README mit go run/test.
```
