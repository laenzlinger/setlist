# Musicians Repertoire and Setlist tool

* Generates set lists
* Generates cheat sheets

## Status of the Project

This is currently just a PoC. I am using it for my musical hobby
and develop it on the as i need more features.
Therefore - for all 0.X versions - the command line interface might change,
there is no guarantee fro backwards compatibility yet.

## CLI documentation

Run `setlist --help` to or check the [Markdown CLI docs](./docs/setlist.md)

## Repertoire structure

By convention the Repertoire is organised in the following [directory structure](test/Repertoire):

The repertoire is by default partitioned into substructures for
example to support multiple bands.

```bash
Repertoire
├── Band
│   ├── Gigs
│   │   └── 'Grand Ole Opry.md'
│   ├── Songs
│   │   ├── 'Frankie and Johnnie.odt'
│   │   ├── 'Her Song.md'
│   │   └── 'On the Alamo.pdf'
│   ├── README.md
│   └── Repertoire.md
├── .gitignore
└── .setlist
```

### Repertoire.md

Metadata is maintained in a Markdown
[GFM table](https://github.github.com/gfm/#tables-extension-) on the top level
of the Repertoire.md file. See [example](test/Repertoire/Band/Repertoire.md).

The Table must have a header row. The only mandatory column is the
`Title` column which is used to refer to the song titles
for both `generte sheet` and  `generate list`.

Optional columns used (by default) generate output.

| Column      | Type      | Used by command    |
|-------------|-----------|--------------------|
| Title       | Mandatory | list, cheat, suisa |
| Year        | Optional  | list, suisa        |
| Description | Optional  | list               |
| Arranger    | Optional  | suisa              |
| Composer    | Optional  | suisa              |
| Duration    | Optional  | suisa              |

The output columns can be selected by the `--include-columns` flag,
but the order or the columns is defined by the input Repertoire.md

### Gigs

Each gig is maintained in a Markdown file within the [Gigs](test/Repertoire/Band/Gigs) subdirectory.
The name of the Markdown file is the name of the gig. Each song title is listed on top level of the Markdown file as an
[unordered list](https://www.markdownguide.org/basic-syntax/#unordered-lists) element.

### Songs

Each song is maintained in a .pdf file within the [Songs](test/Repertoire/Band/Songs) subdirectory.
Optionally a .pdf can also be generated out of an Open Document (.odt) or a Markdwon (.md) file.
The filename must the same as the song title within the Gig Markdown file in order to be picked up by the cheat `sheet`
generator.

## Dependencies

* Cheat sheets can be designed in the odt format. LibreOffice is used to generate pdf.
* To convert html pdf, a chrome (tested with chromium) browser is required.

### Docker image

In case you don't want to install the dependencies locally, you can use the [docker image](https://github.com/laenzlinger/setlist/pkgs/container/setlist):

Example:

```bash
docker run --rm --user "$(id -u)":"$(id -g)" -v $(pwd):/repertoire ghcr.io/laenzlinger/setlist --help
```
