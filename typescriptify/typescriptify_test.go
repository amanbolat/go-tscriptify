package typescriptify

import (
	"bitbucket.org/amanbolat/caconsole/shipment/model"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Address struct {
	// Used in html
	Duration float64 `json:"duration"`
	Text1    string  `json:"text,omitempty"`
	// Ignored:
	Text2 string `json:",omitempty"`
	Text3 string `json:"-"`
}

type Dummy struct {
	Something     string      `json:"something"`
	SomeInterface interface{} `json:"some_interface"`
}

type HasName struct {
	Name string `json:"name"`
}

type Person struct {
	HasName
	Nicknames  []string             `json:"nicknames"`
	Addresses  []Address            `json:"addresses"`
	Dummy      Dummy                `json:"a"`
	Ptr        *Dummy               `json:"b"`
	SlicePtr   []*Dummy             `json:"slice_ptr"`
	Map        map[string]*Dummy    `json:"map"`
	Birthday   time.Time            `json:"birthday"`
}

type Shipment struct {
	Statuses map[string]time.Time `json:"statuses"`
}

func TestTypeMapStringTime(t *testing.T)  {
	converter := New()
	converter.AddType(reflect.TypeOf(Shipment{}))
	converter.CreateFromMethod = false

	desiredResult := `export class Shipment {
		statuses: {[key: string]: Date};
}`

	testConverter(t, converter, desiredResult)
}

func TestTypescriptifyWithTypes(t *testing.T) {
	converter := New()

	converter.AddType(reflect.TypeOf(Person{}))
	converter.CreateFromMethod = false
	converter.DoExportClass = true

	desiredResult := `export class Dummy {
        something: string;
		some_interface: any;
}
export class Address {
        duration: number;
        text: string;
}
export class Person {
        name: string;
        nicknames: string[];
        addresses: Address[];
        a: Dummy;
		b: Dummy;
		slice_ptr: Dummy[];
		map: {[key: string]: Dummy};
		birthday: Date;
}`
	testConverter(t, converter, desiredResult)
}

func TestTypescriptifyWithInstances(t *testing.T) {
	converter := New()

	converter.Add(Person{})
	converter.Add(Dummy{})
	converter.CreateFromMethod = false

	desiredResult := `export class Dummy {
        something: string;
		some_interface: any;
}
export class Address {
        duration: number;
        text: string;
}
export class Person {
        name: string;
        nicknames: string[];
        addresses: Address[];
        a: Dummy;
		b: Dummy;
		slice_ptr: Dummy[];
		map: {[key: string]: Dummy};
		birthday: Date;
}`
	testConverter(t, converter, desiredResult)
}

func TestTypescriptifyWithDoubleClasses(t *testing.T) {
	converter := New()

	converter.AddType(reflect.TypeOf(Person{}))
	converter.AddType(reflect.TypeOf(Person{}))
	converter.CreateFromMethod = false
	converter.DoExportClass = true

	desiredResult := `export class Dummy {
        something: string;
		some_interface: any;
}
export class Address {
        duration: number;
        text: string;
}
export class Person {
        name: string;
        nicknames: string[];
        addresses: Address[];
        a: Dummy;
		b: Dummy;
		slice_ptr: Dummy[];
		map: {[key: string]: Dummy};
		birthday: Date;
}`
	testConverter(t, converter, desiredResult)
}

func TestWithPrefixes(t *testing.T) {
	converter := New()

	converter.Prefix = "test_"

	converter.Add(Address{})
	converter.Add(Dummy{})
	converter.CreateFromMethod = false

	desiredResult := `export class test_Address {
        duration: number;
        text: string;
}
export class test_Dummy {
        something: string;
		some_interface: any;
}`
	testConverter(t, converter, desiredResult)
}

func testConverter(t *testing.T, converter *TypeScriptify, desiredResult string) {
	typeScriptCode, err := converter.Convert(nil)
	if err != nil {
		panic(err.Error())
	}

	typeScriptCode = strings.Trim(typeScriptCode, " \t\n\r")
	if typeScriptCode != desiredResult {
		lines1 := strings.Split(typeScriptCode, "\n")
		lines2 := strings.Split(desiredResult, "\n")

		if len(lines1) != len(lines2) {
			os.Stderr.WriteString(fmt.Sprintf("Lines: %d != %d\n", len(lines1), len(lines2)))
			os.Stderr.WriteString(fmt.Sprintf("Expected:\n%s\n\nGot:\n%s\n", desiredResult, typeScriptCode))
			t.Fail()
		} else {
			for i := 0; i < len(lines1); i++ {
				line1 := strings.Trim(lines1[i], " \t\r\n")
				line2 := strings.Trim(lines2[i], " \t\r\n")
				if line1 != line2 {
					os.Stderr.WriteString(fmt.Sprintf("%d. line don't match: `%s` != `%s`\n", i+1, line1, line2))
					os.Stderr.WriteString(fmt.Sprintf("Expected:\n%s\n\nGot:\n%s\n", desiredResult, typeScriptCode))
					t.Fail()
				}
			}
		}
	}
}

func TestEnum(t *testing.T) {
	typeOf := reflect.TypeOf((*model.PaymentMethod)(nil))
	t.Logf("%+v", typeOf.Elem().Name())
}
