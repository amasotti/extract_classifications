/* Read and parse xml (format mods) from sru.k10plus/gvk */

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Define structs : requires knowledge of the specs of the xml format.
// TODO: Posssibly convert the data array into pointers array
type Root struct {
	XMLName xml.Name
	Records Records `xml:"records"` // root
}

type Records struct { // list of records
	XMLName xml.Name
	Record  []SingleResult `xml:"record"`
}

type SingleResult struct { // each result taken alone
	Position int  `xml:"recordPosition"` // index
	Mods     Mods `xml:"recordData"`
}

type Mods struct { // Format tag, contains info about the version of the format and the URL to the specs
	RecordData RecordData `xml:"mods"`
}

type RecordData struct { // This is the tag of greatest interest for us; it contains all the importants subtags
	Title  Title    `xml:"titleInfo"`
	Person []Person `xml:"name"`
}

type Title struct {
	XMLName  xml.Name
	Titel    string `xml:"title"`
	Subtitle string `xml:"subTitle"`
}

type Person struct {
	Name string `xml:"namePart"`
	Role string `xml:"role>roleTerm"`
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

		if i < 10 {
			titelField := recs[i].Mods.RecordData.Title
			titel := titelField.Titel
			subtitel := titelField.Subtitle
			fmt.Println(titel + "\n" + subtitel + "\n")

			autor := recs[i].Mods.RecordData.Person
			for j := 0; j < len(autor); j++ {
				fmt.Println("Name: " + autor[j].Name)
				fmt.Println("Funktion: " + autor[j].Role + "\n\n")
			}
		}

	}

}
