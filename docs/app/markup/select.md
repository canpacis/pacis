# Select

{.text-muted-foreground}
Displays a list of options for the user to pick from, triggered by a button.

{plate=0}
```go
Select(
  Name("library"),
  Class("min-w-[200px]"),

  SelectTrigger(
    Span(Text("Select a library")),
    Span(X("text", "value")),
  ),
  SelectContent(
    SelectItem(Text("Templ"), Value("Templ")),
    SelectItem(Text("Gomponents"), Value("Gomponents")),
    SelectItem(Text("Pacis"), Value("Pacis")),
  ),
)
```

## Usage

```go
Select(
  Name("name"),

  SelectTrigger(
    Span(Text("Empty Trigger")),
    Span(X("text", "value")), // <- Selected value
  ),
  SelectContent(
    SelectItem(Text("Item 1"), Value("item-1")),
    SelectItem(Text("Item 2"), Value("item-2")),
  ),
)
```

## Examples

### Clearable

{plate=1}
```go
Select(
  Clearable
  Name("library"),
  Class("min-w-[200px]"),

  SelectTrigger(
    Span(Text("Select a library")),
    Span(X("text", "value")),
  ),
  SelectContent(
    SelectItem(Text("Templ"), Value("Templ")),
    SelectItem(Text("Gomponents"), Value("Gomponents")),
    SelectItem(Text("Pacis"), Value("Pacis")),
  ),
)
```