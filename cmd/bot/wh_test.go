package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_whItem_matchHash(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		item  whItem
		want  bool
	}{
		{"empty-str",
			[]string{""},
			whItem{},
			false},
		{"empty-array",
			[]string{},
			whItem{},
			false},
		{"single",
			[]string{"живот"},
			whItem{Tags: []string{"живот"}},
			true},
		{"many-on",
			[]string{"живот", "   голова", "  "},
			whItem{Tags: []string{"живот"}},
			true},
		{"many-many",
			[]string{"живот ", "голова"},
			whItem{Tags: []string{"отравление", "живот"}},
			true},
		{"many-many-name-true",
			[]string{"живот ", "голова"},
			whItem{Name: "супрастин"},
			false},
		{"many-many-name-false",
			[]string{"живот ", "супрастин"},
			whItem{Name: "супрастин"},
			true},
		{"many-many-exp-date",
			[]string{"живот ", " ", "2023 "},
			whItem{ExpDate: "2023-06-05"},
			true},
		{"many-many-place",
			[]string{"кухня"},
			whItem{
				Place: []string{"кухня", "кладовка"},
				Tags:  []string{"отравление", "живот"}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.updateHash()
			if got := tt.item.matchHash(tt.words); got != tt.want {
				t.Errorf("matchHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_whItem_updateHash(t *testing.T) {
	tests := []struct {
		name string
		item *whItem
		want string
	}{
		{
			"empty", &whItem{}, "||",
		},
		{
			"tags",
			&whItem{
				Tags: []string{"tag0", "tag1"},
			},
			"|||tag0|tag1",
		},
		{
			"filled-0",
			&whItem{
				Tags: []string{"tag0", "tag1"},
				Form: "таблетки",
			},
			"||таблетки|tag0|tag1",
		},
		{
			"filled-1",
			&whItem{
				Tags:  []string{"tag0", "tag1"},
				Place: []string{"холодильник", "кладовка"},
				Form:  "таблетки",
			},
			"||таблетки|холодильник|кладовка|tag0|tag1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.updateHash()
			if got := tt.item.Hash; got != tt.want {
				t.Errorf("whItem.updateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNotionPage(t *testing.T) {
	tests := []struct {
		name     string
		pageJson string
		want     whItem
	}{
		{
			"null",
			`{
			"object":"page","parent":{"database_id":"some-db-id"},
			"properties":{
				"Место":{"multi_select":[]},
				"Название":{"title":[]},
				"Срок годности":{"date":{"start":""}},
				"Тэги":{"multi_select":[]},
				"Форма":{"select":{}},
				"Фото":{"files":[]}}
			}`,
			whItem{Place: []string{}, Tags: []string{}, Hash: "||"},
		},
		{
			"full",
			`{
			"object":"page","parent":{"database_id":"some-db-id"},
			"properties":{
				"Место":{"multi_select":[{"name":"холодильник"}]},
				"Название":{"title":[{"plain_text":"Уголь активированный"}]},
				"Срок годности":{"date":{"start":"2023-06-01"}},
				"Тэги":{"multi_select":[{"color":"brown","name":"живот"},{"color":"green","name":"отравление"}]},
				"Форма":{"select":{"color":"default","name":"таблетки"}},
				"Фото":{"files":[{"file":{"url":"http://some-image-url"}}]}}
			}`,
			whItem{
				Name:    "Уголь активированный",
				Form:    "таблетки",
				Place:   []string{"холодильник"},
				Photo:   "http://some-image-url",
				ExpDate: "2023-06-01",
				Tags:    []string{"живот", "отравление"},
				Hash:    "уголь активированный|2023-06-01|таблетки|холодильник|живот|отравление",
			},
		},
	}
	for _, tt := range tests {
		page := NtnPage{}
		if err := json.Unmarshal([]byte(tt.pageJson), &page); err != nil {
			t.Error(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNotionPage(page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNotionPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
