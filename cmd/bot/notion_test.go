package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNtnPage_GetExpDate(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"empty",
			`{"properties": {"Срок годности": {"date": {"start": ""}}}}`, ""},
		{"valid",
			`{"properties": {"Срок годности": {"date": {"start": "2022-11-10"}}}}`, "2022-11-10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetExpDate(); got != tt.want {
				t.Errorf("GetExpDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNtnPage_GetForm(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"empty",
			`{"properties": {"Форма": {"select": {"name": ""}}}}`, ""},
		{"valid",
			`{"properties": {"Форма": {"select": {"name": "таблетки"}}}}`, "таблетки"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetForm(); got != tt.want {
				t.Errorf("GetForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNtnPage_GetName(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"empty",
			`{"properties": {"Название": {"title": []}}}`, ""},
		{"null",
			`{"properties": {"Название": {"title": null}}}`, ""},
		{"null-value",
			`{"properties": {"Название": {"title": [{"plain_text": null}]}}}`, ""},
		{"valid0",
			`{"properties": {"Название": {"title": [{"plain_text": "анальгин"}]}}}`, "анальгин"},
		{"valid1",
			`{"properties": {"Название": {"title": [{"plain_text": "анальгин"}, {"plain_text": "аспирин"}]}}}`, "анальгин"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNtnPage_GetPlaces(t *testing.T) {
	tests := []struct {
		name string
		json string
		want []string
	}{
		{"null",
			`{"properties": {"Место": {"multi_select": null}}}`, []string{}},
		{"empty",
			`{"properties": {"Место": {"multi_select": [{"name": ""}]}}}`, []string{}},
		{"null-value",
			`{"properties": {"Место": {"multi_select": [{"name": null}]}}}`, []string{}},
		{"null-array",
			`{"properties": {"Место": {"multi_select": [null]}}}`, []string{}},
		{"valid0",
			`{"properties": {"Место": {"multi_select": [{"name": "кладовка"}]}}}`, []string{"кладовка"}},
		{"valid1",
			`{"properties": {"Место": {"multi_select": [{"name": "холодильник"}, {"name": "кладовка"}]}}}`, []string{"холодильник", "кладовка"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetPlaces(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNtnPage_GetPhotoHref(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{"no-photo-attr",
			`{"properties": {}}`, ""},
		{"no-files-attr",
			`{"properties": {"Фото": {}}}`, ""},
		{"empty-files-array",
			`{"properties": {"Фото": {"files": []}}}`, ""},
		{"no-valid-files",
			`{"properties": {"Фото": {"files": [{}]}}}`, ""},
		{"no-url",
			`{"properties": {"Фото": {"files": [{"file": {}}]}}}`, ""},
		{"empty-url",
			`{"properties": {"Фото": {"files": [{"file": {"url": ""}}]}}}`, ""},
		{"valid",
			`{"properties": {"Фото": {"files": [{"file": {"url": "http://some-url"}}]}}}`, "http://some-url"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetPhotoHref(); got != tt.want {
				t.Errorf("GetPhotoHref() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNtnPage_GetTags(t *testing.T) {
	tests := []struct {
		name string
		json string
		want []string
	}{
		{"null",
			`{"properties": {"Тэги": {"multi_select": null}}}`, []string{}},
		{"empty",
			`{"properties": {"Тэги": {"multi_select": []}}}`, []string{}},
		{"null-name",
			`{"properties": {"Тэги": {"multi_select": [null]}}}`, []string{}},
		{"empty-name",
			`{"properties": {"Тэги": {"multi_select": [{"name": ""}]}}}`, []string{}},
		{"valid0",
			`{"properties": {"Тэги": {"multi_select": [{"name": "давление"}]}}}`, []string{"давление"}},
		{"valid1",
			`{"properties": {"Тэги": {"multi_select": [{"name": "живот"}, {"name": "голова"}]}}}`, []string{"живот", "голова"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var page NtnPage
			if err := json.Unmarshal([]byte(tt.json), &page); err != nil {
				t.Error(err)
			}
			if got := page.GetTags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
