# Selector

Selector is a implementation of JavaScripts querySelector. It is split between two files this file and the `element` package to prevent circular dependancys, however this document will be the source for it. The best way to explain how this works is to start in the `element` package with the `querySelector` method and then take a look at the parts that make it up.

> func (n *Node) QuerySelector(selectString string) *Node {

## QuerySelector

| Arguments             | Description                     |
| --------------------- | ------------------------------- |
| node \*element.Node   | Target \*element.Node           |
| selectString string   | CSS querySelector string        |
| return \*element.Node | element.Node matching the query |

`QuerySelector` works almost the same as JavaScripts [querySelector method](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_selectors) with a far limited scope. After a document is loaded from a HTML file it is compiled into `element.Node`'s which is a custom implementation of `net/html.Node`. The reason the `net/html` node is not used is it has already defined features that stray away from JavaScripts DOM.

> if TestSelector(selectString, n) {return n}

To start out, we check the current element to see if the `selectString` matches the `element.Node` we called the method on using the [`TestSelector`](./#testselector) function. If it does we can end the function there and return itself. If it does not we can continue and check its children. We do this process recursively to simplify the code.

> if cr.Properties.Id != "" {return cr}

We also do a check to see if the `element.Node.Properties.Id` has been assigned. This is a importaint step as this id is the the `#id` used in html but a unqiue id generated at run time to be used as a internal reference. If it has not been assigned then the element does not exist.

> func TestSelector(selectString string, n \*Node) bool {

## TestSelector

| Arguments           | Description                                     |
| ------------------- | ----------------------------------------------- |
| selectString string | CSS querySelector string                        |
| node \*element.Node | Target \*element.Node                           |
| return bool         | returns true if the selector matches the string |

`TestSelector` is the foundation of the [`QuerySelector`](./#queryselector) and [`QuerySelectorAll`](./#queryselectorall) as seen above.

> parts := strings.Split(selectString, ">")

It first starts off by splitting the `selectString` in to parts divided by `>` this is becuase when you have a selector like `blockquote > p` you need to start at the first level (`p`) to compare the current node to see if you will need to continue to check the parents of the element with the next selector.

> s := []string{}
> if n.Properties.Focusable {
> if n.Properties.Focused {
> s = append(s, ":focus")
> }
> }
> classes := n.ClassList.Classes
> for \_, v := range classes {
> s = append(s, "."+v)
> }

Then we need to build the selectors, so we start by creating an array to store them in (`s`) and we check to see if the element is focusable and if the element is focused. If so we add the `:focus` selector to the list. This is important because when targeting a `:focus`ed element with a querySelector that is the text that is past. We then do the same for classes.

> selectors := selector.GetCSSSelectors(n.Properties.Node, s)

Next we use the [`GetCSSSelectors`](./#getcssselectors) method in this package to generate any selectors assigned to the `net/html` Node.

> if n.Id != "" {
> selectors = append(selectors, "#"+n.Id)
> }

Then we add the id to the array to complete the current Nodes selectors.

> part := selector.SplitSelector(strings.TrimSpace(parts[len(parts)-1]))
> has := selector.Contains(part, selectors)

After we have the current Nodes selectors we can use the [SplitSelector](./#splitselector) and [Contains](./#contains) methods to process the passed query (selectString) and compare the two arrays.

> func GetCSSSelectors(node \*html.Node, selectors []string) []string {

## GetCSSSelectors

| Arguments          | Description          |
| ------------------ | -------------------- |
| node \*html.Node   | Target net/html Node |
| selectors []string | Previous Selctors    |
| return []string    | Output of selectors  |

`GetCSSSelectors` purpose is to generate all possible selectors for a `net/html` Node. It is used inside of the element package interally to the [`TestSelector`](./#testselector) function. It does this buy taking the classes, id's, and attributes and creating an array of their string equalivents (.class, #id, and [value="somevalue"]).

> func SplitSelector(s string) []string {

## SplitSelector

| Arguments       | Description         |
| --------------- | ------------------- |
| s string        | Selector string     |
| return []string | Output of selectors |

`SplitSelector` works by simply spliting a CSS selector into it's individual parts see below for an example:

```go
func main() {
	fmt.Println(SplitSelector("p.text[name='first']"))
}
```

Result

```text
[p .text [name='first']]
```

> func Contains(selector []string, node []string) bool {

## Contains

| Arguments         | Description                                 |
| ----------------- | ------------------------------------------- |
| selector []string | Array of selectors from the target selector |
| node []string     | Array of selectors from the target element  |
| return bool       | boolean value                               |

`Contains` compares two arrays of selectors, the first argument is the array of the selector that will be use to detirmine if the Node is a match or not. The second argument is the selecter of the targeted Node, the Node need to have all of the selectors of the `selector` array, however it can have additional selectors and it will still match.

```go
package selector

import (
 "slices"
 "strings"

 "golang.org/x/net/html"
)

func GetCSSSelectors(node *html.Node, selectors []string) []string {
 if node.Type == html.ElementNode {
  selectors = append(selectors, node.Data)
  for _, attr := range node.Attr {
   if attr.Key == "class" {
    classes := strings.Split(attr.Val, " ")
    for _, class := range classes {
     selectors = append(selectors, "."+class)
    }
   } else if attr.Key == "id" {
    selectors = append(selectors, "#"+attr.Val)
   } else {
    selectors = append(selectors, "["+attr.Key+"=\""+attr.Val+"\"]")
   }
  }
 }

 return selectors
}

func SplitSelector(s string) []string {
 var result []string
 var current string

 for _, char := range s {
  switch char {
  case '.', '#', '[', ']', ':':
   if current != "" {
    if string(char) == "]" {
     current += string(char)
    }
    result = append(result, current)
   }
   current = ""
   if string(char) != "]" {
    current += string(char)
   }
  default:
   current += string(char)
  }
 }

 if current != "" && current != "]" {
  result = append(result, current)
 }

 return result
}

func Contains(selector []string, node []string) bool {
 has := true
 for _, s := range selector {
  if !slices.Contains(node, s) {
   has = false
  }
 }
 return has
}

```

```go
func (n *Node) QuerySelectorAll(selectString string) *[]*Node {
 results := []*Node{}
 if TestSelector(selectString, n) {
  results = append(results, n)
 }

 for i := range n.Children {
  el := &n.Children[i]
  cr := el.QuerySelectorAll(selectString)
  if len(*cr) > 0 {
   results = append(results, *cr...)
  }
 }
 return &results
}

func (n *Node) QuerySelector(selectString string) *Node {
 if TestSelector(selectString, n) {
  return n
 }

 for i := range n.Children {
  el := &n.Children[i]
  cr := el.QuerySelector(selectString)
  if cr.Properties.Id != "" {
   return cr
  }
 }

 return &Node{}
}

func TestSelector(selectString string, n *Node) bool {
 parts := strings.Split(selectString, ">")

 s := []string{}
 if n.Properties.Focusable {
  if n.Properties.Focused {
   s = append(s, ":focus")
  }
 }

 classes := n.ClassList.Classes

 for _, v := range classes {
  s = append(s, "."+v)
 }
 // fmt.Println(n.Properties.Node)
 selectors := selector.GetCSSSelectors(n.Properties.Node, s)
 if n.Id != "" {
  selectors = append(selectors, "#"+n.Id)
 }

 part := selector.SplitSelector(strings.TrimSpace(parts[len(parts)-1]))

 has := selector.Contains(part, selectors)

 if len(parts) == 1 || !has {
  return has
 } else {
  return TestSelector(strings.Join(parts[0:len(parts)-1], ">"), n.Parent)
 }
}
```

<script src="/plugins/sidebyside.js"></script>
