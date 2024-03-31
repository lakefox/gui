# Parser

Parser is the CSS parser for this project, it is made up of two primary functions `ParseCSS`, `ParseStyleAttribute`, and a few other functions designed to help with the parsing.

## ParseCSS?(go)

`ParseCSS` is the function for reading CSS files. It is a RegExp based parser which it converts CSS definitions into a 2d map of strings. If you want to convert the values into absoulte values there are helper functions in the [utils](/utils) documentation.

## parseSelectors?(go)

## parseStyles?(go)

## ParseStyleAttribute?(go)

## removeComments?(go)

<{./main.go}>
