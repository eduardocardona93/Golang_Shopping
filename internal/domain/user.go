package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a customer (Usuario) able to make purchases.
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Nombre         string         `gorm:"type:varchar(100);not null" json:"nombre"`
	Apellido       string         `gorm:"type:varchar(100);not null" json:"apellido"`
	Telefono       string         `gorm:"type:varchar(20)" json:"telefono"`
	Direccion      string         `gorm:"type:varchar(255)" json:"direccion"`
	Identificacion string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"identificacion"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
