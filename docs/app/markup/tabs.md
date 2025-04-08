# Tabs

{.text-muted-foreground}
A set of layered sections of content—known as tab panels—that are displayed one at a time.

{plate=0}

```go
Tabs(
  Value("tab-item-1"), // <- Default value

  TabList(
    TabTrigger(Text("Tab Item 1"), Value("tab-item-1")), // Value attributes are required
    TabTrigger(Text("Tab Item 2"), Value("tab-item-2")),
  ),
  TabContent(
    Value("tab-item-1"), // Value attributes are required

    P(Text("Tab item 1 content")),
  ),
  TabContent(
    Value("tab-item-2"),

    P(Text("Tab item 2 content")),
  ),
)
```

## Usage

```go
Tabs(
  TabList(
    TabTrigger(Text("Trigger"), Value("value")),
  ),
  TabContent(
    Value("value"),

    // Content
  ),
)
```
