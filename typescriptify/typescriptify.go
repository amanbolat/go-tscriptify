package typescriptify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
	"regexp"
	"github.com/guregu/null"
)

type TypeScriptify struct {
	Prefix           string
	Suffix           string
	Indent           string
	CreateFromMethod bool
	DoExportClass    bool
	BackupExtension  string // If empty no backup
	UseInterface     bool

	golangTypes []reflect.Type
	types       map[reflect.Kind]string
	dateTypes   []reflect.Type

	// throwaway, used when converting
	alreadyConverted map[reflect.Type]bool
}

func New() *TypeScriptify {
	result := new(TypeScriptify)
	result.Indent = "\t"
	result.BackupExtension = "backup"

	types := make(map[reflect.Kind]string)

	types[reflect.Bool] = "boolean"

	types[reflect.Int] = "number"
	types[reflect.Int8] = "number"
	types[reflect.Int16] = "number"
	types[reflect.Int32] = "number"
	types[reflect.Int64] = "number"
	types[reflect.Uint] = "number"
	types[reflect.Uint8] = "number"
	types[reflect.Uint16] = "number"
	types[reflect.Uint32] = "number"
	types[reflect.Uint64] = "number"
	types[reflect.Float32] = "number"
	types[reflect.Float64] = "number"

	types[reflect.String] = "string"
	types[reflect.Interface] = "any"

	result.types = types
	result.dateTypes = []reflect.Type{
		reflect.TypeOf(time.Now()),
		reflect.TypeOf(null.NewTime(time.Now(), true)),
	}

	result.Indent = "    "
	result.CreateFromMethod = true
	result.DoExportClass = true

	return result
}

func deepFields(typeOf reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)

	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < typeOf.NumField(); i++ {
		f := typeOf.Field(i)

		kind := f.Type.Kind()
		if f.Anonymous && kind == reflect.Struct {
			//fmt.Println(v.Interface())
			fields = append(fields, deepFields(f.Type)...)
		} else {
			fields = append(fields, f)
		}
	}

	return fields
}

func (t *TypeScriptify) Add(obj interface{}) {
	t.AddType(reflect.TypeOf(obj))
}

func (t *TypeScriptify) AddType(typeOf reflect.Type) {
	t.golangTypes = append(t.golangTypes, typeOf)
}

func (t *TypeScriptify) Convert(customCode map[string]string) (string, error) {
	t.alreadyConverted = make(map[reflect.Type]bool)

	result := ""
	for _, typeof := range t.golangTypes {
		typeScriptCode, err := t.convertType(typeof, customCode)
		if err != nil {
			return "", err
		}
		result += "\n" + strings.Trim(typeScriptCode, " "+t.Indent+"\r\n")
	}
	return result, nil
}

func loadCustomCode(fileName string) (map[string]string, error) {
	result := make(map[string]string)
	f, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return result, err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return result, err
	}

	var currentName string
	var currentValue string
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//[") && strings.HasSuffix(trimmedLine, ":]") {
			currentName = strings.Replace(strings.Replace(trimmedLine, "//[", "", -1), ":]", "", -1)
			currentValue = ""
		} else if trimmedLine == "//[end]" {
			result[currentName] = strings.TrimRight(currentValue, " \t\r\n")
			currentName = ""
			currentValue = ""
		} else if len(currentName) > 0 {
			currentValue += line + "\n"
		}
	}

	return result, nil
}

