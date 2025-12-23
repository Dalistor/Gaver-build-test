package patterns

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DTO interface {
	Validate() (err error)
	ToModel() *Model
}

type Model interface {
	TableName() string
	Validate() (err error)
	ToDTO() *DTO
}

type DefaultModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
