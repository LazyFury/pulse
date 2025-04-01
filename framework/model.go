package framework

import (
	"github.com/lazyfury/pulse/framework/types"
	"gorm.io/gorm"
)

type IModel interface {
	TableName() string
	GetId() uint

	AbleToDelete() bool
	AbleToEdit() bool

	IsDeleted() bool
}

type Model struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt types.LocalTime `json:"created_at" format:"2006-01-02 15:04:05"`
	UpdatedAt types.LocalTime `json:"updated_at" format:"2006-01-02 15:04:05"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

var _ IModel = (*Model)(nil)

func (m *Model) GetId() uint {
	return m.ID
}

func (m *Model) TableName() string {
	return "-"
}

func (m *Model) AbleToDelete() bool {
	return true
}

func (m *Model) AbleToEdit() bool {
	return true
}

func (m *Model) IsDeleted() bool {
	return m.DeletedAt != nil
}
