package framework

import (
	"os"
	"text/template"
)

func ParseTemplate(tpl string, data interface{}) error {
	// Create a new template object
	t, err := template.New("test").Parse(tpl)
	if err != nil {
		return err
	}

	// Inject the data into the template and print the result
	err = t.Execute(os.Stdout, data)
	if err != nil {
		return err
	}

	return nil
}