func (t TypeScriptify) backup(fileName string) error {
	fileIn, err := os.Open(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// No neet to backup, just return:
		return nil
	}
	defer fileIn.Close()

	bytes, err := ioutil.ReadAll(fileIn)
	if err != nil {
		return err
	}

	fileOut, err := os.Create(fmt.Sprintf("%s-%s.%s", fileName, time.Now().Format("2006-01-02T15_04_05.99"), t.BackupExtension))
	if err != nil {
		return err
	}
	defer fileOut.Close()

	_, err = fileOut.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (t TypeScriptify) ConvertToFile(fileName string) error {
	if len(t.BackupExtension) > 0 {
		err := t.backup(fileName)
		if err != nil {
			return err
		}
	}

	customCode, err := loadCustomCode(fileName)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	converted, err := t.Convert(customCode)
	if err != nil {
		return err
	}

	f.WriteString("/* Do not change, this code is generated from Golang structs */\n\n")
	f.WriteString(converted)
	if err != nil {
		return err
	}

	return nil
}

func (t *TypeScriptify) convertType(typeOf reflect.Type, customCode map[string]string) (string, error) {
	for _, v := range t.dateTypes {
		//fmt.Printf("Time is equal: %v\n", reflect.TypeOf(time.Now()) == reflect.TypeOf(null.Time{}))
		//fmt.Printf("Type is %s\n", typeOf)
		if v.String() == typeOf.String() {
			//fmt.Printf("Found Time: %s %s\n", typeOf, v)
			return "", nil
		}
	}

	if _, found := t.alreadyConverted[typeOf]; found { // Already converted
		return "", nil
	}
	t.alreadyConverted[typeOf] = true

	entityName := fmt.Sprintf("%s%s%s", t.Prefix, t.Suffix, typeOf.Name())

	// Set type of typescript kind
	// class, interface, enum
	typeKind := "class"
	if t.UseInterface {
		typeKind = "interface"
	}

	stringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	if typeOf.Kind() == reflect.Int && typeOf.Implements(stringer) {
		typeKind = "enum"
	}

	result := fmt.Sprintf("%s %s {\n", typeKind, entityName)
	if t.DoExportClass {
		result = "export " + result
	}
	builder := typeScriptClassBuilder{
		types:  t.types,
		indent: t.Indent,
	}

	fields := deepFields(typeOf)
	for _, field := range fields {
		jsonTag := field.Tag.Get("json")
		jsonFieldName := ""
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = field.Type.Elem()
		}

		if len(jsonTag) > 0 {
			jsonTagParts := strings.Split(jsonTag, ",")
			if len(jsonTagParts) > 0 {
				jsonFieldName = strings.Trim(jsonTagParts[0], t.Indent)
			}
		}

		if len(jsonFieldName) > 0 && jsonFieldName != "-" {
			var err error
			switch fieldType.Kind() {
			case reflect.Map:
				keyType := "string"
				if k, ok := t.types[fieldType.Key().Kind()]; ok {
					keyType = k
				}

				valType := "any"
				mapValType := fieldType.Elem()

				if mapValType.Kind() == reflect.Ptr {
					mapValType = mapValType.Elem()
				}
				if mapValType.Kind() == reflect.Struct {
					valType = mapValType.Name()

					typeScriptChunk, err := t.convertType(mapValType, customCode)
					if err != nil {
						return "", err
					}
					result = typeScriptChunk + "\n" + result
				}
				if v, ok := t.types[mapValType.Kind()]; ok {
					valType = v
				}

				builder.AddStructField(jsonFieldName, fmt.Sprintf("{[key: %s]: %s}", keyType, valType))
			case reflect.Interface:
				builder.AddStructField(jsonFieldName, "any")
			case reflect.Struct:
				name := fieldType.Name()
				typeScriptChunk, err := t.convertType(fieldType, customCode)
				if err != nil {
					return "", err
				}

				for _, v := range t.dateTypes {
					if v.String() != fieldType.String() {
						continue
					}

					name = "Date"
				}

				result = typeScriptChunk + "\n" + result
				builder.AddStructField(jsonFieldName, name)
			case reflect.Slice:
				elemType := fieldType.Elem()
				if elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}

				switch elemType.Kind() {
				case reflect.Struct:
					typeScriptChunk, err := t.convertType(elemType, customCode)
					if err != nil {
						return "", err
					}
					result = typeScriptChunk + "\n" + result
					builder.AddArrayOfStructsField(jsonFieldName, elemType.Name())
				default:
					err = builder.AddSimpleArrayField(jsonFieldName, elemType.Name(), elemType.Kind())
				}
			case reflect.Int:
				// If it is custom type of int, then it could be enum
				if fieldType.Name() != fieldType.Kind().String() {
					tsChunk, err := t.convertType(fieldType, customCode)
					if err != nil {
						return "", err
					}
					result = tsChunk + "\n" + result
					builder.AddStructField(jsonFieldName, fieldType.Name())
				} else {
					err = builder.AddSimpleField(jsonFieldName, fieldType.Name(), fieldType.Kind())
				}
			default:
				err = builder.AddSimpleField(jsonFieldName, fieldType.Name(), fieldType.Kind())
			}

			if err != nil {
				return "", err
			}
		}
	}

	result += builder.fields
	if t.CreateFromMethod && !t.UseInterface {
		result += fmt.Sprintf("\n%sstatic createFrom(source: any) {\n", t.Indent)
		result += fmt.Sprintf("%s%slet result = new %s();\n", t.Indent, t.Indent, entityName)
		result += builder.createFromMethodBody
		result += fmt.Sprintf("%s%sreturn result;\n", t.Indent, t.Indent)
		result += fmt.Sprintf("%s}\n\n", t.Indent)
	}

	// Set all enum values
	if typeKind == "enum" {
		val := reflect.New(typeOf).Elem()
		for i := 0; i < 10000; i++ {
			val.SetInt(int64(i))
			arr := val.MethodByName("String").Call(nil)
			if len(arr) == 1 {
				enumVal := fmt.Sprintf("%+v", arr[0])

				if strings.Contains(enumVal, fmt.Sprintf("(%d)", i)) {
					continue
				}
				result += fmt.Sprintf("%s%s = '%s',\n", t.Indent, ToCamel(enumVal), enumVal)
			}
		}
	}

	if customCode != nil {
		code := customCode[entityName]
		result += t.Indent + "//[" + entityName + ":]\n" + code + "\n\n" + t.Indent + "//[end]\n"
	}

	// Type is enum
	// Loop through all possible values
	//if typeKind == "enum" {
	//
	//}

	result += "}"

	return result, nil
}

