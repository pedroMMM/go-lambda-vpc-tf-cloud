package main

import (
	"fmt"
	"strings"
)

type diff struct {
	removed   []string
	added     []string
	message   *string
	noChanges *bool
}

func (d diff) NoChanges() bool {
	if d.noChanges != nil {
		return *d.noChanges
	}

	tmp := len(d.added) == 0 && len(d.removed) == 0
	d.noChanges = &tmp
	return *d.noChanges
}

func (d diff) ToString() string {
	if d.message != nil {
		return *d.message
	}

	if d.NoChanges() {
		tmp := "no changes"
		d.message = &tmp
		return *d.message
	}

	stringify := func(in []string, action string) string {
		if len(in) == 0 {
			return fmt.Sprintf("none were %s", action)
		}

		return fmt.Sprintf("%s: %s", action, strings.Join(in, ", "))
	}

	tmp := fmt.Sprintf("%s\n%s", stringify(d.removed, "removed"), stringify(d.added, "added"))
	d.message = &tmp
	return *d.message
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
