package main

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
	if len(titles) < 1 {
		return
	}
	t0 := titles[0].(map[string]interface{})
	if nil == t0 || nil == t0["plain_text"] {
		return
	}
	return t0["plain_text"].(string)
}

func (p *NtnPage) GetPlaces() (ps []string) {
	ps = make([]string, 0)
	if nil == p.Properties.Place.MSelect {
		return
	}
	places := p.Properties.Place.MSelect.([]interface{})
	for _, place := range places {
		if nil == place {
			continue
		}
		placeStruct := place.(map[string]interface{})
		if nil != placeStruct["name"] {
			if sVal := placeStruct["name"].(string); sVal != "" {
				ps = append(ps, placeStruct["name"].(string))
			}
		}
	}
	return
}

func (p *NtnPage) GetPhotoHref() (url string) {
	if nil == p.Properties.Photo.Files {
		return
	}
	files := p.Properties.Photo.Files.([]interface{})
	if len(files) < 1 || nil == files[0] {
		return
	}
	fileStruct := files[0].(map[string]interface{})
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

func (p *NtnPage) GetTags() (res []string) {
	res = make([]string, 0)
	if p.Properties.Tags.MSelect == nil {
		return
	}
	tags := p.Properties.Tags.MSelect.([]interface{})
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		tagStruct := tag.(map[string]interface{})
		if tagStruct == nil {
			continue
		}

		if tagStruct["name"] != nil {
			if sVal := tagStruct["name"].(string); sVal != "" {
				res = append(res, sVal)
			}
		}
	}
	return
}
