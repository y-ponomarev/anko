package vm

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/mattn/anko/parser"
)

type testStruct struct {
	script     string
	parseError error
	types      map[string]interface{}
	input      map[string]interface{}
	runError   error
	runOutput  interface{}
	ouput      map[string]interface{}
}

var (
	testVarValue    = reflect.Value{}
	testVarValueP   = &reflect.Value{}
	testVarBool     = true
	testVarBoolP    = &testVarBool
	testVarInt32    = int32(1)
	testVarInt32P   = &testVarInt32
	testVarInt64    = int64(1)
	testVarInt64P   = &testVarInt64
	testVarFloat32  = float32(1)
	testVarFloat32P = &testVarFloat32
	testVarFloat64  = float64(1)
	testVarFloat64P = &testVarFloat64
	testVarString   = "a"
	testVarStringP  = &testVarString
	testVarFunc     = func() int64 { return 1 }
	testVarFuncP    = &testVarFunc

	testVarValueBool    = reflect.ValueOf(true)
	testVarValueInt32   = reflect.ValueOf(int32(1))
	testVarValueInt64   = reflect.ValueOf(int64(1))
	testVarValueFloat32 = reflect.ValueOf(float32(1.1))
	testVarValueFloat64 = reflect.ValueOf(float64(1.1))
	testVarValueString  = reflect.ValueOf("a")
)

