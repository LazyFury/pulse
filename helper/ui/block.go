package ui

type Block struct {
	Widget
	Header Widget
	Footer Widget
}

var _ IWidget = (*Block)(nil)
