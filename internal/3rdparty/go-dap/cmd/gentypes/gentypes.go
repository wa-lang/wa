// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// gentypes generates Go types from debugProtocol.json
//
// Usage:
//
// $ gentypes <path to debugProtocol.json>
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

var (
	uFlag = flag.Bool("u", false, "updates the debugProtocol.json file before generating the code")
	oFlag = flag.String("o", "", "specifies the output file name. If unspecified, outputs to stdout")
)

// parseRef parses the value of a "$ref" key.
// For example "#definitions/ProtocolMessage" => "ProtocolMessage".
func parseRef(refValue interface{}) string {
	refContents := refValue.(string)
	if !strings.HasPrefix(refContents, "#/definitions/") {
		log.Fatal("want ref to start with '#/definitions/', got ", refValue)
	}

	return replaceGoTypename(refContents[14:])
}

// goFieldName converts a property name from its JSON representation to an
// exported Go field name.
// For example "__some_property_name" => "SomePropertyName".
func goFieldName(jsonPropName string) string {
	clean := strings.ReplaceAll(jsonPropName, "_", " ")
	titled := strings.Title(clean)
	return strings.ReplaceAll(titled, " ", "")
}

// parsePropertyType takes the JSON value of a property field and extracts
// the Go type of the property. For example, given this map:
//
//	{
//	  "type": "string",
//	  "description": "The command to execute."
//	},
//
// It will emit "string".
func parsePropertyType(propValue map[string]interface{}) string {
	if ref, ok := propValue["$ref"]; ok {
		return parseRef(ref)
	}

	if _, ok := propValue["oneOf"]; ok {
		return "interface{}"
	}
	propType, ok := propValue["type"]
	if !ok {
		log.Fatal("property with no type or ref:", propValue)
	}

	switch propType.(type) {
	case string:
		switch propType {
		case "string":
			return "string"
		case "number":
			return "int"
		case "integer":
			return "int"
		case "boolean":
			return "bool"
		case "array":
			propItems, ok := propValue["items"]
			if !ok {
				log.Fatal("missing items type for property of array type:", propValue)
			}
			propItemsMap := propItems.(map[string]interface{})
			return "[]" + parsePropertyType(propItemsMap)
		case "object":
			// When the type of a property is "object", we'll emit a map with a string
			// key and a value type that depends on the type of the
			// additionalProperties field.
			additionalProps, ok := propValue["additionalProperties"]
			if !ok {
				log.Fatal("missing additionalProperties field when type=object:", propValue)
			}
			var valueType string
			switch actual := additionalProps.(type) {
			case bool:
				valueType = "interface{}"
			case map[string]interface{}:
				valueType = parsePropertyType(actual)
			default:
				log.Fatal("unexpected additionalProperties value:", additionalProps)
			}
			return fmt.Sprintf("map[string]%v", valueType)
		case "interface{}":
			return "interface{}"
		default:
			log.Fatalf("unknown property type value %v in %v", propType, propValue)
		}

	case []interface{}:
		return "interface{}"

	default:
		log.Fatal("unknown property type", propType)
	}

	panic("unreachable")
}

// maybeParseInheritance helps parse types that inherit from other types.
// A type description can have an "allOf" key, which means it inherits from
// another type description. Returns the name of the base type specified in
// allOf, and the description of the inheriting type.
//
// Example:
//
//	"allOf": [ { "$ref": "#/definitions/ProtocolMessage" },
//	           {... type description ...} ]
//
// Returns base type ProtocolMessage and a map representing type description.
// If there is no "allOf", returns an empty baseTypeName and descMap itself.
func maybeParseInheritance(descMap map[string]json.RawMessage) (baseTypeName string, typeDescJson map[string]json.RawMessage) {
	allOfListJson, ok := descMap["allOf"]
	if !ok {
		return "", descMap
	}

	var allOfSliceOfJson []json.RawMessage
	if err := json.Unmarshal(allOfListJson, &allOfSliceOfJson); err != nil {
		log.Fatal(err)
	}
	if len(allOfSliceOfJson) != 2 {
		log.Fatal("want 2 elements in allOf list, got", allOfSliceOfJson)
	}

	var baseTypeRef map[string]interface{}
	if err := json.Unmarshal(allOfSliceOfJson[0], &baseTypeRef); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(allOfSliceOfJson[1], &typeDescJson); err != nil {
		log.Fatal(err)
	}
	return parseRef(baseTypeRef["$ref"]), typeDescJson
}

