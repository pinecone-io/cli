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

func TestJSONObject_Set_InlineReplace(t *testing.T) {
	var obj JSONObject
	// seed object
	_ = obj.Set(`{"a":1}`)

	// inline should replace (not merge)
	if err := obj.Set(`{"b":2}`); err != nil {
		t.Fatalf("Set inline: %v", err)
	}
	if len(obj) != 1 || obj["b"] != float64(2) {
		t.Fatalf("expected replace on inline, got: %#v", obj)
	}
}

func TestJSONObject_Set_EmptyClears(t *testing.T) {
	var obj JSONObject
	_ = obj.Set(`{"a":1}`)
	if len(obj) == 0 {
		t.Fatalf("seed failed")
	}
	// empty clears
	if err := obj.Set(""); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if len(obj) != 0 {
		t.Fatalf("expected cleared, got: %#v", obj)
	}
}

func TestJSONObject_Set_FileReplace(t *testing.T) {
	var obj JSONObject
	_ = obj.Set(`{"x":1}`)

	dir := t.TempDir()
	path := writeTemp(t, dir, "obj.json", `{"y":2}`)
	if err := obj.Set("@" + path); err != nil {
		t.Fatalf("Set file: %v", err)
	}
	if len(obj) != 1 || obj["y"] != float64(2) {
		t.Fatalf("expected replace from file, got: %#v", obj)
	}
}

func TestFloat32List_Set_EmptyClears(t *testing.T) {
	var l Float32List
	_ = l.Set(`[1,2,3]`)
	if len(l) != 3 {
		t.Fatalf("seed failed")
	}
	if err := l.Set(""); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if len(l) != 0 {
		t.Fatalf("expected cleared, got %v", l)
	}
}

func TestFloat32List_Set_File(t *testing.T) {
	var l Float32List
	dir := t.TempDir()
	path := writeTemp(t, dir, "f32.json", `[1.5,2.5]`)
	if err := l.Set("@" + path); err != nil {
		t.Fatalf("Set file: %v", err)
	}
	if len(l) != 2 || l[0] != 1.5 || l[1] != 2.5 {
		t.Fatalf("unexpected list: %#v", l)
	}
}

func TestUInt32List_Set_File(t *testing.T) {
	var l UInt32List
	dir := t.TempDir()
	path := writeTemp(t, dir, "u32.json", `[10,20]`)
	if err := l.Set("@" + path); err != nil {
		t.Fatalf("Set file: %v", err)
	}
	if len(l) != 2 || l[0] != 10 || l[1] != 20 {
		t.Fatalf("unexpected list: %#v", l)
	}
}

func TestStringList_Set_File(t *testing.T) {
	var l StringList
	dir := t.TempDir()
	path := writeTemp(t, dir, "s.json", `["a","b"]`)
	if err := l.Set("@" + path); err != nil {
		t.Fatalf("Set file: %v", err)
	}
	if len(l) != 2 || l[0] != "a" || l[1] != "b" {
		t.Fatalf("unexpected list: %#v", l)
	}
}
