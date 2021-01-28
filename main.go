/* Read and parse xml (format mods) from sru.k10plus/gvk */

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
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

/*********************************/
// 			AUXILIARY FUNCTIONS
/*********************************/

// ReadMarshalXML reads a xml file and returns the marshalled content
func ReadMarshalXML(path string) *Root {
	// Read the file
	xmlFile, err := os.Open(path)

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
func ExtractClassifications(results []SingleResult, verbose bool) (OrderedClassification, map[string]int) {
	var subjects []string
	var classes = new(FinalClassification)

	for i := 0; i < len(results); i++ {
		schlagwoerter := results[i].Mods.Infos.Schlagwoerter
		classifications := results[i].Mods.Infos.Subjects

		// over schlagwoerter
		for j := 0; j < len(schlagwoerter); j++ {
			if schlagwoerter[j].SchlagwortGeneral != "" {
				if verbose {
					fmt.Println("\nTopic/Schlagwort found: ")
					fmt.Println("----------------------------------------")
					fmt.Println("Subject: " + schlagwoerter[j].SchlagwortGeneral)
				}
				subjects = append(subjects, schlagwoerter[j].SchlagwortGeneral)
			}
			if schlagwoerter[j].SchlagwortGeographic != "" {
				subjects = append(subjects, schlagwoerter[j].SchlagwortGeographic)
			}
			if schlagwoerter[j].SchlagwortTemporal != "" {
				subjects = append(subjects, schlagwoerter[j].SchlagwortTemporal)
			}
		}
		// over explicit subjects
		if len(classifications) != 0 {
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
	Countedsubjects := CountUnique(subjects)
	CountedClassification := orderClassification(*classes)
	return CountedClassification, Countedsubjects
}
/*********************************/
//				COUNTERS
/*********************************/
// CountUnique counts the number of unique elem in a list and returns a map [elemen] : [occurrences]
func CountUnique(list []string) map[string]int {

	counter := make(map[string]int)
	check := make(map[string]bool)

	for _, el := range list {
		_, counted := check[el]
		if !counted {
			counter[el] = 1
			check[el] = true
		} else {
			counter[el] += 1
		}
	}
	return counter
}

// orderClassification transforms a list of classifications into a counted list of classifications
func orderClassification(c FinalClassification) OrderedClassification {
	sorted := OrderedClassification{Lcc: CountUnique(c.Lcc), Ddc: CountUnique(c.Ddc), Bisacsh: CountUnique(c.Bisacsh), BISAC: CountUnique(c.BISAC), Bkl: CountUnique(c.Bkl), Rvk: CountUnique(c.Rvk)}

	return sorted
}

/*********************************/
//	 HTTP FUNCTIONS
/*********************************/

func getXML(queryText, queryIndex, path string, save bool) *Root {
	body := sendRequest(queryText, queryIndex)
	if save {
		_ = ioutil.WriteFile(path, body, 0666)
	}

	parsedXML := new(Root) // TODO: Avoid new keyword
	xml.Unmarshal(body, &parsedXML)
	return parsedXML
}

// sendRequest sends the query to the sru address
func sendRequest(queryText, queryIndex string) (body []byte) {
	request := buildQuery(queryText, queryIndex)

	//request := fmt.Sprintf(address, url.PathEscape(queryText))
	log.Println("GET", request)
	response, err := http.Get(request)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return body
}

func buildQuery(queryText, queryIndex string) (request string) {

	switch queryIndex {
	case "slw":
		fallthrough
	case "tit":
		fallthrough
	case "bkl":
		fallthrough
	case "per":
		fullQuery := analyzeInput(queryText, queryIndex)
		request = fmt.Sprintf("https://sru.gbv.de/gvk?version=1.1&operation=searchRetrieve&query=%s&recordSchema=mods&maximumRecords=500",fullQuery)
	default: // Use the full text search
		GeneralQuery := analyzeInput(queryText, "all")
		request = fmt.Sprintf("https://sru.gbv.de/gvk?version=1.1&operation=searchRetrieve&query=%s&recordSchema=mods&maximumRecords=500",GeneralQuery)
	}
	return
}
// analyzeInput looks at the string given and formats this, ready to be used in the a query
func analyzeInput(query, key string) string {
	splittedString := strings.Split(query, "AND")
	newString := ""
	if len(splittedString) == 1 {
		newString = "pica." + url.PathEscape(key) + "=" + url.PathEscape(query)
		return newString
	}
	fmt.Println("More than 1 keyword")
	newString = strings.TrimSpace(splittedString[0])
	for i := 1; i < len(splittedString); i++ {
		newString += ("+and+pica." + key + "=" + strings.TrimSpace(url.PathEscape(splittedString[i])))
	}
	return newString
}
/***************************************/
//	UTIL FUNCTIONS AND TYPES
/***************************************/

// FinalClassification is the struct for storing the extracted classifications
type FinalClassification struct {
	Lcc     []string `json:"lcc"`
	Ddc     []string `json:"ddc"`
	Bisacsh []string `json:"bisacsh"`
	BISAC   []string `json:"BISAC"`
	Bkl     []string `json:"bkl"`
	Rvk     []string `json:"rvk"`
}

type OrderedClassification struct {
	Lcc     map[string]int `json:"lcc"`
	Ddc     map[string]int `json:"ddc"`
	Bisacsh map[string]int `json:"bisacsh"`
	BISAC   map[string]int `json:"BISAC"`
	Bkl     map[string]int `json:"bkl"`
	Rvk     map[string]int `json:"rvk"`
}
/*********************************/
//			OUTPUT ANALYZER
/*********************************/


// keyWordAnalyzer prints the n most used classifications for the query
func keyWordAnalyzer(keys map[string]int, n int) {

	type kv struct {
		Key string
		Value int
	}

	sortingList := []kv{}
	for k, v := range keys {
		sortingList = append(sortingList, kv{k, v})
	}
	sort.Slice(sortingList, func(i, j int) bool { return sortingList[i].Value > sortingList[j].Value })

	for i, pairs := range sortingList {
		if i >= n { break
		}
		fmt.Printf("%s : %d\n",pairs.Key, pairs.Value)
	}
}

func classificationAnalyzer(cls OrderedClassification, n int) {
	// Uses reflect to get the struct fields
	w := reflect.ValueOf(cls)
	typ := w.Type()

	for i := 0; i < w.NumField(); i++ {
		fmt.Printf("The %d most common %v classifications for your query\n",n,typ.Field(i).Name)

		var interfaceConverter = w.Field(i).Interface()
		var newMap = interfaceConverter.(map[string]int)
		keyWordAnalyzer(newMap, n)

	}
}

// quickAnalysis prints out the n most used classifications fount in the query results
func quickAnalysis(subjs map[string]int, cls OrderedClassification, n int) {
	fmt.Println("SUBJECT HEADINGS")
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("Printing the %d most common Subject Headings for your query:\n", 5)
	keyWordAnalyzer(subjs, n)
	fmt.Println("\n\nCLASSIFICATIONS")
	fmt.Println("------------------------------------------------------------------\n")
	classificationAnalyzer(cls, n)
}


/*********************************/
//	MAIN FUNCTION
/*********************************/
func main() {

	/*
		+ for reading an existing XML-file : use the function ReadMarshalXML(path)
		+ for downloading data from the web use:
	*/
	


	queryIndex := "all"
	queryText := "Dante"
	//queryText = analyzeInput(queryText, queryIndex)

	path := "outputs/testXML.xml" // hardcoded for now

	XMLFile := getXML(queryText,queryIndex, path, false)

	// Extract list of records
	recs := XMLFile.Records.Record

	fmt.Printf("Number of results: %d \n -----------------------------\n", len(recs))

	finalClass, finalSubjs := ExtractClassifications(recs, false)

	//Save to json
	file, _ := json.MarshalIndent(finalClass, "", "\t")
	_ = ioutil.WriteFile("outputs/testClasses.json", file, 0644)
	file, _ = json.MarshalIndent(finalSubjs, "", " ")
	_ = ioutil.WriteFile("outputs/testSubjects.json", file, 0644)

	quickAnalysis(finalSubjs, finalClass, 7)

}
