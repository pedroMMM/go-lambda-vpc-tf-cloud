package main

import (
	"fmt"
	"strings"
)

type diff struct {
	removed []string
	added   []string
}

func (d diff) NoChanges() bool {
	return len(d.added) == 0 && len(d.removed) == 0
}

func (d diff) ToString() string {
	if d.NoChanges() {
		return "no changes"
	}

	stringify := func(in []string, action string) string {
		if len(in) == 0 {
			return fmt.Sprintf("none were %s", action)
		}

		return fmt.Sprintf("%s: %s", action, strings.Join(in, ", "))
	}

	return fmt.Sprintf("%s\n%s", stringify(d.removed, "removed"), stringify(d.added, "added"))
}

func calculateDiff(old, new []string) diff {
	counter := make(map[string]int)
	out := diff{
		removed: make([]string, 0),
		added:   make([]string, 0),
	}

	for _, e := range old {
		e := e
		counter[e] = -1
	}

	for _, e := range new {
		e := e
		if v, ok := counter[e]; ok {
			counter[e] = v + 1
		} else {
			out.added = append(out.added, e)
		}
	}

	for k, v := range counter {
		k := k
		if v == -1 {
			out.removed = append(out.removed, k)
		}
	}

	return out
}
