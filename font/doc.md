# Font

Font rasterizes all text in a document and handles the wrapping and positionion of text within a texture.

```mermaid
flowchart LR;
    LoadFont-->State.Text.Font;
    State.Text.Font-->State;
    State-->Render;
```

## GetFontPath?(go)

## tryLoadSystemFont?(go)

## sortByLength?(go)

## GetFontSize?(go)

## LoadFont?(go)

## MeasureText?(go)

## MeasureSpace?(go)

## MeasureLongest?(go)

## getSystemFonts?(go)

## getWindowsFontPaths?(go)

## getMacFontPaths?(go)

## getLinuxFontPaths?(go)

## getFontsRecursively?(go)

## Render?(go)

## drawString?(go)

## wrap?(go)

## drawLine?(go)

## getLines?(go)

<{./main.go}>