// emitToplevelType emits a single type into a string. It takes the type name
// and a serialized json object representing the type. The json representation
// will have fields: "type", "properties" etc.
func emitToplevelType(typeName string, descJson json.RawMessage, goTypeIsStruct map[string]bool) string {
	var b strings.Builder
	var baseType string

	// We don't parse the description all the way to map[string]interface{}
	// because we have to retain the original JSON-order of properties (in this
	// type as well as any nested types like "body").
	var descMap map[string]json.RawMessage
	if err := json.Unmarshal(descJson, &descMap); err != nil {
		log.Fatal(err)
	}
	baseType, descMap = maybeParseInheritance(descMap)

	typeJson, ok := descMap["type"]
	if !ok {
		log.Fatal("want description to have 'type', got ", descMap)
	}

	var descTypeString string
	if err := json.Unmarshal(typeJson, &descTypeString); err != nil {
		log.Fatal(err)
	}

	var comment string
	descriptionJson, ok := descMap["description"]
	if ok {
		if err := json.Unmarshal(descriptionJson, &comment); err != nil {
			log.Fatal(err)
		}
	}

	if len(comment) > 0 {
		comment = commentOutEachLine(fmt.Sprintf("%s: %s", typeName, comment))
		fmt.Fprint(&b, comment)
	}

	if descTypeString == "string" {
		fmt.Fprintf(&b, "type %s string\n", typeName)
		return b.String()
	} else if descTypeString == "object" {
		fmt.Fprintf(&b, "type %s struct {\n", typeName)
		if len(baseType) > 0 {
			fmt.Fprintf(&b, "\t%s\n\n", baseType)
		}
	} else {
		log.Fatal("want description type to be object or string, got ", descTypeString)
	}

	var propsMapOfJson map[string]json.RawMessage
	if propsJson, ok := descMap["properties"]; ok {
		if err := json.Unmarshal(propsJson, &propsMapOfJson); err != nil {
			log.Fatal(err)
		}
	} else {
		b.WriteString("}\n")
		return b.String()
	}

	propsNamesInOrder, err := keysInOrder(descMap["properties"])
	if err != nil {
		log.Fatal(err)
	}

	// Stores the properties that are required.
	requiredMap := make(map[string]bool)

	if requiredJson, ok := descMap["required"]; ok {
		var required []interface{}
		if err := json.Unmarshal(requiredJson, &required); err != nil {
			log.Fatal(err)
		}
		for _, r := range required {
			requiredMap[r.(string)] = true
		}
	}

	// Some types will have a "body" which should be emitted as a separate type.
	// Since we can't emit a whole new Go type while in the middle of emitting
	// another type, we save it for later and emit it after the current type is
	// done.
	bodyType := ""

	for _, propName := range propsNamesInOrder {
		// The JSON schema is designed for the TypeScript type system, where a
		// subclass can redefine a field in a superclass with a refined type (such
		// as specific values for a field). To ensure we emit Go structs that can
		// be unmarshaled from JSON messages properly, we must limit each field
		// to appear only once in hierarchical types.
		if propName == "type" && (typeName == "Request" || typeName == "Response" || typeName == "Event") {
			continue
		}
		if propName == "command" && typeName != "Request" && typeName != "Response" {
			continue
		}
		if propName == "event" && typeName != "Event" {
			continue
		}
		if propName == "arguments" && typeName == "Request" {
			continue
		}

		var propDesc map[string]interface{}
		if err := json.Unmarshal(propsMapOfJson[propName], &propDesc); err != nil {
			log.Fatal(err)
		}

		if propName == "body" {
			if typeName == "Response" || typeName == "Event" {
				continue
			}

			var bodyTypeName string
			if ref, ok := propDesc["$ref"]; ok {
				bodyTypeName = parseRef(ref)
			} else {
				bodyTypeName = typeName + "Body"
				bodyType = emitToplevelType(bodyTypeName, propsMapOfJson["body"], goTypeIsStruct)
			}

			if requiredMap["body"] {
				fmt.Fprintf(&b, "\t%s %s `json:\"body\"`\n", "Body", bodyTypeName)
			} else {
				fmt.Fprintf(&b, "\t%s %s `json:\"body,omitempty\"`\n", "Body", bodyTypeName)
			}
		} else if propName == "arguments" && (typeName == "LaunchRequest" || typeName == "AttachRequest") {
			// Special case for LaunchRequest or AttachRequest arguments, which are implementation
			// defined and don't have pre-set field names in the specification.
			fmt.Fprintln(&b, "\tArguments json.RawMessage `json:\"arguments\"`")
		} else {
			// Go type of this property.
			goType := parsePropertyType(propDesc)

			jsonTag := fmt.Sprintf("`json:\"%s", propName)
			if requiredMap[propName] {
				jsonTag += "\"`"
			} else if typeName == "ContinueResponseBody" && propName == "allThreadsContinued" {
				// This one special field must not have the omitempty tag, despite being
				// optional. If this attribute is missing the client will (according to
				// the specification) assume a value of 'true' for backward
				// compatibility. See: https://github.com/google/go-dap/issues/39
				jsonTag += "\"`"
			} else if typeName == "InitializeRequestArguments" && (propName == "linesStartAt1" || propName == "columnsStartAt1") {
				// These two special fields must not have the omitempty tag, despite being
				// optional. If this attribute is missing the server will (according to
				// the specification) assume a value of 'true'.
				jsonTag += "\"`"
			} else if typeName == "ErrorMessage" && propName == "showUser" {
				// For launch/attach errors, vscode will treat omitted values the same way as true,
				// so to suppress visible reporting, we must report false explicitly.
				jsonTag += "\"`"
			} else {
				jsonTag += ",omitempty\"`"
				// If the field should be omitted when empty and is a struct type in Go, make it a pointer,
				// because non-pointer structs get initialized with default values in Go (and not nil), and
				// are then indistinguishable from structs with values actually set to zero when serializing
				// to JSON. Making them a pointer makes them initialize to nil, which is then indeed omitted
				// during serialization.
				if _, ok := propDesc["$ref"]; ok {
					// If we have a ref, then goType is the parsed ref
					if goTypeIsStruct[goType] {
						goType = "*" + goType
					}
				}
			}
			fmt.Fprintf(&b, "\t%s %s %s\n", goFieldName(propName), goType, jsonTag)

		}
	}

	b.WriteString("}\n")

	if len(bodyType) > 0 {
		b.WriteString("\n")
		b.WriteString(bodyType)
	}

	return b.String()
}

