package hal

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Link struct {
	Href   string `json:"href"`
	Title  string `json:"title,omitempty"`
	Method string `json:"method,omitempty"`
}

type Resource struct {
	Type     string          `json:"_type,omitempty"`
	Links    map[string]Link `json:"_links"`
	Embedded map[string]any  `json:"_embedded,omitempty"`
	Payload
}

type Payload interface{}

type NewResourceOption func(r *Resource)

func NewResource(res any, selfLink string, opts ...NewResourceOption) *Resource {
	r := &Resource{
		Links: map[string]Link{
			"self": {
				Href: selfLink,
			},
		},
		Payload: res,
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

func WithType(t string) NewResourceOption {
	return func(r *Resource) {
		r.Type = t
	}
}

func WithSelfTitle(t string) NewResourceOption {
	return func(r *Resource) {
		if l, ok := r.Links["self"]; ok {
			l.Title = t

			r.Links["self"] = l
		}
	}
}

func (r *Resource) Embed(relation string, resource *Resource) {
	if r.Embedded == nil {
		r.Embedded = make(map[string]any, 1)
	}

	r.Embedded[relation] = resource
}

func (r *Resource) AddLink(relation string, link Link) {
	if r.Links == nil {
		r.Links = make(map[string]Link, 1)
	}

	r.Links[relation] = link
}

func (r *Resource) SetType(t string) {
	r.Type = t
}

func (r *Resource) MarshalJSON() ([]byte, error) {

	rootData := struct {
		Type     string          `json:"_type,omitempty"`
		Links    map[string]Link `json:"_links"`
		Embedded map[string]any  `json:"_embedded,omitempty"`
	}{}

	rootData.Type = r.Type

	if len(r.Links) > 0 {
		rootData.Links = r.Links
	}

	if len(r.Embedded) > 0 {
		rootData.Embedded = r.Embedded
	}

	dataRoot, err := json.Marshal(rootData)
	if err != nil {
		return nil, fmt.Errorf("error marshal root data. %w", err)
	}

	dataChild, err := json.Marshal(r.Payload)
	if err != nil {
		return nil, fmt.Errorf("error marshal payload data. %w", err)
	}

	if len(dataRoot) == 2 {
		return dataChild, nil
	}

	var b bytes.Buffer

	if cap := b.Cap(); cap < (len(dataRoot) + len(dataChild)) {
		b.Grow((len(dataRoot) + len(dataChild)) - cap)
	}

	b.Write(dataRoot[:len(dataRoot)-1])

	if len(dataChild) != 2 {
		b.Write([]byte(`,`))
	}

	b.Write(dataChild[1:])

	return b.Bytes(), nil
}
