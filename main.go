/* Read and parse xml (format mods) from sru.k10plus/gvk */

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

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

type Schlagwoerter struct {
	XMLName    xml.Name `xml:"subject"`
	Authority  string   `xml:"authority,attr"`
	Schlagwort string   `xml:"topic"`
}

type Subjects struct {
	XMLName   xml.Name
	Authority string `xml:"authority,attr"`
	Subject   string `xml:",innerxml"`
}

func main() {
	// Read the file
	xmlFile, err := os.Open("dostoevsky.xml")

	if err != nil {
		log.Println(err)
	}

	defer xmlFile.Close()

	contents, err := ioutil.ReadAll(xmlFile)
	parsedXML := new(Root)

	if err != nil {
		log.Println(err)
	} else {
		xml.Unmarshal(contents, &parsedXML)
	}
	//fmt.Println()
	recs := parsedXML.Records.Record
	for i := 0; i < len(recs); i++ {

		// Loop over schlagwoerter
		schlagwoerter := recs[i].Mods.Infos.Schlagwoerter
		if len(schlagwoerter) != 0 {
			fmt.Println("\n\n\nTopic/Schlagwort found: ")
			fmt.Println("----------------------------------------")
			for j := 0; j < len(schlagwoerter); j++ {

				var (
					subject    = schlagwoerter[j]
					schalgwort = subject.Schlagwort
					authority  = "-"
				)
				if subject.Authority != "" {
					authority = subject.Authority
				}

				fmt.Println("Subject: " + schalgwort)
				if authority != "-" {
					fmt.Println("Assigned by: " + authority)
				}
			}
		}

		// Loop over Classifications
		classes := recs[i].Mods.Infos.Subjects
		if len(classes) != 0 {
			fmt.Println("\nSubject found: ")
			fmt.Println("--------------------------")
			for j := 0; j < len(classes); j++ {
				fmt.Println("Authority: '" + classes[j].Authority + "'\tValue: " + classes[j].Subject)
			}
		}
		fmt.Println("\n\n")
	}

}
