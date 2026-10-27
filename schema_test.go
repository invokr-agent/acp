package acp

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// TestTypesFollowSchema breaks when types.go drifts from the vendored v1 schema
// (refresh it with `make schema`): a type is missing, an object's fields or a
// union's tags differ from the schema's, an enum loses a value, or a method in
// meta.json is not wired into a connection.
func TestTypesFollowSchema(t *testing.T) {
	defs := loadDefs(t)
	gotypes := parseTypes(t)
	for name, def := range defs {
		if jsonRPCEnvelope[name] {
			continue
		}
		got, ok := gotypes[strings.ToLower(name)]
		if !ok {
			t.Errorf("%s: no Go type", name)
			continue
		}
		if props := propertyNames(def); props != nil {
			missing, extra := diff(props, got.fields)
			if len(missing) > 0 || (len(extra) > 0 && def["anyOf"] == nil && def["oneOf"] == nil) {
				t.Errorf("%s: fields missing %v, not in schema %v", name, missing, extra)
			}
		}
		// Raw-JSON and flattened types carry no tags or constants to compare.
		if want := constValues(def); want != nil && got.values != nil {
			if missing, _ := diff(want, got.values); len(missing) > 0 {
				t.Errorf("%s: variants or values missing %v", name, missing)
			}
		}
	}
	src := readAll(t, "agent.go", "client.go", "conn.go")
	for _, side := range []string{"agentMethods", "clientMethods", "protocolMethods"} {
		for _, method := range loadMeta(t)[side] {
			if !strings.Contains(src, strconv.Quote(method)) {
				t.Errorf("method %s is not wired into a connection", method)
			}
		}
	}
}

// jsonRPCEnvelope names schema types the JSON-RPC library covers.
var jsonRPCEnvelope = map[string]bool{"RequestId": true, "AgentRequest": true, "AgentResponse": true,
	"AgentNotification": true, "ClientRequest": true, "ClientResponse": true, "ClientNotification": true,
	"ExtRequest": true, "ExtResponse": true, "ExtNotification": true}

type goType struct {
	fields []string // json names, including embedded structs' fields
	values []string // union tags and enum constant values
}

func loadDefs(t *testing.T) map[string]map[string]any {
	var schema struct {
		Defs map[string]map[string]any `json:"$defs"`
	}
	if err := json.Unmarshal([]byte(readAll(t, "schema/v1/schema.json")), &schema); err != nil {
		t.Fatal(err)
	}
	return schema.Defs
}

func loadMeta(t *testing.T) map[string]map[string]string {
	var meta map[string]json.RawMessage
	if err := json.Unmarshal([]byte(readAll(t, "schema/v1/meta.json")), &meta); err != nil {
		t.Fatal(err)
	}
	out := map[string]map[string]string{}
	for k, raw := range meta {
		var methods map[string]string
		if json.Unmarshal(raw, &methods) == nil {
			out[k] = methods
		}
	}
	return out
}

// propertyNames is an object's property names, nil for a type without any.
func propertyNames(def map[string]any) []string {
	props, ok := def["properties"].(map[string]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(props))
	for name := range props {
		out = append(out, name)
	}
	return out
}

// constValues is a union's tags or an enum's values: every const in its variants.
func constValues(def map[string]any) []string {
	var out []string
	for _, key := range []string{"oneOf", "anyOf"} {
		alts, _ := def[key].([]any)
		for _, a := range alts {
			alt, _ := a.(map[string]any)
			if c, ok := alt["const"]; ok {
				out = append(out, jsonString(c))
			}
			props, _ := alt["properties"].(map[string]any)
			for _, p := range props {
				if c, ok := p.(map[string]any)["const"]; ok {
					out = append(out, jsonString(c))
				}
			}
		}
	}
	return out
}

func jsonString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// parseTypes reads types.go: each struct's json field names (embedded structs
// flattened in) and union tags, and each named type's constant values.
func parseTypes(t *testing.T) map[string]*goType {
	file, err := parser.ParseFile(token.NewFileSet(), "types.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	structs := map[string]*ast.StructType{}
	out := map[string]*goType{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gen.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				out[strings.ToLower(s.Name.Name)] = &goType{}
				if st, ok := s.Type.(*ast.StructType); ok {
					structs[s.Name.Name] = st
				}
			case *ast.ValueSpec:
				addConst(out, s)
			}
		}
	}
	for name, st := range structs {
		g := out[strings.ToLower(name)]
		g.fields, g.values = fieldsOf(st, structs), append(g.values, unionTags(st)...)
	}
	return out
}

func addConst(out map[string]*goType, s *ast.ValueSpec) {
	typ, ok := s.Type.(*ast.Ident)
	if !ok || len(s.Values) == 0 {
		return
	}
	expr, neg := s.Values[0], ""
	if u, ok := expr.(*ast.UnaryExpr); ok && u.Op == token.SUB {
		expr, neg = u.X, "-"
	}
	if lit, ok := expr.(*ast.BasicLit); ok {
		v := neg + lit.Value
		if unq, err := strconv.Unquote(v); err == nil {
			v = unq
		}
		if g := out[strings.ToLower(typ.Name)]; g != nil {
			g.values = append(g.values, v)
		}
	}
}

func fieldsOf(st *ast.StructType, structs map[string]*ast.StructType) []string {
	var out []string
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 { // embedded
			if id, ok := f.Type.(*ast.Ident); ok && structs[id.Name] != nil {
				out = append(out, fieldsOf(structs[id.Name], structs)...)
			}
			continue
		}
		if name, _, _ := strings.Cut(structTag(f).Get("json"), ","); name != "" && name != "-" {
			out = append(out, name)
		}
	}
	return out
}

func unionTags(st *ast.StructType) []string {
	var out []string
	for _, f := range st.Fields.List {
		if tag, ok := structTag(f).Lookup("tag"); ok && tag != "" {
			out = append(out, tag)
		}
	}
	return out
}

func structTag(f *ast.Field) reflect.StructTag {
	if f.Tag == nil {
		return ""
	}
	raw, _ := strconv.Unquote(f.Tag.Value)
	return reflect.StructTag(raw)
}

// diff reports what want has that got lacks, and what got has beyond want.
func diff(want, got []string) (missing, extra []string) {
	for _, w := range want {
		if !slices.Contains(got, w) {
			missing = append(missing, w)
		}
	}
	for _, g := range got {
		if !slices.Contains(want, g) {
			extra = append(extra, g)
		}
	}
	return missing, extra
}

func readAll(t *testing.T, paths ...string) string {
	var b strings.Builder
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(data)
	}
	return b.String()
}
