package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"regexp"
)

const (
	registryTemplate = `// This file is auto generated, DO NOT EDIT.
package libdnsregistry

import (
	{{- range .}}
	{{.ModuleAlias}} "{{.ModulePath}}"
	{{- end}}
)

var registry = RegistryStore{
	{{- range .}}
	"{{.Name}}": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[{{.ModuleAlias}}.Provider](conf)
		},
	},
	{{- end}}
}
`
)

type providerSource struct {
	Repo              string `json:"repo"`
	ModulePath        string `json:"module_path"`
	LibdnsRequirement string `json:"libdns_requirement"`
	Status            string `json:"status"`
	Reason            string `json:"reason"`
}

type providerTemplate struct {
	Name        string
	ModuleAlias string
	ModulePath  string
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		panic("Usage: go run generate code <json path> <output path>")
	}

	providers := new([]providerSource)

	listJSON, err := os.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(listJSON, providers)
	if err != nil {
		panic(err)
	}

	filteredProviders := []providerTemplate{}
	re := regexp.MustCompile(`[^a-z0-9]`)
	for _, provider := range *providers {
		if provider.Status == "compatible" {
			filteredProviders = append(filteredProviders, providerTemplate{
				Name:        provider.Repo,
				ModulePath:  provider.ModulePath,
				ModuleAlias: fmt.Sprintf("libdns%s", re.ReplaceAllString(provider.Repo, `_`)),
			})
		}
	}

	switch args[0] {
	case "code":
		generateCode(args[2], filteredProviders)

	default:
		panic("type not supported")
	}
}

func generateCode(path string, providers []providerTemplate) {
	registryTpl, err := template.New("libdns").Parse(registryTemplate)
	if err != nil {
		panic(err)
	}

	content := &bytes.Buffer{}
	if err := registryTpl.Execute(content, providers); err != nil {
		panic(err)
	}

	// formattedContent, err := format.Source(content.Bytes())
	// if err != nil {
	// 	panic(err)
	// }

	//nolint:gosec,mnd
	if err := os.WriteFile(path, content.Bytes(), 0o644); err != nil {
		panic(err)
	}
}
