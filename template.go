package main

import (
	flag "github.com/spf13/pflag"
	"fmt"
	"strings"
	"html/template"
	"io"
	"os"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	inputFile = flag.StringP("input", "i", "", "the input template file to parse")
	outputFile = flag.StringP("output", "o", "", "the desired output file location")
	variables = flag.StringArrayP("var", "v", nil, "variables used to fill in the template")
	jsonPath = flag.StringP("json", "j", "", "pass the template variables in JSON format or reference a JSON file")
	yamlPath = flag.StringP("yaml", "y", "", "pass the template variables in YAML format or reference a YAML file")
)

var (
	inputTemplate *template.Template
	outputWriter io.Writer
	variableMap map[string]string
)

func main() {
	flag.Parse()
	// create a map with the variables in the cmd line flags
	fillVariableMap()
	// parse the input template
	parseTemplate()
	// create the output writer
	specifyOutput()
	// execute the template into the output writer
	if err := inputTemplate.Execute(outputWriter, variableMap); err != nil {
		fmt.Printf("Error applying template: %v\n", err)
		os.Exit(1)
	}
}

// Creates a template.Template instance depending on the --input flag.
// Panics when no input file is provided or when failing to parse the template.
func parseTemplate() {
	if *inputFile == "" {
		// no input file given
		fmt.Println("No input file provided!")
		os.Exit(1)
	}
	tmpl, err := template.ParseFiles(*inputFile)
	if err != nil {
		// parsing error
		fmt.Printf("Error parsing the input file: %v\n", err)
		os.Exit(1)
	}
	inputTemplate = tmpl
}

// Parses all the provided variable flags and puts them into a map structure.
func fillVariableMap() {
	variableMap = make(map[string]string)
	// add regular variables, if any
	for _, v := range *variables {
		s := strings.SplitN(v, "=", 2)
		variableMap[s[0]] = s[1]
	}
	// add JSON variables, if provided
	if *jsonPath != "" {
		// JSON parameter provided
		if jsonFile, err := os.Open(*jsonPath); err != nil {
			// parse JSON file
			if err := json.NewDecoder(jsonFile).Decode(&variableMap); err != nil {
				fmt.Printf("Error parsing JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			// something went wrong opening the file
			fmt.Printf("Could not open the JSON file at %v: %v\n", jsonPath, err)
			os.Exit(1)
		}
	}
	// add YAML variables, if provided
	if *yamlPath != "" {
		// YAML parameter provided
		if content, err := ioutil.ReadFile(*yamlPath); err == nil {
			// read file
			// parse YAML file
			if err := yaml.Unmarshal(content, &variableMap); err != nil {
				fmt.Printf("Error parsing YAML: %v\n", err)
				os.Exit(1)
			}
		} else {
			// something went wrong opening the file
			fmt.Printf("Could not open the YAML file at %v: %v\n", yamlPath, err)
			os.Exit(1)
		}
	}
}

// Creates an output io.Writer depending on the --output flag.
// If the flag is not provided, returns os.Stdout otherwise it tries to open the provided file.
func specifyOutput() {
	if *outputFile == "" {
		// no output file, write to Stdout
		outputWriter = os.Stdout
	} else {
		w, err := os.Create(*outputFile)
		if err != nil {
			// error creating output file
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		outputWriter = w
	}
}
