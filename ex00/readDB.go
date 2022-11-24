package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Item struct {
	Itemname  string `xml:"itemname" json:"ingredient_name"`
	Itemcount string `xml:"itemcount" json:"ingredient_count"`
	Itemunit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type Cake struct {
	Name       string `xml:"name" json:"name"`
	Stovetime  string `xml:"stovetime" json:"time"`
	Ingredient []Item `xml:"ingredients>item" json:"ingredients"`
}

type Recipe struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Recipes []Cake   `xml:"cake" json:"cake"`
}

type XML Recipe
type JSON Recipe

type DBReader interface {
	Read(content []byte) Recipe
	Rewrite(recipes Recipe)
}

func (f *XML) Read(content []byte) Recipe {
	err := xml.Unmarshal(content, f)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return Recipe(*f)
}

func (f *JSON) Read(content []byte) Recipe {
	err := json.Unmarshal(content, f)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return Recipe(*f)
}

func (f *XML) Rewrite(recipes Recipe) {
	file, err := os.Create("file.json")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	out, err := json.MarshalIndent(recipes, "", "	")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprint(file, string(out))
	file.Close()
}

func (f *JSON) Rewrite(recipes Recipe) {
	file, err := os.Create("file.xml")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	out, err := xml.MarshalIndent(recipes, "", "	")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprint(file, string(out))
	file.Close()
}

func parseFile(reader DBReader, filename string) {
	var recipes Recipe
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	filecontent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	recipes = reader.Read(filecontent)
	reader.Rewrite(recipes)
}

func checkFormat(filename string) string {
	if strings.HasSuffix(filename, ".xml") {
		return "xml"
	} else if strings.HasSuffix(filename, ".json") {
		return "json"
	} else {
		return ""
	}
}

func start(filename string) {
	format := checkFormat(filename)
	switch format {
	case "xml":
		myStruct := new(XML)
		parseFile(myStruct, filename)
	case "json":
		myStruct := new(JSON)
		parseFile(myStruct, filename)
	default:
		fmt.Fprint(os.Stderr, "error: invalid file extension\n")
		os.Exit(1)
	}
}

func main() {
	useFile := flag.String("f", "", "parse file")
	flag.Parse()

	if *useFile != "" {
		start(*useFile)
	} else {
		fmt.Println("Use '-f' flag to pass argument")
	}
}
