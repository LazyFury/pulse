package ui

type IPage interface {
	GetTitle() string
	GetWidgets() []IWidget
	GetComponentName() string
}

type Page struct {
	Title     string    `json:"title"`
	Component string    `json:"component"`
	Widgets   []IWidget `json:"widgets"` //暂未实现，处理统计页面，form 和 table 有时间改造成统一的
	Table     *Table    `json:"table"`   // 特殊处理 for TableView.vue
}

var _ IPage = (*Page)(nil)

func (p *Page) GetTitle() string {
	return p.Title
}

func (p *Page) GetWidgets() []IWidget {
	return p.Widgets
}

func (p *Page) GetComponentName() string {
	return p.Component
}