type typeScriptClassBuilder struct {
	types                map[reflect.Kind]string
	indent               string
	fields               string
	createFromMethodBody string
}

func (t *typeScriptClassBuilder) AddSimpleArrayField(fieldName, fieldType string, kind reflect.Kind) error {
	if typeScriptType, ok := t.types[kind]; ok {
		if len(fieldName) > 0 {
			t.fields += fmt.Sprintf("%s%s: %s[];\n", t.indent, fieldName, typeScriptType)
			t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"];\n", t.indent, t.indent, fieldName, fieldName)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Cannot find type for: %s (%s/%s)", kind.String(), fieldName, fieldType))
}

func (t *typeScriptClassBuilder) AddSimpleField(fieldName, fieldType string, kind reflect.Kind) error {
	if typeScriptType, ok := t.types[kind]; ok {
		if len(fieldName) > 0 {
			t.fields += fmt.Sprintf("%s%s: %s;\n", t.indent, fieldName, typeScriptType)
			t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"];\n", t.indent, t.indent, fieldName, fieldName)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Cannot find type '%s' for field '%s' ", fieldType, fieldName))
}

func (t *typeScriptClassBuilder) AddStructField(fieldName, fieldType string) {
	t.fields += fmt.Sprintf("%s%s: %s;\n", t.indent, fieldName, fieldType)
	t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"] ? %s.createFrom(source[\"%s\"]) : null;\n", t.indent, t.indent, fieldName, fieldName, fieldType, fieldName)
}

func (t *typeScriptClassBuilder) AddArrayOfStructsField(fieldName, fieldType string) {
	t.fields += fmt.Sprintf("%s%s: %s[];\n", t.indent, fieldName, fieldType)
	t.createFromMethodBody += fmt.Sprintf("%s%sresult.%s = source[\"%s\"] ? source[\"%s\"].map(function(element) { return %s.createFrom(element); }) : null;\n", t.indent, t.indent, fieldName, fieldName, fieldName, fieldType)
}


// Helpers for camel case
var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}
