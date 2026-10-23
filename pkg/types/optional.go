package types

import (
	"encoding/json"
)

// A field of a partial update that can also be cleared. Set reports whether the
// field came at all. When it did, a nil Value means it came empty and clears it.
type Optional[T any] struct {
	Set   bool
	Value *T
}

// Only called for a field present in the JSON, so a missing field stays unset
// and an explicit null is set with a nil Value.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	o.Value = &value
	return nil
}
