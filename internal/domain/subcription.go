package domain

import "github.com/google/uuid"

type Subcription struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServiceName string    `gorm:"not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	StartDate   MonthYear `gorm:"not null"`
}