// keysInOrder returns the keys in json object in b, in their original order.
// Based on https://github.com/golang/go/issues/27179#issuecomment-415559968
func keysInOrder(b []byte) ([]string, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, errors.New("expected start of object")
	}
	var keys []string
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		if t == json.Delim('}') {
			return keys, nil
		}
		keys = append(keys, t.(string))
		if err := skipValue(d); err != nil {
			return nil, err
		}
	}
}

// replaceGoTypename replaces conflicting type names in the JSON schema with
// proper Go type names.
func replaceGoTypename(typeName string) string {
	// Since we have a top-level interface named Message, we replace the DAP
	// message type Message with ErrorMessage.
	if typeName == "Message" {
		return "ErrorMessage"
	}
	return typeName
}

var errEnd = errors.New("invalid end of array or object")

func skipValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t {
	case json.Delim('['), json.Delim('{'):
		for {
			if err := skipValue(d); err != nil {
				if err == errEnd {
					break
				}
				return err
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return errEnd
	}
	return nil
}

// commentOutEachLine returns s such that a Go comment marker ("//") is
// prepended to each line.
func commentOutEachLine(s string) string {
	parts := strings.Split(s, "\n")
	var sb strings.Builder

	for _, p := range parts {
		fmt.Fprintf(&sb, "// %s\n", p)
	}
	return sb.String()
}

// emitMethodsForType may emit methods for typeName into sb.
func emitMethodsForType(sb *strings.Builder, typeName string) {
	if typeName == "ProtocolMessage" {
		fmt.Fprintln(sb, "func (m *ProtocolMessage) GetSeq() int {return m.Seq}")
	}
	if typeName == "Request" {
		fmt.Fprintln(sb, "func (r *Request) GetRequest() *Request {return r}")
	}
	if typeName == "Response" {
		fmt.Fprintln(sb, "func (r *Response) GetResponse() *Response {return r}")
	}
	if typeName == "Event" {
		fmt.Fprintln(sb, "func (e *Event) GetEvent() *Event {return e}")
	}
	if typeName == "LaunchRequest" || typeName == "AttachRequest" {
		fmt.Fprintf(sb, "func (r *%s) GetArguments() json.RawMessage { return r.Arguments }\n", typeName)
	}
}

func emitCtor(sb *strings.Builder, reqs, resps, events []string) {
	fmt.Fprint(sb, `
// Mapping of request commands and corresponding struct constructors that
// can be passed to json.Unmarshal.
var requestCtor = map[string]messageCtor{`)
	for _, r := range reqs {
		req := strings.TrimSuffix(firstToLower(r), "Request")
		var msg string
		if req == "initialize" {
			msg = `
	Arguments: InitializeRequestArguments{
		// Set the default values specified here: https://microsoft.github.io/debug-adapter-protocol/specification#Requests_Initialize.
		LinesStartAt1:   true,
		ColumnsStartAt1: true,
		PathFormat:      "path",
},
`
		}
		fmt.Fprintf(sb, "\n\t\"%s\":\tfunc() Message { return &%s{%s} },", req, r, msg)
	}
	fmt.Fprint(sb, "\n}")

	fmt.Fprint(sb, `
// Mapping of response commands and corresponding struct constructors that
// can be passed to json.Unmarshal.
var responseCtor = map[string]messageCtor{`)
	for _, r := range resps {
		resp := strings.TrimSuffix(firstToLower(r), "Response")

		fmt.Fprintf(sb, "\n\t\"%s\":\tfunc() Message { return &%s{} },", resp, r)
	}
	fmt.Fprint(sb, "\n}")

	fmt.Fprint(sb, `
// Mapping of event ids and corresponding struct constructors that
// can be passed to json.Unmarshal.
var eventCtor = map[string]messageCtor{`)
	for _, e := range events {
		ev := strings.TrimSuffix(firstToLower(e), "Event")
		fmt.Fprintf(sb, "\n\t\"%s\":\tfunc() Message { return &%s{} },", ev, e)
	}
	fmt.Fprint(sb, "\n}\n")
}

func firstToLower(s string) string {
	r := []rune(s)
	return string(unicode.ToLower(r[0])) + string(r[1:])
}

const preamble = `// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "cmd/gentypes/gentypes.go"; DO NOT EDIT.
// DAP spec: https://microsoft.github.io/debug-adapter-protocol/specification
// See cmd/gentypes/README.md for additional details.

package dap

import "encoding/json"

// Message is an interface that all DAP message types implement with pointer
// receivers. It's not part of the protocol but is used to enforce static
// typing in Go code and provide some common accessors.
//
// Note: the DAP type "Message" (which is used in the body of ErrorResponse)
// is renamed to ErrorMessage to avoid collision with this interface.
type Message interface {
	GetSeq() int
}

// RequestMessage is an interface implemented by all Request-types.
type RequestMessage interface {
	Message
	// GetRequest provides access to the embedded Request.
	GetRequest() *Request
}

// ResponseMessage is an interface implemented by all Response-types.
type ResponseMessage interface {
	Message
	// GetResponse provides access to the embedded Response.
	GetResponse() *Response
}

// EventMessage is an interface implemented by all Event-types.
type EventMessage interface {
	Message
	// GetEvent provides access to the embedded Event.
	GetEvent() *Event
}

// LaunchAttachRequest is an interface implemented by
// LaunchRequest and AttachRequest as they contain shared
// implementation specific arguments that are not part of
// the specification.
type LaunchAttachRequest interface {
	RequestMessage
	// GetArguments provides access to the Arguments map.
	GetArguments() json.RawMessage
}
`

// typesExcludeList is an exclude list of type names we don't want to emit.
var typesExcludeList = map[string]bool{
	// LaunchRequest and AttachRequest arguments can be arbitrary maps.
	// Therefore, this type is not used anywhere.
	"LaunchRequestArguments": true,
	"AttachRequestArguments": true,
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Path to the DAP specification json file is required.")
		fmt.Fprintln(os.Stderr, "gentypes <path/to/debugProtocol.json>")
		os.Exit(1)
	}

	inputFilename := flag.Arg(0)

	if *uFlag {
		if err := updateInput(inputFilename); err != nil {
			log.Fatalf("Failed to update the input file: %v", err)
		}
	}

	inputData, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(inputData, &m); err != nil {
		log.Fatal(err)
	}
	var typeMap map[string]json.RawMessage
	if err := json.Unmarshal(m["definitions"], &typeMap); err != nil {
		log.Fatal(err)
	}

	goTypesIsStruct := make(map[string]bool)
	for typeName, descJson := range typeMap {
		var descMap map[string]json.RawMessage
		if err := json.Unmarshal(descJson, &descMap); err != nil {
			log.Fatal(err)
		}
		_, descMap = maybeParseInheritance(descMap)

		typeJson, ok := descMap["type"]
		if !ok {
			log.Fatal("want description to have 'type', got ", descMap)
		}

		var descTypeString string
		if err := json.Unmarshal(typeJson, &descTypeString); err != nil {
			log.Fatal(err)
		}

		goTypesIsStruct[replaceGoTypename(typeName)] = descTypeString == "object"
	}

	var b strings.Builder
	b.WriteString(preamble)

	typeNames, err := keysInOrder(m["definitions"])
	if err != nil {
		log.Fatal(err)
	}

	var requests, responses, events []string
	for _, typeName := range typeNames {
		if _, ok := typesExcludeList[typeName]; !ok {
			b.WriteString(emitToplevelType(replaceGoTypename(typeName), typeMap[typeName], goTypesIsStruct))
			b.WriteString("\n")
		}

		emitMethodsForType(&b, replaceGoTypename(typeName))
		// Add the typename to the appropriate list.
		if strings.HasSuffix(typeName, "Request") && typeName != "Request" {
			requests = append(requests, typeName)
		}
		if strings.HasSuffix(typeName, "Response") && typeName != "Response" && typeName != "ErrorResponse" {
			responses = append(responses, typeName)
		}
		if strings.HasSuffix(typeName, "Event") && typeName != "Event" {
			events = append(events, typeName)
		}
	}

	// Emit the maps from id to response and event types.
	emitCtor(&b, requests, responses, events)

	wholeFile := []byte(b.String())
	formatted, err := format.Source(wholeFile)
	if err != nil {
		log.Fatal(err)
	}
	if *oFlag == "" {
		fmt.Print(string(formatted))
	} else {
		if err := ioutil.WriteFile(*oFlag, formatted, 0644); err != nil {
			log.Fatalf("Failed to write the generated file: %v", err)
		}
	}
}

func updateInput(inputFilename string) error {
	resp, err := http.Get("https://raw.githubusercontent.com/microsoft/vscode-debugadapter-node/main/debugProtocol.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(inputFilename, data, 0644)
}
