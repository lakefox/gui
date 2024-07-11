# CStyle

CStyle is a single pass style computer.

> WARNING: if you are building a transformer, use c.QuickStyles to get the styles of the element as it speed it up by over 50%. However, QuickStyles does not add any style sheet (master.css) styles to the tags, you will have to add them manually.

## StyleSheet?(go)

## StyleTag?(go)

## GetStyles?(go)

## AddPlugin?(go)

See [/cstyle/plugins/](/cstyle/plugins/)

## ComputeNodeStyle?(go)

> if !utils.ChildrenHaveText(n)
> The utils.ChildrenHaveText function is called here instead of checking if the node directly has text is because in the case below

```html
<em><b>Text</b></em>
```

The `em` element does not have text but it has a element with text insid, but the element will still need to be rendered as text.

## parseBorderShorthand?(go)

## CompleteBorder?(go)

## genTextNode?(go)

<{./main.go}>
