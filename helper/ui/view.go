package ui

type View struct {
	Widget
}

var _ IWidget = (*View)(nil)
