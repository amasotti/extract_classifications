/* Read and parse xml (format mods) from sru.k10plus/gvk */

package main

import (
	"encoding/json"
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

// FinalClassification is the struct for storing the extracted classifications
type FinalClassification struct {
	Lcc     []string `json:"lcc"`
	Ddc     []string `json:"ddc"`
	Bisacsh []string `json:"bisacsh"`
	BISAC   []string `json:"BISAC"`
	Bkl     []string `json:"bkl"`
	Rvk     []string `json:"rvk"`
}

/*********************************/
//	MAIN FUNCTION
/*********************************/
func main() {
	var verbose bool = false

	// Read the file
	xmlFile, err := os.Open("dostoevsky.xml") // TODO: implement http GET to automatically retrieve data from the sru address

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
	recs := parsedXML.Records.Record

	//finalClass := new(FinalClassification)
	var finalSubjs []string
	var finalClass = new(FinalClassification)

	for i := 0; i < len(recs); i++ {

		// Loop over schlagwoerter
		schlagwoerter := recs[i].Mods.Infos.Schlagwoerter
		if len(schlagwoerter) != 0 {
			for j := 0; j < len(schlagwoerter); j++ {

				var (
					subject    = schlagwoerter[j]
					schlagwort = subject.Schlagwort
					authority  = "-"
				)
				if subject.Authority != "" {
					authority = subject.Authority
				}
				// Save "Schlagwörter" independently of the authority
				finalSubjs = append(finalSubjs, schlagwort)

				if verbose {
					fmt.Println("\nTopic/Schlagwort found: ")
					fmt.Println("----------------------------------------")
					fmt.Println("Subject: " + schlagwort)
					if authority != "-" {
						fmt.Println("Assigned by: " + authority)
					}
				}
			}
		}

		// Loop over Classifications
		classes := recs[i].Mods.Infos.Subjects
		if len(classes) != 0 {
			if verbose {
				fmt.Println("\nSubject found: ")
				fmt.Println("--------------------------")
			}
			for j := 0; j < len(classes); j++ {
				switch classes[j].Authority {

				case "lcc":
					finalClass.Lcc = append(finalClass.Lcc, classes[j].Subject)
				case "ddc":
					finalClass.Ddc = append(finalClass.Ddc, classes[j].Subject)
				case "BISAC":
					finalClass.BISAC = append(finalClass.BISAC, classes[j].Subject)
				case "bisacsh":
					finalClass.Bisacsh = append(finalClass.Bisacsh, classes[j].Subject)
				case "rvk":
					finalClass.Rvk = append(finalClass.Rvk, classes[j].Subject)
				case "bkl":
					finalClass.Bkl = append(finalClass.Bkl, classes[j].Subject)
				}

				if verbose {
					fmt.Println("Authority: '" + classes[j].Authority + "'\tValue: " + classes[j].Subject)
				}
			}
		}
	}

	// print results
	//fmt.Println(finalSubjs)
	//fmt.Println(finalClass)

	//Save to json
	file, err := json.MarshalIndent(finalClass, "", "\t")
	err = ioutil.WriteFile("outputs/testClasses.json", file, 0644)

}
