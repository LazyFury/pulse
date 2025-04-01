package ui

type TableColumn struct {
	Prop     string         `json:"prop"`
	Label    string         `json:"label"`
	Width    string         `json:"width"`
	Props    map[string]any `json:"props"`
	Sortable bool           `json:"sortable"`
}

type TableActionType string

var (
	TableActionTypeForm = TableActionType("form")
	TableActionTypeApi  = TableActionType("api")
	TableActionTypeLink = TableActionType("link")
)

type Confirm struct {
	Title          string `json:"title"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	ConfirmBtnText string `json:"confirmButtonText"`
	CancelBtnText  string `json:"cancelButtonText"`
}

type Action struct {
	Label string          `json:"label"`
	Icon  string          `json:"icon"`
	Props map[string]any  `json:"props"`
	Type  TableActionType `json:"type"`

	FormKey string `json:"form_key"`

	ApiKey    string              `json:"api_key"`
	ParamKeys []map[string]string `json:"param_keys"`
	Confirm   *Confirm            `json:"confirm"`

	Path string `json:"path"`
}

func NewEditAction() *Action {
	return &Action{
		Label:   "编辑",
		Icon:    "ep:edit",
		Type:    TableActionTypeForm,
		FormKey: "user-form",
	}
}

func NewDeleteAction() *Action {
	return &Action{
		Label:  "删除",
		Icon:   "ep:delete",
		Type:   TableActionTypeApi,
		ApiKey: "delete",
		ParamKeys: []map[string]string{
			{
				"id": "id",
			},
		},
		Props: map[string]any{"color": "red", "type": "danger"},
		Confirm: &Confirm{
			Title:          "确认删除?",
			Message:        "删除后不可恢复",
			Type:           "warning",
			ConfirmBtnText: "确定",
			CancelBtnText:  "取消",
		},
	}
}

func (a *Action) WithApiKey(key string) *Action {
	a.ApiKey = key
	return a
}

func (a *Action) WithFormKey(key string) *Action {
	a.FormKey = key
	return a
}

func (a *Action) WithConfirm(title, message string) *Action {
	a.Confirm = &Confirm{
		Title:          title,
		Type:           "warning",
		Message:        message,
		ConfirmBtnText: "确定",
		CancelBtnText:  "取消",
	}
	return a
}

func (a *Action) WidthCustomCofirm(c *Confirm) *Action {
	a.Confirm = c
	return a
}

func (a *Action) WithPath(path string) *Action {
	a.Path = path
	return a
}

func (a *Action) WithParamKeys(keys ...map[string]string) *Action {
	a.ParamKeys = keys
	return a
}

func (a *Action) WithProps(props map[string]any) *Action {
	a.Props = props
	return a
}

func (a *Action) WithLabel(label string) *Action {
	a.Label = label
	return a
}

type Form struct {
	Key       string        `json:"key"`
	Title     string        `json:"title"`
	Rows      [][]*FormItem `json:"rows"`
	SubmitApi string        `json:"submit_api"`
}

type FormItem struct {
	Label       string `json:"label"`
	Type        string `json:"type"`
	Prop        string `json:"prop"`
	Placeholder string `json:"placeholder"`
	Width       string `json:"width"`
	Required    bool   `json:"required"`
}

func NewFormItem(label string, prop string, t string, placeholder string) *FormItem {
	return &FormItem{
		Label:       label,
		Type:        t,
		Placeholder: placeholder,
		Prop:        prop,
	}
}

func (f *FormItem) WithWidth(width string) *FormItem {
	f.Width = width
	return f
}

func (f *FormItem) WithRequired(required bool) *FormItem {
	f.Required = required
	return f
}

func NewForm(key, title, api string) *Form {
	return &Form{
		Key:       key,
		Title:     title,
		SubmitApi: api,
	}
}

func (f *Form) WithApi(api string) *Form {
	f.SubmitApi = api
	return f
}

func (f *Form) WithRows(rows [][]*FormItem) *Form {
	f.Rows = rows
	return f
}

type Table struct {
	Widget
	Columns []TableColumn    `json:"columns"`
	Forms   map[string]*Form `json:"forms"`
	Actions []*Action        `json:"actions"`
	Filters *Form            `json:"filters"`
	Search  *Form            `json:"search"`
}

var _ IWidget = (*Table)(nil)

func NewTable() *Table {
	return &Table{
		Widget: Widget{
			Name: "table",
		},
	}
}

func (t *Table) WithForms(forms map[string]*Form) *Table {
	t.Forms = forms
	return t
}

func (t *Table) WithFilters(form *Form) *Table {
	t.Filters = form
	return t
}

func (t *Table) WithSearch(form *Form) *Table {
	t.Search = form
	return t
}

func (t *Table) WithColumns(columns []TableColumn) *Table {
	t.Columns = columns
	return t
}

func (t *Table) WithActions(actions []*Action) *Table {
	t.Actions = actions
	return t
}
