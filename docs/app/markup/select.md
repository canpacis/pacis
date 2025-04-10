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
    SelectItem(Value("Templ"), Text("Templ")),
    SelectItem(Value("Gomponents"), Text("Gomponents")),
    SelectItem(Value("Pacis"), Text("Pacis")),
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
    SelectItem(Value("item-1"), Text("Item 1")),
    SelectItem(Value("item-2"), Text("Item 2")),
  ),
)
```

## Examples

### Clearable

{plate=1}
```go
Select(
  Name("library"),
  Class("min-w-[200px]"),
  Clearable

  SelectTrigger(
    Span(Text("Select a library")),
    Span(X("text", "value")),
  ),
  SelectContent(
    SelectItem(Value("Templ"), Text("Templ")),
    SelectItem(Value("Gomponents"), Text("Gomponents")),
    SelectItem(Value("Pacis"), Text("Pacis")),
  ),
)
```

### Default value

{plate=2}
```go
Select(
  Name("library"),
  Value("Pacis"),
  Class("min-w-[200px]"),
  Clearable

  SelectTrigger(
    Span(Text("Select a library")),
    Span(X("text", "value")),
  ),
  SelectContent(
    SelectItem(Value("Templ"), Text("Templ")),
    SelectItem(Value("Gomponents"), Text("Gomponents")),
    SelectItem(Value("Pacis"), Text("Pacis")),
  ),
)
```