package acp

import (
	"encoding/json"
	"errors"
)

// marshalTagged encodes v (an object) with its discriminator key added.
func marshalTagged(key, tag string, v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	if fields == nil {
		fields = map[string]json.RawMessage{}
	}
	fields[key], _ = json.Marshal(tag)
	return json.Marshal(fields)
}

func marshalUnknown(raw json.RawMessage) ([]byte, error) {
	if raw == nil {
		return nil, errors.New("acp: union has no variant set")
	}
	return raw, nil
}

// readTag returns the discriminator's value, "" when the key is absent.
func readTag(data []byte, key string) (string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return "", err
	}
	var tag string
	if raw, ok := fields[key]; ok {
		if err := json.Unmarshal(raw, &tag); err != nil {
			return "", err
		}
	}
	return tag, nil
}
