package main

import (
	"github.com/amanbolat/go-tscriptify/typescriptify"
	"bitbucket.org/amanbolat/caconsole/shipment/model"
)

type Address struct {
	// Used in html
	City    string  `json:"city"`
	Number  float64 `json:"number"`
	Country string  `json:"country,omitempty"`
}

type PersonalInfo struct {
	Hobbies []string `json:"hobby"`
	PetName string   `json:"pet_name"`
}

type Person struct {
	Name         string       `json:"name"`
	PersonalInfo PersonalInfo `json:"personal_info"`
	Nicknames    []string     `json:"nicknames"`
	Addresses    []Address    `json:"addresses"`
	Shipments []*model.Shipment `json:"shipments"`
}

func main() {
	converter := typescriptify.New()
	converter.CreateFromMethod = true
	converter.Indent = "    "
	converter.UseInterface = true

	converter.Add(Person{})

	err := converter.ConvertToFile("browser_test/example_output.ts")
	if err != nil {
		panic(err.Error())
	}
}
