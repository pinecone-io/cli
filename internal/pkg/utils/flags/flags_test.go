package flags

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, dir, name, contents string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(contents), 0o600); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	return p
}

func Test_JSONObject_Set_InlineReplace(t *testing.T) {
	var obj JSONObject
	_ = obj.Set(`{"a":1}`)

	if err := obj.Set(`{"b":2}`); err != nil {
		t.Fatalf("error setting inline in JSONObject: %v", err)
	}
	if len(obj) != 1 || obj["b"] != float64(2) {
		t.Fatalf("expected replace from inline, got: %#v", obj)
	}
}

func Test_JSONObject_Set_FileReplace(t *testing.T) {
	var obj JSONObject
	_ = obj.Set(`{"x":1}`)

	dir := t.TempDir()
	path := writeTemp(t, dir, "obj.json", `{"y":2}`)
	if err := obj.Set(path); err != nil {
		t.Fatalf("error setting file path in JSONObject: %v", err)
	}
	if len(obj) != 1 || obj["y"] != float64(2) {
		t.Fatalf("expected replace from file, got: %#v", obj)
	}
}

func Test_JSONObject_Set_EmptyClears(t *testing.T) {
	var obj JSONObject
	_ = obj.Set(`{"a":1}`)

	if err := obj.Set(""); err != nil {
		t.Fatalf("error clearing JSONObject: %v", err)
	}
	if len(obj) != 0 {
		t.Fatalf("expected cleared JSONObject, got: %#v", obj)
	}
}

func Test_Float32List_Set_InlineReplace(t *testing.T) {
	var l Float32List
	_ = l.Set(`[1.5,2.5]`)
	replacement := "[3.5,4.5]"

	if err := l.Set(replacement); err != nil {
		t.Fatalf("error setting inline in Float32List: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from inline, got: %#v", l)
	}
}

func Test_Float32List_Set_FileReplace(t *testing.T) {
	var l Float32List
	_ = l.Set(`[1.5,2.5]`)
	replacement := "[1.5,2.5]"

	dir := t.TempDir()
	path := writeTemp(t, dir, "f32.json", replacement)
	if err := l.Set(path); err != nil {
		t.Fatalf("error setting file in Float32List: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from file, got: %#v", l)
	}
}

func Test_Float32List_Set_EmptyClears(t *testing.T) {
	var l Float32List
	_ = l.Set(`[1,2,3]`)
	if len(l) != 3 {
		t.Fatalf("seed failed")
	}
	if err := l.Set(""); err != nil {
		t.Fatalf("error clearing Float32List: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("expected cleared Float32List, got %#v", l)
	}
}

func Test_UInt32List_Set_InlineReplace(t *testing.T) {
	var l UInt32List
	_ = l.Set(`[1,2,3]`)
	replacement := "[4,5,6]"

	if err := l.Set(replacement); err != nil {
		t.Fatalf("error setting inline in UInt32List: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from inline, got: %#v", l)
	}
}

func Test_UInt32List_Set_FileReplace(t *testing.T) {
	var l UInt32List
	_ = l.Set(`[1,2,3]`)
	replacement := "[4,5,6]"

	dir := t.TempDir()
	path := writeTemp(t, dir, "u32.json", replacement)
	if err := l.Set(path); err != nil {
		t.Fatalf("error setting file in UInt32List: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from file, got: %#v", l)
	}
}

func Test_UInt32List_Set_EmptyClears(t *testing.T) {
	var l UInt32List
	_ = l.Set(`[1,2,3]`)

	if err := l.Set(""); err != nil {
		t.Fatalf("error clearing UInt32List: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("expected cleared UInt32List, got %#v", l)
	}
}

func Test_StringList_Set_InlineReplace(t *testing.T) {
	var l StringList
	_ = l.Set(`["a","b","c"]`)
	replacement := `["d","e","f"]`

	if err := l.Set(replacement); err != nil {
		t.Fatalf("error setting inline in StringList: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from inline, got: %#v", l)
	}
}

func Test_StringList_Set_File(t *testing.T) {
	var l StringList
	_ = l.Set(`["a","b","c"]`)
	replacement := `["d","e","f"]`

	dir := t.TempDir()
	path := writeTemp(t, dir, "s.json", replacement)
	if err := l.Set(path); err != nil {
		t.Fatalf("error setting file in StringList: %v", err)
	}
	if l.String() != replacement {
		t.Fatalf("expected replace from file, got: %#v", l)
	}
}

func Test_StringList_Set_EmptyClears(t *testing.T) {
	var l StringList
	_ = l.Set(`["a","b","c"]`)

	if err := l.Set(""); err != nil {
		t.Fatalf("error clearing StringList: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("expected cleared StringList, got %#v", l)
	}
}
