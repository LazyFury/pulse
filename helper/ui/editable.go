package ui

type EditableWidget struct {
	Widget
	Prop string // property
	Type string // input, textarea
}

var _ IWidget = (*EditableWidget)(nil)
