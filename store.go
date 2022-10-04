package main

import (
	"fmt"
	"sort"
	"strings"
)

type (
	Record struct {
		Name    string
		Form    string
		Place   []string
		Photo   string
		ExpDate string
		Tags    []string
		Hash    string
	}
	Store []Record
)

func (r *Record) ContainsAll(words []string) bool {
	for _, w := range words {
		if "" == strings.TrimSpace(w) {
			continue
		}
		if !strings.Contains(r.Hash, w) {
			return false
		}
	}
	return true
}

func (s *Store) Filter(query string) (res Store) {
	qWords := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return ' ' == r || ',' == r
	})
	res = make([]Record, 0)
	for _, r := range *s {
		if r.ContainsAll(qWords) {
			res = append(res, r)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return
}

func (r *Record) AsHTML() (res string) {
	res = fmt.Sprintf(`<b>%s</b> %s
📦 [%s]
🤢️ <code>%s</code>
📆 <b>%s</b>
<a href="%s">&#8205;</a>
`,
		r.Name,
		r.Form,
		strings.Join(r.Place, ", "),
		strings.Join(r.Tags, " "),
		r.ExpDate,
		r.Photo,
	)
	return
}
