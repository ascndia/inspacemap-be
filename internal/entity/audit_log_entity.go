package entity

import (
	"time"

	"github.com/google/uuid"
)
type AuditLog struct {
	ID        uint      `gorm:"primaryKey"` 
	CreatedAt time.Time `gorm:"index"`
	
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID         uuid.UUID `gorm:"type:uuid;index"`
	User           *User     `gorm:"foreignKey:UserID"`
	
	Action   string `gorm:"type:varchar(50);index"`
	Entity   string `gorm:"type:varchar(50);index"`
	EntityID string `gorm:"type:varchar(50)"`     
	
	Details JSONMap `gorm:"type:jsonb"`
	
	IPAddress string `gorm:"type:varchar(50)"`
	UserAgent string `gorm:"type:varchar(255)"`
}