func TestIf(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testStruct{
		{script: "if true {}", input: map[string]interface{}{"a": nil}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if true {}", input: map[string]interface{}{"a": true}, runOutput: nil, ouput: map[string]interface{}{"a": true}},
		{script: "if true {}", input: map[string]interface{}{"a": int64(1)}, runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "if true {}", input: map[string]interface{}{"a": float64(1.1)}, runOutput: nil, ouput: map[string]interface{}{"a": float64(1.1)}},
		{script: "if true {}", input: map[string]interface{}{"a": "a"}, runOutput: nil, ouput: map[string]interface{}{"a": "a"}},

		{script: "if true {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if true {a = true}", input: map[string]interface{}{"a": int64(2)}, runOutput: true, ouput: map[string]interface{}{"a": true}},
		{script: "if true {a = 1}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(1), ouput: map[string]interface{}{"a": int64(1)}},
		{script: "if true {a = 1.1}", input: map[string]interface{}{"a": int64(2)}, runOutput: float64(1.1), ouput: map[string]interface{}{"a": float64(1.1)}},
		{script: "if true {a = \"a\"}", input: map[string]interface{}{"a": int64(2)}, runOutput: "a", ouput: map[string]interface{}{"a": "a"}},

		{script: "if a == 1 {a = 1}", input: map[string]interface{}{"a": int64(2)}, runOutput: false, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "if a == 2 {a = 1}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(1), ouput: map[string]interface{}{"a": int64(1)}},
		{script: "if a == 1 {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: false, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "if a == 2 {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},

		{script: "if a == 1 {a = 1} else {a = 3}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(3), ouput: map[string]interface{}{"a": int64(3)}},
		{script: "if a == 2 {a = 1} else {a = 3}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(1), ouput: map[string]interface{}{"a": int64(1)}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3}", input: map[string]interface{}{"a": int64(2)}, runOutput: false, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "if a == 1 {a = 1} else if a == 2 {a = 3}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(3), ouput: map[string]interface{}{"a": int64(3)}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else {a = 4}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(4), ouput: map[string]interface{}{"a": int64(4)}},

		{script: "if a == 1 {a = 1} else {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if a == 2 {a = nil} else {a = 3}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if a == 1 {a = nil} else if a == 3 {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: false, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "if a == 1 {a = 1} else if a == 2 {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},

		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else if a == 4 {a = 4} else {a = 5}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(5), ouput: map[string]interface{}{"a": int64(5)}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else if a == 4 {a = 4} else {a = nil}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else if a == 2 {a = 4} else {a = 5}", input: map[string]interface{}{"a": int64(2)}, runOutput: int64(4), ouput: map[string]interface{}{"a": int64(4)}},
		{script: "if a == 1 {a = 1} else if a == 3 {a = 3} else if a == 2 {a = nil} else {a = 5}", input: map[string]interface{}{"a": int64(2)}, runOutput: nil, ouput: map[string]interface{}{"a": nil}},
	}
	runTests(t, tests)
}

func TestStrings(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testStruct{
		{script: "a", input: map[string]interface{}{"a": 'a'}, runOutput: 'a', ouput: map[string]interface{}{"a": 'a'}},
		{script: "a.b", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support member operation"), ouput: map[string]interface{}{"a": 'a'}},
		{script: "a[0]", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support index operation"), runOutput: nil, ouput: map[string]interface{}{"a": 'a'}},
		{script: "a[0:1]", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support slice operation"), runOutput: nil, ouput: map[string]interface{}{"a": 'a'}},

		{script: "a.b = \"a\"", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support member operation"), runOutput: nil, ouput: map[string]interface{}{"a": 'a'}},
		{script: "a[0] = \"a\"", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support index operation"), runOutput: nil, ouput: map[string]interface{}{"a": 'a'}},
		{script: "a[0:1] = \"a\"", input: map[string]interface{}{"a": 'a'}, runError: fmt.Errorf("type int32 does not support slice operation"), runOutput: nil, ouput: map[string]interface{}{"a": 'a'}},

		{script: "a.b = \"a\"", input: map[string]interface{}{"a": "test"}, runError: fmt.Errorf("type string does not support member operation"), ouput: map[string]interface{}{"a": "test"}},
		{script: "a[0] = \"a\"", input: map[string]interface{}{"a": "test"}, runError: fmt.Errorf("type string does not support index operation for assignment"), ouput: map[string]interface{}{"a": "test"}},
		{script: "a[0:1] = \"a\"", input: map[string]interface{}{"a": "test"}, runError: fmt.Errorf("type string does not support slice operation for assignment"), ouput: map[string]interface{}{"a": "test"}},

		{script: "a", input: map[string]interface{}{"a": "test"}, runOutput: "test", ouput: map[string]interface{}{"a": "test"}},
		{script: "a[0]", input: map[string]interface{}{"a": "test"}, runOutput: 't', ouput: map[string]interface{}{"a": "test"}},
		{script: "a[1]", input: map[string]interface{}{"a": "test"}, runOutput: 'e', ouput: map[string]interface{}{"a": "test"}},
		{script: "a[3]", input: map[string]interface{}{"a": "test"}, runOutput: 't', ouput: map[string]interface{}{"a": "test"}},

		{script: "a[1:0]", input: map[string]interface{}{"a": "test data"}, runError: fmt.Errorf("invalid slice index"), ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[-1:2]", input: map[string]interface{}{"a": "test data"}, runError: fmt.Errorf("index out of range"), ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:-2]", input: map[string]interface{}{"a": "test data"}, runError: fmt.Errorf("index out of range"), ouput: map[string]interface{}{"a": "test data"}},

		{script: "a[0:0]", input: map[string]interface{}{"a": "test data"}, runOutput: "", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:1]", input: map[string]interface{}{"a": "test data"}, runOutput: "t", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:2]", input: map[string]interface{}{"a": "test data"}, runOutput: "te", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:3]", input: map[string]interface{}{"a": "test data"}, runOutput: "tes", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:7]", input: map[string]interface{}{"a": "test data"}, runOutput: "test da", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:8]", input: map[string]interface{}{"a": "test data"}, runOutput: "test dat", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:9]", input: map[string]interface{}{"a": "test data"}, runOutput: "test data", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[0:10]", input: map[string]interface{}{"a": "test data"}, runError: fmt.Errorf("index out of range"), ouput: map[string]interface{}{"a": "test data"}},

		{script: "a[1:1]", input: map[string]interface{}{"a": "test data"}, runOutput: "", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:2]", input: map[string]interface{}{"a": "test data"}, runOutput: "e", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:3]", input: map[string]interface{}{"a": "test data"}, runOutput: "es", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:7]", input: map[string]interface{}{"a": "test data"}, runOutput: "est da", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:8]", input: map[string]interface{}{"a": "test data"}, runOutput: "est dat", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:9]", input: map[string]interface{}{"a": "test data"}, runOutput: "est data", ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[1:10]", input: map[string]interface{}{"a": "test data"}, runError: fmt.Errorf("index out of range"), ouput: map[string]interface{}{"a": "test data"}},

		{script: "a[:2]", input: map[string]interface{}{"a": "test data"}, parseError: fmt.Errorf("syntax error"), ouput: map[string]interface{}{"a": "test data"}},
		{script: "a[2:]", input: map[string]interface{}{"a": "test data"}, parseError: fmt.Errorf("syntax error"), ouput: map[string]interface{}{"a": "test data"}},
	}
	runTests(t, tests)
}

func TestVar(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testStruct{
		{script: "var a = 1", runOutput: int64(1), ouput: map[string]interface{}{"a": int64(1)}},
		{script: "var a, b = 1, 2", runOutput: int64(2), ouput: map[string]interface{}{"a": int64(1), "b": int64(2)}},
		{script: "var a, b, c = 1, 2, 3", runOutput: int64(3), ouput: map[string]interface{}{"a": int64(1), "b": int64(2), "c": int64(3)}},
	}
	runTests(t, tests)
}

func TestMake(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testStruct{
		{script: "make(nilT)", types: map[string]interface{}{"nilT": nil}, runError: fmt.Errorf("invalid type for make")},

		{script: "make(bool)", types: map[string]interface{}{"bool": true}, runOutput: false},
		{script: "make(int32)", types: map[string]interface{}{"int32": int32(1)}, runOutput: int32(0)},
		{script: "make(int64)", types: map[string]interface{}{"int64": int64(1)}, runOutput: int64(0)},
		{script: "make(float32)", types: map[string]interface{}{"float32": float32(1.1)}, runOutput: float32(0)},
		{script: "make(float64)", types: map[string]interface{}{"float64": float64(1.1)}, runOutput: float64(0)},
		{script: "make(string)", types: map[string]interface{}{"string": "a"}, runOutput: ""},
	}
	runTests(t, tests)
}

func TestForLoop(t *testing.T) {
	os.Setenv("ANKO_DEBUG", "1")
	tests := []testStruct{
		{script: "break", runError: fmt.Errorf("Unexpected break statement")},
		{script: "continue", runError: fmt.Errorf("Unexpected continue statement")},

		{script: "for { break }", runOutput: nil},
		{script: "for {a = 1; if a == 1 { break } }", runOutput: nil},
		{script: "a = 1; for { if a == 1 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for { if a == 1 { break }; a++ }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},

		{script: "a = 1; for { a++; if a == 2 { continue } else { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},
		{script: "a = 1; for { a++; if a == 2 { continue }; if a == 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},

		{script: "for a in [1] { if a == 1 { break } }", runOutput: nil},
		{script: "for a in [1, 2] { if a == 2 { break } }", runOutput: nil},
		{script: "for a in [1, 2, 3] { if a == 3 { break } }", runOutput: nil},

		{script: "a = [1]; for b in a { if b == 1 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1)}}},
		{script: "a = [1, 2]; for b in a { if b == 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1), int64(2)}}},
		{script: "a = [1, 2, 3]; for b in a { if b == 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1), int64(2), int64(3)}}},

		{script: "a = [1]; b = 0; for c in a { b = c }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1)}, "b": int64(1)}},
		{script: "a = [1, 2]; b = 0; for c in a { b = c }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1), int64(2)}, "b": int64(2)}},
		{script: "a = [1, 2, 3]; b = 0; for c in a { b = c }", runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1), int64(2), int64(3)}, "b": int64(3)}},

		{script: "a = 1; for a < 2 { a++ }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for a < 3 { a++ }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},

		{script: "a = 1; for nil { a++; if a > 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for nil { a++; if a > 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for true { a++; if a > 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},
		{script: "a = 1; for true { a++; if a > 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(4)}},

		{script: "func x() { return [1] }; for b in x() { if b == 1 { break } }", runOutput: nil},
		{script: "func x() { return [1, 2] }; for b in x() { if b == 2 { break } }", runOutput: nil},
		{script: "func x() { return [1, 2, 3] }; for b in x() { if b == 3 { break } }", runOutput: nil},

		{script: "func x() { a = 1; for { if a == 1 { return } } }; x()", runOutput: nil},
		{script: "func x() { a = 1; for { if a == 1 { return nil } } }; x()", runOutput: nil},
		{script: "func x() { a = 1; for { if a == 1 { return true } } }; x()", runOutput: true},
		{script: "func x() { a = 1; for { if a == 1 { return 1 } } }; x()", runOutput: int64(1)},
		{script: "func x() { a = 1; for { if a == 1 { return 1.1 } } }; x()", runOutput: float64(1.1)},
		{script: "func x() { a = 1; for { if a == 1 { return \"a\" } } }; x()", runOutput: "a"},

		{script: "func x() { for a in [1, 2, 3] { if a == 3 { return } } }; x()", runOutput: nil},
		{script: "func x() { for a in [1, 2, 3] { if a == 1 { continue } } }; x()", runOutput: nil},
		{script: "func x() { for a in [1, 2, 3] { if a == 1 { continue };  if a == 3 { return } } }; x()", runOutput: nil},

		{script: "func x() { return [1, 2] }; a = 1; for i in x() { a++ }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},
		{script: "func x() { return [1, 2, 3] }; a = 1; for i in x() { a++ }", runOutput: nil, ouput: map[string]interface{}{"a": int64(4)}},

		{script: "for a = 1; nil; nil { if a == 1 { break } }", runOutput: nil},
		{script: "for a = 1; nil; nil { if a == 2 { break }; a++ }", runOutput: nil},
		{script: "for a = 1; nil; nil { a++; if a == 3 { break } }", runOutput: nil},

		{script: "for a = 1; a < 1; nil { }", runOutput: nil},
		{script: "for a = 1; a > 1; nil { }", runOutput: nil},
		{script: "for a = 1; a == 1; nil { break }", runOutput: nil},

		{script: "for a = 1; a == 1; a++ { }", runOutput: nil},
		{script: "for a = 1; a < 2; a++ { }", runOutput: nil},
		{script: "for a = 1; a < 3; a++ { }", runOutput: nil},

		{script: "a = 1; for b = 1; a < 1; a++ { }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 2; a++ { }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for b = 1; a < 3; a++ { }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},

		{script: "a = 1; for b = 1; a < 1; a++ {  if a == 1 { continue } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 2; a++ {  if a == 1 { continue } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for b = 1; a < 3; a++ {  if a == 1 { continue } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},

		{script: "a = 1; for b = 1; a < 1; a++ {  if a == 1 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 2; a++ {  if a == 1 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 3; a++ {  if a == 1 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 1; a++ {  if a == 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 2; a++ {  if a == 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for b = 1; a < 3; a++ {  if a == 2 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for b = 1; a < 1; a++ {  if a == 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(1)}},
		{script: "a = 1; for b = 1; a < 2; a++ {  if a == 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(2)}},
		{script: "a = 1; for b = 1; a < 3; a++ {  if a == 3 { break } }", runOutput: nil, ouput: map[string]interface{}{"a": int64(3)}},

		{script: "func x() { for a = 1; a < 3; a++ { if a == 1 { return a } } }; x()", runOutput: int64(1)},
		{script: "func x() { for a = 1; a < 3; a++ { if a == 2 { return a } } }; x()", runOutput: int64(2)},
		{script: "func x() { for a = 1; a < 3; a++ { if a == 3 { return a } } }; x()", runOutput: nil},
		{script: "func x() { for a = 1; a < 3; a++ { if a == 4 { return a } } }; x()", runOutput: nil},

		{script: "func x() { a = 1; for b = 1; a < 3; a++ { if a == 1 { continue } }; return a }; x()", runOutput: int64(3)},
		{script: "func x() { a = 1; for b = 1; a < 3; a++ { if a == 2 { continue } }; return a }; x()", runOutput: int64(3)},
		{script: "func x() { a = 1; for b = 1; a < 3; a++ { if a == 3 { continue } }; return a }; x()", runOutput: int64(3)},
		{script: "func x() { a = 1; for b = 1; a < 3; a++ { if a == 4 { continue } }; return a }; x()", runOutput: int64(3)},

		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{reflect.Value{}}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{reflect.Value{}}, "b": reflect.Value{}}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{nil}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{nil}, "b": nil}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{true}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{true}, "b": true}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{int32(1)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int32(1)}, "b": int32(1)}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{int64(1)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{int64(1)}, "b": int64(1)}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{float32(1.1)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{float32(1.1)}, "b": float32(1.1)}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{float64(1.1)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{float64(1.1)}, "b": float64(1.1)}},

		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(reflect.Value{})}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(reflect.Value{})}, "b": interface{}(reflect.Value{})}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(nil)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(nil)}, "b": interface{}(nil)}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(true)}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(true)}, "b": interface{}(true)}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(int32(1))}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(int32(1))}, "b": interface{}(int32(1))}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(int64(1))}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(int64(1))}, "b": interface{}(int64(1))}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(float32(1.1))}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(float32(1.1))}, "b": interface{}(float32(1.1))}},
		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": []interface{}{interface{}(float64(1.1))}}, runOutput: nil, ouput: map[string]interface{}{"a": []interface{}{interface{}(float64(1.1))}, "b": interface{}(float64(1.1))}},

		{script: "b = 0; for i in a { b = i }", input: map[string]interface{}{"a": interface{}([]interface{}{nil})}, runOutput: nil, ouput: map[string]interface{}{"a": interface{}([]interface{}{nil}), "b": nil}},

		{script: "for i in nil { }", runError: fmt.Errorf("for cannot loop over type interface")},
		{script: "for i in true { }", runError: fmt.Errorf("for cannot loop over type bool")},
		{script: "a = {}; for i in a { }", runError: fmt.Errorf("for cannot loop over type map")},
		{script: "for i in a { }", input: map[string]interface{}{"a": reflect.Value{}}, runError: fmt.Errorf("for cannot loop over type invalid"), ouput: map[string]interface{}{"a": reflect.Value{}}},
		{script: "for i in a { }", input: map[string]interface{}{"a": interface{}(nil)}, runError: fmt.Errorf("for cannot loop over type interface"), ouput: map[string]interface{}{"a": interface{}(nil)}},
		{script: "for i in a { }", input: map[string]interface{}{"a": interface{}(true)}, runError: fmt.Errorf("for cannot loop over type bool"), ouput: map[string]interface{}{"a": interface{}(true)}},

		{script: "a = 1; b = [{\"c\": \"c\"}]; for i in b { a = i }", runOutput: nil, ouput: map[string]interface{}{"a": map[string]interface{}{"c": "c"}, "b": []interface{}{map[string]interface{}{"c": "c"}}}},
		{script: "a = 1; b = {\"x\": [{\"y\": \"y\"}]};  for i in b.x { a = i }", runOutput: nil, ouput: map[string]interface{}{"a": map[string]interface{}{"y": "y"}, "b": map[string]interface{}{"x": []interface{}{map[string]interface{}{"y": "y"}}}}},
	}
	runTests(t, tests)
}

