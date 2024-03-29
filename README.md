# Musicians Repertoire and Setlist tool

* Generates set lists
* Generates cheat sheets

## Status of the Project

This is currently just a PoC. I am using it for my musical hobby and develop it on the go.

## Dependencies

- Cheat sheets can be designed in the odt format. LibreOffice is used to generate pdf.
- To convert html pdf, a chrome (tested with chromium) browser is required.

## Repertoire

By convention the Repertoire is organised in the following [directory structure](test/Repertoire):

The repertoire is partitioned into substructures for example to support multiple bands.

```
Band
├── Gigs
│   └── 'Grand Ole Opry.md'
├── Songs
│   ├── 'Frankie and Johnnie.odt'
│   ├── 'On the Alamo.pdf'
│   └── README.md
└── Repertoire.md
```

### Repertoire.md

Metadata is maintained in a Markdown [GFM table](https://github.github.com/gfm/#tables-extension-) on the top level of the
Repertoire.md file. See [example]{test/Repertoire/Band/Repertoire.md).

The Table must have a header row. The only mandatory column is the `Title` column which is used to refer to the song titles
for generating both cheat `sheet` and set `list`.

Optional columns are used to generate output.

| Column      | Type      | Used by command |
|-------------|-----------|-----------------|
| Title       | Mandatory | list, cheat     |
| Year        | Optional  | list            |
| Description | Optional  | list            |


### Gigs

Each gig is maintained in a Markdown file within the [Gigs](test/Repertoire/Band/Gigs) subdirectory.
The name of the Markdown file is the name of the gig. Each song title is listed on top level of the Markdown file as an
[unordered list](https://www.markdownguide.org/basic-syntax/#unordered-lists) element.

### Songs

Each song is maintained in a .pdf file within the [Songs](test/Repertoire/Band/Songs) subdirectory.
Optionally a .pdf can also be generated out of an Open Document (.odt) file.
The filename must the same as the song title within the Gig Markdown file in order to be picked up by the cheat `sheet` 
generator.
