# Parser

Parser is the CSS parser for this project, it is made up of two primary functions `ParseCSS`, `ParseStyleAttribute`, and a few other functions designed to help with the parsing.

## ParseCSS?(go)

`ParseCSS` is the function for reading CSS files. It is a RegExp based parser which it converts CSS definitions into a 2d map of strings. If you want to convert the values into absoulte values there are helper functions in the [utils](/utils) documentation.

> matches := selectorRegex.FindAllStringSubmatch(css, -1)

First we start off by using a RegExp to find the individual CSS blocks and to sort them into the block selector and the styles for the selector.

> selectors := parseSelectors(selectorBlock)

The mapped values are defined by the selector pulled from the [parseSelectors](./#parseselectorsgo) function and will include the entire name (this includes the symbol ".","#", and ",")

> selectorMap[selector] = parseStyles(styleBlock)

Once the selectors of the file have been parsed, the styles are mapped to the second level of the map with their respective key and value pulled from [parseStyles](./#parsestylesgo).

> NOTE: When parsing duplicate selectors and styles will be merged with the last conflicting selector/style overriding the prevous.

### Implementation

> styles := parser.ParseCSS(string(dat))

The only time `ParseCSS` is used is in the `cstyle` package, and it used to add css files in the first example

> styles := parser.ParseCSS(css)

and style tags in the next. As you can see in both examples those functions are for appending the new styles to the current global CSS stylesheet held within the instance of the CSS struct (`CSS.StyleSheets`).

> NOTE: Style tag is refering to the below

```html
<style>
  table td.r,
  table th.r {
    text-align: center;
  }
</style>
```

## parseSelectors?(go)

`parseSelectors` takes the first output of the RegExp match in [ParseCSS](./#parsecssgo) and splits it up by commas.

### parseSelectors Example

```go

selectorBlock := `table td.r,
table th.r`

parseSelectors(selectorBlock)

// Output
[table td.r table th.r]

```

## parseStyles?(go)

> styleRegex := regexp.MustCompile

`parseStyles` takes the second output of the RegExp match in [ParseCSS](./#parsecssgo) and splits it up using this RegExp:

> styleMap := make(map[string]string)
> for \_, match := range matches {
> propName := strings.TrimSpace(match[1])
> propValue := strings.TrimSpace(match[2])
> styleMap[propName] = propValue
> }

It then takes the split styles and inserts them into a `map[string]string`.

### parseStyles Example

```go

selectorBlock := `text-align: center;
color: red;`

parseStyles(selectorBlock)

// Output
map[string]string=map[text-align:center color:red]

```

## ParseStyleAttribute?(go)

> inline := parser.ParseStyleAttribute(n.GetAttribute("style") + ";")

`ParseStyleAttribute` is for parsing inline styles from elements in the html document on the inital load. It is also used to parse the local styles applied by the "script" via the `.style` attribute. It will only be applied to a `element.Node`'s local styles and will not be add to the global stylesheets. It is used with the `cstyle.GetStyles` function that is ran on every cycle.

### ParseStyleAttribute Example

```go

styleAttribute := "color:#f8f8f2;background-color:#272822;"

ParseStyleAttribute(styleAttribute)

//Output
map[string]string=map[color:#f8f8f2 background-color:#272822]

```

## removeComments?(go)

<{./main.go}>
<{../cstyle/main.go}>