func runTests(t *testing.T, tests []testStruct) {
	var value reflect.Value
loop:
	for _, test := range tests {
		stmts, err := parser.ParseSrc(test.script)
		if err != nil && test.parseError != nil {
			if err.Error() != test.parseError.Error() {
				t.Errorf("ParseSrc error - received: %v expected: %v - script: %v", err, test.parseError, test.script)
				continue
			}
		} else if err != test.parseError {
			t.Errorf("ParseSrc error - received: %v expected: %v - script: %v", err, test.parseError, test.script)
			continue
		}

		env := NewEnv()

		for inputName, inputValue := range test.input {
			err = env.Define(inputName, inputValue)
			if err != nil {
				t.Errorf("Define error: %v - inputName: %v - script: %v", err, inputName, test.script)
				continue loop
			}
		}

		for typeName, typeValue := range test.types {
			err = env.DefineType(typeName, typeValue)
			if err != nil {
				t.Errorf("DefineType error: %v - typeName: %v - script: %v", err, typeName, test.script)
				continue loop
			}
		}

		value, err = Run(stmts, env)
		if err != nil && test.runError != nil {
			if err.Error() != test.runError.Error() {
				t.Errorf("Run error - received: %v expected: %v - script: %v", err, test.runError, test.script)
				continue
			}
		} else if err != test.runError {
			t.Errorf("Run error - received: %v expected: %v - script: %v", err, test.runError, test.script)
			continue
		}
		if !value.IsValid() || !value.CanInterface() {
			if !reflect.DeepEqual(value, test.runOutput) {
				t.Errorf("Run output - received: %#v expected: %#v - script: %v", value, test.runOutput, test.script)
				continue
			}
		} else if value.Kind() == reflect.Func {
			if !reflect.DeepEqual(value, reflect.ValueOf(test.runOutput)) {
				t.Errorf("Run output - received: %#v expected: %#v - script: %v", value, test.runOutput, test.script)
				continue
			}
		} else {
			if !reflect.DeepEqual(value.Interface(), test.runOutput) {
				t.Errorf("Run output - received: %#v expected: %#v - script: %v", value, reflect.ValueOf(test.runOutput), test.script)
				continue
			}
		}

		for outputName, outputValue := range test.ouput {
			value, err = env.Get(outputName)
			if err != nil {
				t.Errorf("Get error: %v - outputName: %v - script: %v", err, outputName, test.script)
				continue loop
			}

			if !value.IsValid() || !value.CanInterface() {
				if !reflect.DeepEqual(value, outputValue) {
					t.Errorf("outputName %v - received: %#v expected: %#v - script: %v", outputName, value, outputValue, test.script)
					continue
				}
			} else if value.Kind() == reflect.Func {
				if !reflect.DeepEqual(value, reflect.ValueOf(outputValue)) {
					t.Errorf("outputName %v - received: %#v expected: %#v - script: %v", outputName, value, outputValue, test.script)
					continue
				}
			} else {
				if !reflect.DeepEqual(value.Interface(), outputValue) {
					t.Errorf("outputName %v - received: %#v expected: %#v - script: %v", outputName, value, reflect.ValueOf(outputValue), test.script)
					continue
				}
			}
		}

		env.Destroy()
	}
}

