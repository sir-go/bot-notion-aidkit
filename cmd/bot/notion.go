package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type (
	NtnSelect struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	NtnPage struct {
		Object string `json:"object"`
		Parent struct {
			DatabaseId string `json:"database_id,omitempty"`
		} `json:"parent"`
		Archived   bool `json:"archived,omitempty"`
		Properties struct {
			Name struct {
				Title interface{} `json:"title"`
			} `json:"Название"`
			Form struct {
				Select NtnSelect `json:"select"`
			} `json:"Форма"`
			Place struct {
				MSelect interface{} `json:"multi_select"`
			} `json:"Место"`
			Photo struct {
				Files interface{} `json:"files"`
			} `json:"Фото"`
			ExpDate struct {
				Date struct {
					Start string `json:"start"`
				} `json:"date"`
			} `json:"Срок годности"`
			Tags struct {
				MSelect interface{} `json:"multi_select"`
			} `json:"Тэги"`
		} `json:"properties"`
	}

	NtnSearchResult struct {
		Results []NtnPage `json:"results"`
	}
)

func (p *NtnPage) GetForm() string {
	return p.Properties.Form.Select.Name
}

func (p *NtnPage) GetName() (name string) {
	if nil == p.Properties.Name.Title {
		return
	}
	titles := p.Properties.Name.Title.([]interface{})
	if 1 > len(titles) {
		return
	}
	t0 := titles[0].(map[string]interface{})
	if nil == t0 || nil == t0["plain_text"] {
		return
	}
	return t0["plain_text"].(string)
}

func (p *NtnPage) GetPlaces() (ps []string) {
	if nil == p.Properties.Place.MSelect {
		return
	}
	places := p.Properties.Place.MSelect.([]interface{})
	for _, place := range places {
		if nil == place {
			continue
		}
		placeStruct := place.(map[string]interface{})
		if nil == placeStruct {
			continue
		}

		if nil != placeStruct["name"] {
			ps = append(ps, placeStruct["name"].(string))
		}
	}
	return
}

func (p *NtnPage) GetPhotoHref() (url string) {
	if nil == p.Properties.Photo.Files {
		return
	}
	files := p.Properties.Photo.Files.([]interface{})
	if 1 > len(files) || nil == files[0] {
		return
	}
	fileStruct := files[0].(map[string]interface{})
	if nil == fileStruct {
		return
	}
	fileRecord := fileStruct["file"]
	if nil == fileRecord {
		return
	}
	fileUrl := fileRecord.(map[string]interface{})
	if nil == fileUrl || nil == fileUrl["url"] {
		return
	}
	return fileUrl["url"].(string)
}

func (p *NtnPage) GetExpDate() string {
	return p.Properties.ExpDate.Date.Start
}

func (p *NtnPage) GetTags() (ts []string) {
	if nil == p.Properties.Tags.MSelect {
		return
	}
	tags := p.Properties.Tags.MSelect.([]interface{})
	for _, tag := range tags {
		if nil == tag {
			continue
		}
		tagStruct := tag.(map[string]interface{})
		if nil == tagStruct {
			continue
		}

		if nil != tagStruct["name"] {
			ts = append(ts, tagStruct["name"].(string))
		}
	}
	return
}

func (p *NtnPage) ToRecord() (r Record) {
	r = Record{
		Name:    p.GetName(),
		Photo:   p.GetPhotoHref(),
		ExpDate: p.GetExpDate(),
		Form:    p.GetForm(),
		Place:   p.GetPlaces(),
		Tags:    p.GetTags(),
	}

	r.Hash = strings.ToLower(strings.Join(append(append([]string{
		r.Name,
		r.ExpDate,
		r.Form,
	}, r.Place...), r.Tags...), "|"))
	return
}

func getData() (store Store, err error) {
	var (
		req       *http.Request
		resp      *http.Response
		body      []byte
		searchRes *NtnSearchResult
	)
	client := &http.Client{Timeout: CFG.NotionAPI.Timeout.Duration}
	req, err = http.NewRequest("POST", CFG.NotionAPI.SearchUrl, nil)
	req.Header.Set("authorization", CFG.NotionAPI.Token)
	req.Header.Set("notion-version", CFG.NotionAPI.Version)

	resp, err = client.Do(req)
	eh(err)

	body, err = io.ReadAll(resp.Body)
	eh(err)
	defer func(Body io.ReadCloser) {
		eh(Body.Close())
	}(resp.Body)

	searchRes = new(NtnSearchResult)
	eh(json.Unmarshal(body, searchRes))
	store = make([]Record, 0)
	for _, rec := range searchRes.Results {
		if rec.Object != "page" ||
			rec.Parent.DatabaseId != CFG.NotionAPI.DbId ||
			rec.Archived {
			continue
		}
		store = append(store, rec.ToRecord())
	}
	LOG.Printf("data updated, got %d records", len(store))
	return
}
