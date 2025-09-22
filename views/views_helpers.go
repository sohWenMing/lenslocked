package views

import (
	"fmt"
	"html/template"
)

/*
used to load up and parse all templates at the beginning of execution of the program
TemplateMust function will panic any error found during parsing, which will shut down execution of the program
*/

func getTemplatePaths(tplStrings []string, baseFolderName string) []string {
	fullPaths := make([]string, len(tplStrings))
	for i, tplString := range tplStrings {
		fullPath := fmt.Sprintf("%s/%s", baseFolderName, tplString)
		fullPaths[i] = fullPath
	}
	return fullPaths
}

// helper function - used to create the final string slice of template paths that will be parsed

func TemplateMust(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return t
}
