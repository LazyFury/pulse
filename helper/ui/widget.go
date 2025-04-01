package ui

type IWidget interface {
	GetName() string
}

type Widget struct {
	Name        string         `json:"name"`
	Label       string         `json:"label"`
	Description string         `json:"description"`
	Children    []Widget       `json:"children"` // children
	Params      map[string]any `json:"params"`   // params
}

var _ IWidget = (*Widget)(nil)

func (w *Widget) GetName() string {
	return w.Name
}
