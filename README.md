# Musicians Repertoire and Setlist tool

* Generates set lists
* Generates cheat sheets

## Status of the Project

This is currently just a PoC. I am using it for my musical hobby and develop it on the go.

## Dependencies

- Cheat sheets can be designed in the odt format. LibreOffice is used to generate pdf.
- To convert html pdf, a chrome (tested with chromium) browser is required.

## Repertoire Directory Structure

By convention the Repertoire is organised in the following [structure](test/Repertoire):

```
 Band
├──  gigs
│   └──  'Grand Ole Opry.md'
├──  songs
│   ├──  'Frankie and Johnnie.odt'
│   ├──  'On the Alamo.pdf'
│   └──  README.md
└──  Repertoire.md 
```

