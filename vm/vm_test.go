package vm

import (
	"fmt"
	"testing"

	"compiler/ast"
	"compiler/compiler"
	"compiler/lexer"
	"compiler/object"
	"compiler/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

type vmTestCase struct {
	input    string
	expected any
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		vm := New(comp.Bytecode())
		vmerr := vm.Run()
		if vmerr != nil {
			t.Fatalf("vm rrror: %s", vmerr)
		}

		stackElem := vm.StackTop()

		testExpectedObject(t, tt.expected, stackElem)

	}
}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not integer. got=%t (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2}, // TODO: edit the expected field to reflect the behavior of a fully integerated addition machine
	}
	runVmTests(t, tests)
}
