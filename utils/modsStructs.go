package utils

import "encoding/xml"

/*********************************/
//			XML STRUCTS
/*********************************/

// TODO: Posssibly convert the data array into array of pointers (should be faster)

// Root tag
type Root struct {
	XMLName xml.Name
	Records Records `xml:"records"`
}

// Records list of found items
type Records struct { // list of records
	XMLName xml.Name
	Record  []SingleResult `xml:"record"`
}

// SingleResult : main node for single items, contains all the relevant informations
type SingleResult struct {
	Position int  `xml:"recordPosition"` // index
	Mods     Mods `xml:"recordData"`
}

// Mods Format tag, contains info about the version of the format and the URL to the specs
type Mods struct {
	Infos Info `xml:"mods"`
}

// Info This is the tag of greatest interest for us; it contains all the importants subtags
type Info struct {
	Schlagwoerter []Schlagwoerter `xml:"subject"`
	Subjects      []Subjects      `xml:"classification"`
}

// Schlagwoerter finds the subjects and the assigning authority
type Schlagwoerter struct {
	XMLName    xml.Name `xml:"subject"`
	Authority  string   `xml:"authority,attr"`
	SchlagwortGeneral string   `xml:"topic"`
	SchlagwortTemporal string   `xml:"temporal"`
	SchlagwortGeographic string   `xml:"geographic"`
}

// Subjects finds the categories (codes) and the assigning authority
type Subjects struct {
	XMLName   xml.Name
	Authority string `xml:"authority,attr"`
	Subject   string `xml:",innerxml"`
}