func TestInterrupts(t *testing.T) {
	scripts := []string{
		`
closeWaitChan()
for {
}
`,
		`
closeWaitChan()
for {
	for {
	}
}
`,
		`
a = []
for i = 0; i < 10000; i++ {
	a += 1
}
closeWaitChan()
for i in a {
}
`,
		`
a = []
for i = 0; i < 10000; i++ {
	a += 1
}
closeWaitChan()
for i in a {
	for j in a {
	}
}
`,
		`
closeWaitChan()
for i = 0; true; nil {
}
`,
		`
closeWaitChan()
for i = 0; true; nil {
	for j = 0; true; nil {
	}
}
`,
	}
	for _, script := range scripts {
		runInterruptTest(t, script)
	}
}

func runInterruptTest(t *testing.T, script string) {
	waitChan := make(chan struct{}, 1)
	closeWaitChan := func() {
		close(waitChan)
	}
	env := NewEnv()
	err := env.Define("closeWaitChan", closeWaitChan)
	if err != nil {
		t.Errorf("Define error: %v", err)
	}

	go func() {
		<-waitChan
		Interrupt(env)
	}()

	_, err = env.Execute(script)
	if err == nil {
		t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
	} else if err.Error() != InterruptError.Error() {
		t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
	}
}

func TestInterruptConcurrency(t *testing.T) {
	var waitGroup sync.WaitGroup
	env := NewEnv()

	waitGroup.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := env.Execute("for { }")
			if err == nil {
				t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
			} else if err.Error() != InterruptError.Error() {
				t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
			}
			waitGroup.Done()
		}()
	}
	time.Sleep(time.Millisecond)
	Interrupt(env)
	waitGroup.Wait()

	_, err := env.Execute("for { }")
	if err == nil {
		t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
	} else if err.Error() != InterruptError.Error() {
		t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
	}

	ClearInterrupt(env)

	_, err = env.Execute("for i = 0; i < 1000; i ++ {}")
	if err != nil {
		t.Errorf("Execute error - received %v expected: %v", err, nil)
	}

	waitGroup.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := env.Execute("for { }")
			if err == nil {
				t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
			} else if err.Error() != InterruptError.Error() {
				t.Errorf("Execute error - received %v expected: %v", err, InterruptError)
			}
			waitGroup.Done()
		}()
	}
	time.Sleep(time.Millisecond)
	Interrupt(env)
	waitGroup.Wait()
}