package types

import (
	"encoding/json"

	"github.com/Meirlan28/myapp/internal/core/optional"
)

type Optional[T any] struct {
	optional.Optional[T]
}

func (n *Optional[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Optional[T]) ToDomain() optional.Optional[T] {
	return optional.Optional[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
