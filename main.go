/* Read and parse xml (format mods) from sru.k10plus/gvk */

package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
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

/*
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
*/

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

func getXML(queryText, queryIndex, path, maxResult string, save bool) []SingleResult {
	body := sendRequest(queryText, queryIndex,maxResult)
	if save { //TODO check if the directory exist, if not mkdir
		_ = ioutil.WriteFile(path, body, 0666)
	}

	parsedXML := Root{} // TODO: Avoid new keyword
	xml.Unmarshal(body, &parsedXML)

	XmlFile := parsedXML.Records.Record
	return XmlFile
}

// sendRequest sends the query to the sru address
func sendRequest(queryText, queryIndex, maxResult string) (body []byte) {
	request := buildQuery(queryText, queryIndex,maxResult)

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

func buildQuery(queryText, queryIndex, maxResult string) (request string) {

	switch queryIndex {
	case "slw":
		fallthrough
	case "tit":
		fallthrough
	case "bkl":
		fallthrough
	case "per":
		fullQuery := analyzeInput(queryText, queryIndex)
		request = fmt.Sprintf("https://sru.gbv.de/gvk?version=1.1&operation=searchRetrieve&query=%s&recordSchema=mods&maximumRecords=%s",fullQuery,maxResult)
	default: // Use the full text search
		GeneralQuery := analyzeInput(queryText, "all")
		request = fmt.Sprintf("https://sru.gbv.de/gvk?version=1.1&operation=searchRetrieve&query=%s&recordSchema=mods&maximumRecords=%s",GeneralQuery,maxResult)
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
		newString += "+and+pica." + key + "=" + strings.TrimSpace(url.PathEscape(splittedString[i]))
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

	var sortingList []kv
	for k, v := range keys {
		sortingList = append(sortingList, kv{k, v})
	}
	sort.Slice(sortingList, func(i, j int) bool { return sortingList[i].Value > sortingList[j].Value })
		for i, pairs := range sortingList {
			if i >= n {
				break
			}
			fmt.Printf("%s : %d\n", pairs.Key, pairs.Value)
		}
	}

func classificationAnalyzer(cls OrderedClassification, n int) {
	// Uses reflect to get the struct fields
	w := reflect.ValueOf(cls)
	typ := w.Type()

	for i := 0; i < w.NumField(); i++ {
			var interfaceConverter = w.Field(i).Interface()
			var newMap = interfaceConverter.(map[string]int)
		if len(newMap) != 0 {
			fmt.Printf("The %d most common %v classifications for your query\n", n, typ.Field(i).Name)
			keyWordAnalyzer(newMap, n)
		}
	}
}

// quickAnalysis prints out the n most used classifications fount in the query results
func quickAnalysis(subjs map[string]int, cls OrderedClassification, n int) {

	if len(subjs) != 0 {
		fmt.Println("SUBJECT HEADINGS")
		fmt.Println("------------------------------------------------------------------")
		fmt.Printf("Printing the %d most common Subject Headings for your query:\n", 5)
		keyWordAnalyzer(subjs, n)
	}
	fmt.Println("\n\nCLASSIFICATIONS")
	fmt.Println("------------------------------------------------------------------")
	classificationAnalyzer(cls, n)
}


func saveJson(i interface{}, fp string) {
	file, _ := json.MarshalIndent(i, "", " ")
	_ = ioutil.WriteFile(fp, file, 0644)
}

/*********************************/
//	MAIN FUNCTION
/*********************************/
func main() {

	/*
		+ for reading an existing XML-file : use the function ReadMarshalXML(path)
		+ for downloading data from the web use:
	*/
	var verbose bool
	var save bool

	buf := bufio.NewReader(os.Stdin)
	queryIndex := flag.String("k","","Pica query field")
	queryText := flag.String("q", "","The query string")
	n := flag.Int( "n",7,"How many results should be printed")
	flag.BoolVar(&verbose, "v", false, "Print all outputs (not recommended)")
	flag.BoolVar(&save, "s", false,"Save the retrieved xml and the json outputs")
	path := flag.String("p","outputs/","Path for the results")
	maxResult := flag.String("m","500","How many results should be retrieved")
	flag.Parse()


	if flag.NFlag() < 2 {
		log.Println("Not enough commands, please use at least -q")
		fmt.Println("Usage:\n\t$", os.Args[0]," + parameters (at least -q): \n\n" +
			"\t\t-k\t\t query Key (bkl, tit, per, slw (default: all)\n" +
			"\t\t-q\t\t query String (the text to search)\n" +
			"\t\t-v\t\t verbose (bool) prints everything on the console, not recommended (debug purpose only)\n" +
			"\t\t-s\t\t save (bool) if true saves the retrieved xml and the json output files (default : false)\n" +
			"\t\t-n\t\t number of Results (integer) specifiers how many results will be printed on the console\n" +
			"\t\t-p\t\t path (string) which subdirectory of the current working directory will be used to save the files (works only if -s is set)" +
			"\t\t-m\t\t maxResult how many entries should be retrieved (default: 120)")
		_ , err := buf.ReadByte()
		log.Println(err)
		fmt.Println("Exiting program")
		return
	}
	fmt.Println("                                                                  \n                                  88      a8P  88    ,a8888a,     \n                                  88    ,88' ,d88  ,8P\"'  `\"Y8,   \n                                  88  ,88\" 888888 ,8P        Y8,  \n                                  88,d88'      88 88          88  \n                                  8888\"88,     88 88          88  \n                                  88P   Y8b    88 `8b        d8'  \n                                  88     \"88,  88  `8ba,  ,ad8'   \n                                  88       Y8b 88    \"Y8888P\"     \n                                                                  \n                                                                  \n                                                                                            \n         ,ad8888ba,  88                                88    ad88 88                        \n        d8\"'    `\"8b 88                                \"\"   d8\"   \"\"                        \n       d8'           88                                     88                              \n       88            88 ,adPPYYba, ,adPPYba, ,adPPYba, 88 MM88MMM 88  ,adPPYba, 8b,dPPYba,  \n       88            88 \"\"     `Y8 I8[    \"\" I8[    \"\" 88   88    88 a8P_____88 88P'   \"Y8  \n       Y8,           88 ,adPPPPP88  `\"Y8ba,   `\"Y8ba,  88   88    88 8PP\"\"\"\"\"\"\" 88          \n        Y8a.    .a8P 88 88,    ,88 aa    ]8I aa    ]8I 88   88    88 \"8b,   ,aa 88          \n         `\"Y8888Y\"'  88 `\"8bbdP\"Y8 `\"YbbdP\"' `\"YbbdP\"' 88   88    88  `\"Ybbd8\"' 88          \n                                                                                            \n                                                                                            ")

	XMLFile := getXML(*queryText,*queryIndex, *path + "Results.xml", *maxResult, save)
	fmt.Printf("Number of results: %d \n -----------------------------\n", len(XMLFile))

	finalClass, finalSubjs := ExtractClassifications(XMLFile, verbose)

	//Save to json TODO: Build a function to save json
	if save {
		saveJson(finalClass, *path + "Classes.json")
		saveJson(finalSubjs, *path + "SubjectHeadings.json")
	}

	// Print analysis
	quickAnalysis(finalSubjs, finalClass, *n)

}

