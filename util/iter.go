package util

import (
	"iter"
)

func CollectAs[T, E any](seq iter.Seq[E]) (s []T, _ bool) {
	for el := range seq {
		v, ok := any(el).(T)
		if !ok { return nil, false }
		s = append(s, v)
	}

	return s, true
}
