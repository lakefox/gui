# CStyle

CStyle is a single pass style computer.

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
