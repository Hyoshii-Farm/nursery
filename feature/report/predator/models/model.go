package predator

import (
	"time"

	"gorm.io/gorm"
)

func (Predator) TableName() string {
	return "Predator"
}

type Predator struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(25);not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PredatorHistory) TableName() string {
	return "PredatorHistory"
}

type PredatorHistory struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PredatorID uint           `json:"predator_id"`
	Datetime   time.Time      `json:"datetime"`
	Quantity   float64        `json:"quantity"`
	Action     string         `json:"action"`
	IsActive   bool           `json:"is_active"`
	CompanyID  uint           `json:"company_id"`
	Pic        string         `json:"pic"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Predator Predator `gorm:"foreignKey:PredatorID" json:"predator"`
}
