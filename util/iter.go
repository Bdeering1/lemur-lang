package util

import (
	"iter"
)

func Collect[T, E any](seq iter.Seq[E]) (s []T, _ bool) {
	for e := range seq {
		v, ok := any(e).(T)
		if !ok { return nil, false }
		s = append(s, v)
	}

	return s, true
}
