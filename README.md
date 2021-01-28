# Classify (K10+)

**Work in progress**


The idea is to create a small flexible tool to extract library classifications (rvk, bkl, ddc, lcc) from the GVK Catalogue given a query (keyword, title, author etc...)

+ [Sru](https://wiki.k10plus.de/display/K10PLUS/SRU) address

      https://sru.k10plus.de/gvk


The SRU protocol supports several output formats. This scritp works (at moment at least) only with the [MODS](https://en.wikipedia.org/wiki/Metadata_Object_Description_Schema) format.

This is a personal project, created mostly for fun, to learn Go and to make some of my daily-tasks easier. It is not perfect at all
but feel free to play with the code or suggest me improvements. Install the code (with correctly installed GO environment) in your directory using

        go get github.com/amasotti/k10_classify

Build the binary with

        go build main.go

Play with it ;) 

#### Usage:

```
Usage:
$ klassify_10.exe  + parameters (at least -q):

    -k               query Key (bkl, tit, per, slw (default: all)
    -q               query String (the text to search)
    -v               verbose (bool) prints everything on the console, not recommended (debug purpose only)
    -s               save (bool) if true saves the retrieved xml and the json output files (default : false)
    -n               number of Results (integer) specifiers how many results will be printed on the console
    -p               path (string) which subdirectory of the current working directory will be used to save the files (works only if -s is set)             
    -m               maxResult how many entries should be retrieved

```

## Subjects Headings

### Authority:

+ **bkl** : [Basisklassifikation](https://www.gbv.de/bibliotheken/verbundbibliotheken/02Verbund/01Erschliessung/02Richtlinien/05Basisklassifikation/index)
+ **bisacsh** : [BISAC Subject Headings](https://bisg.org/page/bisacedition)
+ **bicss** : [BIC Subject Categories](https://bic.org.uk/files/pdfs/101201%20bic2.1%20complete%20rev.pdf)
+ **ddc** : Dewey Class or Division
+ **fid** : Fachinformationsdienst
+ **lcsh** : [Library of Congress Subject Headings](https://id.loc.gov/vocabulary/subjectSchemes/bisacsh.html)
+ **rvk** : [Regensburger Verbundsklassifikation](https://rvk.uni-regensburg.de/regensburger-verbundklassifikation-online)

Check [loc.gov](https://www.loc.gov/standards/sourcelist/subject.html) for a complete list of authority codes.


### Example of raw data in MODS format:

see [test example](https://github.com/amasotti/k10_classify/blob/main/testXML.xml)


## Output examples:

### Json output example

#### Classification

```json
{
	"lcc": {
		"AG": 1,
		"B": 2,
		"B1-5802": 1,
                ...
	},
	"ddc": {
		"000": 6,
		"001": 1,
		"001.51": 1,
                ...
	},
	"bisacsh": {},
	"BISAC": {
		"ART 043000": 1,
		"ART015000": 2,
		"ART015100": 1,
		          ...
	},
	"bkl": {
		"02.01": 3,
		"02.02": 1,
		"02.13": 4,
		          ...
	},
	"rvk": {
		"AK 18000": 2,
		"AK 39500": 2,
		"AK 39540": 1,
                  ...
	}
}

```




### Console output

<pre>
$ go run  main.go -k swl -q "Umberto Eco" -m 580

2021/01/28 16:21:14 GET https://sru.gbv.de/gvk?version=1.1&operation=searchRetrieve&query=pica.all=Umberto%20Eco&recordSchema=mods&maximumRecords=580

<strong>Number of results: 580</strong>
------------------------------------------------------------------
SUBJECT HEADINGS
------------------------------------------------------------------

Printing the 5 most common Subject Headings for your query:

Geschichte : 54
Literatur : 39
Semiotik : 28
Criticism and interpretation : 27
Kunst : 26
20th century : 25
History and criticism : 24


CLASSIFICATIONS
------------------------------------------------------------------
The 7 most common Lcc classifications for your query

PQ4865.C6 : 41
PQ : 18
PN : 10
P99 : 7
PN241 : 5
D : 3
BH81 : 3

The 7 most common Ddc classifications for your query

850 : 40
853 : 24
800 : 19
850 B : 19
853.914 : 15
700 : 14
809 : 14

The 7 most common BISAC classifications for your query

LIT004200 : 2
LIT006000 : 2
LIT 006000 : 2
LIT000000 : 2
ART015000 : 2
LIT 004200 : 2
POL 019000 : 1

The 7 most common Bkl classifications for your query

18.27 : 118
17.97 : 79
18.00 : 44
20.06 : 22
17.73 : 15
17.08 : 15
08.41 : 11

The 7 most common Rvk classifications for your query

IV 25480 : 108
IV 25481 : 35
LH 61040 : 15
CC 6900 : 14
AK 39580 : 8
ER 730 : 8
EC 1070 : 7

</pre>


## ToDO

+ <del>Use ```net/http``` to download the xml from the web instead of loading it manually (or possibly leave both options open)</del>
+ <del>Allow multiple keywords separated by AND</del>
      slw : "Umberto Eco AND Semiotik"
+ <del>Allow for more than one query key (at moment only slw *Schlagwort* possible)</del>
    + Allow more than one key at the same time
+ Give a closer look at the Go idiomatic way of initializing new structs and saving values passing them to the pointers
+ <del>Implement or use a counter to see if some keywords (Schlagwörter) are used more than other</del>
  + <del>Still to do : sort the mapping</del> 


