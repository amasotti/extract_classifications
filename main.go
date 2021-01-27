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
	Schlagwort string   `xml:"topic"`
}

// Subjects finds the categories (codes) and the assigning authority
type Subjects struct {
	XMLName   xml.Name
	Authority string `xml:"authority,attr"`
	Subject   string `xml:",innerxml"`
}

/*********************************/
// 			AUXILIARY FUNCTIONS
/*********************************/

// ReadMarshalXML reads a xml file and returns the marshalled content
func ReadMarshalXML(path string) *Root {
	// Read the file
	xmlFile, err := os.Open(path) // TODO: implement http GET to automatically retrieve data from the sru address

	if err != nil {
		log.Println(err)
	}
	defer xmlFile.Close()

	// Initialize the XML categories
	parsedXML := new(Root)
	contents, err := ioutil.ReadAll(xmlFile)

	if err == nil {
		xml.Unmarshal(contents, &parsedXML)
	} else {
		log.Println(err)
	}

	return parsedXML

}

// ExtractClassifications loops over schlagwörter und classifications and returns string slices
func ExtractClassifications(results []SingleResult, verbose bool) (finalClasses FinalClassification, finalSubjs []string) {
	var subjects []string
	var classes = new(FinalClassification)

	for i := 0; i < len(results); i++ {
		schlagwoerter := results[i].Mods.Infos.Schlagwoerter
		classifications := results[i].Mods.Infos.Subjects

		// over schlagwoerter
		for j := 0; j < len(schlagwoerter); j++ {
			subjects = append(subjects, schlagwoerter[j].Schlagwort)
			if verbose {
				fmt.Println("\nTopic/Schlagwort found: ")
				fmt.Println("----------------------------------------")
				fmt.Println("Subject: " + schlagwoerter[j].Schlagwort)
			}
		}
		// over explicit subjects
		if len(classifications) == 0 {
			log.Println("Item without classification --- skipping it")
		} else {
			if verbose {
				fmt.Println("\nClassification found: ")
				fmt.Println("--------------------------")
			}
			for j := 0; j < len(classifications); j++ {
				switch classifications[j].Authority {
				case "lcc":
					classes.Lcc = append(classes.Lcc, classifications[j].Subject)
				case "ddc":
					classes.Ddc = append(classes.Ddc, classifications[j].Subject)
				case "BISAC":
					classes.BISAC = append(classes.BISAC, classifications[j].Subject)
				case "bisacsh":
					classes.Bisacsh = append(classes.Bisacsh, classifications[j].Subject)
				case "rvk":
					classes.Rvk = append(classes.Rvk, classifications[j].Subject)
				case "bkl":
					classes.Bkl = append(classes.Bkl, classifications[j].Subject)
				}
				if verbose {
					fmt.Println("Authority: '" + classifications[j].Authority + "'\tValue: " + classifications[j].Subject)
				}
			}
		}

	}

	return finalClasses, finalSubjs
}

func ClearDuplicates(ToBeCleaned []string) []string {
	key := make(map[string]bool)
	list := []string{}

	for _, entry := range ToBeCleaned {
		if _, value := key[entry]; !value {
			key[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/*********************************/
//	UTIL FUNCTIONS AND TYPES
/*********************************/

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
	parsedXML := ReadMarshalXML("dostoevsky.xml")

	// Extract list of records
	recs := parsedXML.Records.Record

	finalClass, finalSubjs := ExtractClassifications(recs, false)
	// print results
	//fmt.Println(finalSubjs)
	//fmt.Println(finalClass)

	//Save to json
	file, _ := json.MarshalIndent(finalClass, "", "\t")
	_ = ioutil.WriteFile("outputs/testClasses.json", file, 0644)

	file, _ = json.MarshalIndent(finalSubjs, "", " ")
	_ = ioutil.WriteFile("outputs/testSubjects.json", file, 0644)

}
