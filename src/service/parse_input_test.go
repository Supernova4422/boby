package service

import (
	"errors"
	"testing"
)

func TestStringJoin(t *testing.T) {
	parser := ParserBasic()
	params := []string{"string"}
	tokens := []string{"hello", "world"}
	val, err := ParseInput(parser, tokens, params)
	if err != nil {
		t.Fail()
	}

	if len(val) != 1 || val[0].(string) != "hello world" {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	parser := ParserBasic()
	params := []string{"int"}
	tokens := []string{"100"}
	val, err := ParseInput(parser, tokens, params)
	if err != nil {
		t.Fail()
	}

	if len(val) != 1 || val[0].(int) != 100 {
		t.Fail()
	}
}

func TestBadInt(t *testing.T) {
	parser := ParserBasic()
	params := []string{"int"}
	tokens := []string{"hundred"}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}

func TestBadBool(t *testing.T) {
	parser := ParserBasic()
	params := []string{"bool"}
	tokens := []string{"tama"}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}

func TestBool(t *testing.T) {
	parser := ParserBasic()
	params := []string{"bool", "bool"}
	tokens := []string{"true", "false"}
	val, err := ParseInput(parser, tokens, params)
	if err != nil {
		t.Fail()
	}

	if len(val) != 2 || !val[0].(bool) || val[1].(bool) {
		t.Fail()
	}
}

func TestUnsupported(t *testing.T) {
	parser := ParserBasic()
	params := []string{"unsupported"}
	tokens := []string{"unsupported"}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}

func TestPropogateErrs(t *testing.T) {
	parser := ParserBasic()
	parser["err"] = func(string) (interface{}, error) {
		return nil, errors.New("error")
	}

	params := []string{"err"}
	tokens := []string{"err"}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}

func TestDontJoinBool(t *testing.T) {
	parser := ParserBasic()
	params := []string{"bool"}
	tokens := []string{"true", "false"}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}
func TestTooFew(t *testing.T) {
	parser := ParserBasic()
	params := []string{"int"}
	tokens := []string{}
	_, err := ParseInput(parser, tokens, params)
	if err == nil {
		t.Fail()
	}
}
