package lint

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Links []*Link

func convMapToLink(m map[string]any) (*Link, error) {
	lk := &Link{}
	if e, ok := m["title"]; ok {
		f, ok := e.(string)
		if !ok {
			return nil, errors.New("title must be a string")
		}
		lk.Title = f
	}
	if e, ok := m["link"]; ok {
		f, ok := e.(string)
		if !ok {
			return nil, errors.New("link must be a string")
		}
		lk.Link = f
	}
	return lk, nil
}

func convAnyToLink(a any) (*Link, error) {
	switch d := a.(type) {
	case string:
		return &Link{
			Link: d,
		}, nil
	case map[string]any:
		return convMapToLink(d)
	}
	return nil, errors.New("link must be either a string or map[string]any")
}

func (ls *Links) UnmarshalJSON(b []byte) error {
	var a any
	if err := json.Unmarshal(b, &a); err != nil {
		return fmt.Errorf("unmarshal bytes as JSON: %w", err)
	}
	if a == nil {
		return nil
	}
	switch b := a.(type) {
	case []any:
		links := make([]*Link, len(b))
		for i, c := range b {
			lk, err := convAnyToLink(c)
			if err != nil {
				return err
			}
			links[i] = lk
		}
		*ls = links
		return nil
	case map[string]any:
		links := make([]*Link, 0, len(b))
		for k, v := range b {
			s, ok := v.(string)
			if !ok {
				return errors.New("link must be a string")
			}
			links = append(links, &Link{
				Title: k,
				Link:  s,
			})
		}
		*ls = links
		return nil
	}
	return nil
}

type Link struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"link"`
}
