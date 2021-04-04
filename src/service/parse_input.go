package service

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser is able to convert an input string into a useful type.
// Each key is a string referring to a type.
// Each value is a function that is able to turn a string into the key's type
// (as an interface) or will return an error.
type Parser map[string]func(string) (interface{}, error)

// ParserBasic returns a Parser for input parsing.
func ParserBasic() Parser {
	var parsers = map[string]func(string) (interface{}, error){
		"string": func(input string) (interface{}, error) {
			return input, nil
		},

		"int": func(input string) (interface{}, error) {
			val, err := strconv.Atoi(input)
			if err != nil {
				return nil, err
			}
			return val, nil
		},

		"bool": func(input string) (interface{}, error) {
			if strings.ToLower(input) == "true" {
				return true, nil
			}

			if strings.ToLower(input) == "false" {
				return false, nil
			}

			return false, fmt.Errorf("parsing error")
		},
	}
	return parsers
}

// ParseInput will utilise parsers to parse an input.
func ParseInput(parser Parser, tokens []string, parameters []string) ([]interface{}, error) {

	if len(tokens) != len(parameters) {
		if len(tokens) < len(parameters) {
			return nil, fmt.Errorf("err")
		}

		if len(parameters) > 0 {
			if parameters[len(parameters)-1] != "string" {
				return nil, fmt.Errorf("err")
			}

			tokens[len(parameters)-1] = strings.Join(tokens[len(parameters)-1:], " ")
			tokens = tokens[0:len(parameters)]
		}
	}

	outputs := []interface{}{}
	for i := 0; i < len(tokens) && i < len(parameters); i++ {
		token := tokens[i]
		parameter := parameters[i]

		typeParser, ok := parser[parameter]
		if !ok {
			return nil, fmt.Errorf("unsupported")
		}

		value, err := typeParser(token)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, value)
	}

	return outputs, nil
}
