# Composite Pattern in Go – Fenster-Beispiel

Umsetzung des in der Vorlesung besprochenen Beispiels (Aufgabe 2).
Der verwendete Prompt steht in [`PROMPT.md`](PROMPT.md).

| Rolle im Pattern | Typ      | Bedeutung                          |
|------------------|----------|------------------------------------|
| Component        | `Item`   | Interface mit `Price()`, `Display()` |
| Leaf             | `Glass`  | Glasscheibe                        |
| Leaf             | `Frame`  | Rahmen                             |
| Composite        | `Window` | enthält 1..* Items, auch Fenster   |

## Ausführen

```bash
go run .
go test -v ./...
```

## Ausgabe

```text
Fenster: Hauptfenster (gesamt 235,00 €)
  Rahmen: Hauptrahmen (80,00 €)
  Glas: oben links (30,00 €)
  Glas: oben rechts (30,00 €)
  Glas: unten links (30,00 €)
  Fenster: Unterfenster unten rechts (gesamt 65,00 €)
    Rahmen: Sprossenrahmen (25,00 €)
    Glas: klein a (10,00 €)
    Glas: klein b (10,00 €)
    Glas: klein c (10,00 €)
    Glas: klein d (10,00 €)

Gesamtpreis: 235,00 €
```
