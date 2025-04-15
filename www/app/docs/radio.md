# Radio Group

{.text-muted-foreground}
A set of checkable buttons—known as radio buttons—where no more than one of the buttons can be checked at a time.

{plate=0}
```go
Div(
  D{"submitted": ""},
  Form(
    On("submit.prevent", "submitted = new FormData($event.target).get('radio-group')"),

    RadioGroup(
      Name("radio-group"),
      Value("item-2"),

      RadioGroupItem(Value("item-1"), Text("Radio Item 1")),
      RadioGroupItem(Value("item-2"), Text("Radio Item 2")),
      RadioGroupItem(Value("item-3"), Text("Radio Item 3")),
    ),
    Button(Class("mt-2"), Type("submit"), Text("Submit")),
  ),
  Div(
    Class("mt-4"),

    P(X("show", "submitted.length > 0"), Span(Text("Submitted: ")), Span(X("text", "submitted"))),
  ),
)
```