package skills

import (
	"io/ioutil"
	"net/http"
)

type Skill struct {
	Name   string
	HowTo  string
	Return func(query string) string
}

func (s *Skill) GetJSON(endpoint string) (string, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *Skill) Describe(query string, result string) string {
	return `
		Skill Name: ` + s.Name + `
		Skill Context: ` + s.HowTo + `

		Skill Query (generated from Question): ` + query + `
		Skill Summary: ` + result + `
		------------------------------------------------------------
	`
}
