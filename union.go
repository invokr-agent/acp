package acp

import (
	"encoding/json"
	"errors"
	"reflect"
)

// A union is a struct with one pointer field per variant, each tagged with its
// discriminator value (`tag:""` for the variant sent untagged), and an Unknown
// json.RawMessage that keeps any variant this version does not know.

func marshalUnion(key string, u any) ([]byte, error) {
	v := reflect.ValueOf(u).Elem()
	for i := range v.NumField() {
		tag, ok := v.Type().Field(i).Tag.Lookup("tag")
		if f := v.Field(i); ok && !f.IsNil() {
			return withTag(key, tag, f.Interface())
		}
	}
	if raw := v.FieldByName("Unknown").Bytes(); raw != nil {
		return raw, nil
	}
	return nil, errors.New("acp: " + v.Type().Name() + " has no variant set")
}

func unmarshalUnion(key string, data []byte, u any) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	var tag string
	if raw, ok := fields[key]; ok {
		if err := json.Unmarshal(raw, &tag); err != nil {
			return err
		}
	}
	v := reflect.ValueOf(u).Elem()
	v.SetZero()
	for i := range v.NumField() {
		if t, ok := v.Type().Field(i).Tag.Lookup("tag"); ok && t == tag {
			variant := reflect.New(v.Field(i).Type().Elem())
			v.Field(i).Set(variant)
			return json.Unmarshal(data, variant.Interface())
		}
	}
	v.FieldByName("Unknown").SetBytes(append([]byte(nil), data...))
	return nil
}

// withTag encodes variant, an object, with key set to tag; "" leaves it untagged.
func withTag(key, tag string, variant any) ([]byte, error) {
	data, err := json.Marshal(variant)
	if err != nil || tag == "" {
		return data, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	fields[key], _ = json.Marshal(tag)
	return json.Marshal(fields)
}
