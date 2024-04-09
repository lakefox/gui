# CStyle Plugins

Plugins add the ability to choose what parts of the HTML/CSS spec you add to your application. If you are trying to keep compile sizes small you can remove as many as you need to reach your target size. Here we will go over the basics of how they work and how to use them.

```go
type Plugin struct {
	Styles  map[string]string
	Level   int
	Handler func(*element.Node, *map[string]element.State)
}
```

A plugin is a struct defined in CStyle, in contains three properties:

-   Styles
    -   A map with CSS properties and values that the plugin should match on. There is also support for wildcard properties by setting the value to "\*"
-   Level
    -   Level of priority in which to execute
    -   All library level should be between 0 and 2
        -   All others should be greater than 2
-   Handler
    -   Callback function that provides a pointer to a element that matches the properties of Styles

## AddPlugin?(go)

Add Plugin is the CStyle function that add the plugin to the top level cstyle.Plugins array where it is used within the `ComputeNodeStyle` function.

### Usage

```go
css.AddPlugin(block.Init())
```

The first step in processing the plugins before running there handler functions is to sort them by their levels. We need to sort the plugins by their level because in the example of flex, it relys on a parent element being block position so it can compute styles based of the positioning of its parent elements. If flex was ran before block then it would have nothing to build apon. This is also the reason that if you are building a custom plugin it is reccomended to keep the level above 2 as anything after 2 will have the assumed styles to build apon.

> // Sorting the array by the Level field
> sort.Slice(plugins, func(i, j int) bool {
> return plugins[i].Level < plugins[j].Level
> })

After we have the sorted plugins we can check if the current element matches the `Styles` of the plugin. The matching is a all of nothing matching system, if one property is missing then the plugin wil not be ran. If it does match then a pointer to the `element.Node` (n) is passed to the handler.

> for \_, v := range plugins {
> matches := true
> for name, value := range v.Styles {
> if styleMap[name] != value && !(value == "\*") {
> matches = false
> }
> }
> if matches {
> v.Handler(n)
> }
> }

## plugins/block

Here is the code for the block styling plugin, it is recommended to keep this standard format.

> All other plugins can be found in the cstyle/plugins folder

<{./block/main.go}>
<{../main.go}>
