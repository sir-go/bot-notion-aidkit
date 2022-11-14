package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	zlog "github.com/rs/zerolog/log"
)

type (
	whItem struct {
		Name    string
		Form    string
		Place   []string
		Photo   string
		ExpDate string
		Tags    []string
		Hash    string
	}
	Warehouse struct {
		items       []whItem
		updInterval time.Duration
		dataExpTime time.Time
	}
)

func NewWarehouse(updateInterval time.Duration) *Warehouse {
	return &Warehouse{
		items:       make([]whItem, 0),
		dataExpTime: time.Now(),
		updInterval: updateInterval,
	}
}

func (item *whItem) matchHash(words []string) bool {
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		if strings.Contains(item.Hash, word) {
			return true
		}
	}
	return false
}

func (item *whItem) updateHash() {
	item.Hash = strings.ToLower(strings.Join(append(append([]string{
		item.Name,
		item.ExpDate,
		item.Form,
	}, item.Place...), item.Tags...), "|"))
}

func parseNotionPage(page NtnPage) whItem {
	r := whItem{
		Name:    page.GetName(),
		Photo:   page.GetPhotoHref(),
		ExpDate: page.GetExpDate(),
		Form:    page.GetForm(),
		Place:   page.GetPlaces(),
		Tags:    page.GetTags(),
	}
	r.updateHash()
	return r
}

func (wh *Warehouse) update(nCfg CfgNotionAPI, force bool) (err error) {
	if time.Now().Before(wh.dataExpTime) && len(wh.items) > 0 && !force {
		return
	}

	var (
		req       *http.Request
		resp      *http.Response
		body      []byte
		searchRes *NtnSearchResult
	)
	zlog.Debug().Msg("get new data from notion")
	client := &http.Client{Timeout: nCfg.Timeout.Duration}
	req, err = http.NewRequest("POST", nCfg.SearchUrl, nil)
	req.Header.Set("authorization", nCfg.Token)
	req.Header.Set("notion-version", nCfg.Version)

	if resp, err = client.Do(req); err != nil {
		zlog.Error().Err(err).Msg("notion request")
		return
	}

	if body, err = io.ReadAll(resp.Body); err != nil {
		zlog.Error().Err(err).Msg("read request body")
		return
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			zlog.Error().Err(err).Msg("close request body")
		}
	}(resp.Body)

	searchRes = new(NtnSearchResult)
	if err = json.Unmarshal(body, searchRes); err != nil {
		zlog.Error().Err(err).Msg("unmarshal search result")
		return
	}

	items := make([]whItem, 0)
	for _, rec := range searchRes.Results {
		if rec.Object != "page" ||
			rec.Parent.DatabaseId != nCfg.DbId ||
			rec.Archived {
			continue
		}
		items = append(items, parseNotionPage(rec))
	}
	zlog.Debug().Int("amount", len(items)).Msg("got records")
	wh.items = items
	wh.dataExpTime = time.Now().Add(wh.updInterval)
	return
}

func (wh *Warehouse) Query(query string) []whItem {
	words := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return ' ' == r || ',' == r
	})

	res := make([]whItem, 0)
	for _, r := range wh.items {
		if r.matchHash(words) {
			res = append(res, r)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

func (item *whItem) AsHTML() (res string) {
	res = fmt.Sprintf(`<b>%s</b> %s
📦 [%s]
🤢️ <code>%s</code>
📆 <b>%s</b>
<a href="%s">&#8205;</a>
`,
		item.Name,
		item.Form,
		strings.Join(item.Place, ", "),
		strings.Join(item.Tags, " "),
		item.ExpDate,
		item.Photo,
	)
	return
}
