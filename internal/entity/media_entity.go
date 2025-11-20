package entity

import (
	"time"

	"github.com/google/uuid"
)

type VisibilityStatus string

const (
	VisibilityPublic   VisibilityStatus = "public"   // Muncul di pencarian global
	VisibilityUnlisted VisibilityStatus = "unlisted" // Hanya bisa diakses via Link/QR
	VisibilityPrivate  VisibilityStatus = "private"  // Hanya admin/internal
	VisibilityArchived VisibilityStatus = "archived" // Disembunyikan (Soft Archive)
)

type MediaAsset struct {
	BaseEntity
	OrganizationID uuid.UUID     `gorm:"index;not null"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID"`

	StorageProvider string `gorm:"type:varchar(20);default:'s3'"`
	Bucket          string `gorm:"type:varchar(100);not null"`
	Key             string `gorm:"type:varchar(255);not null;index"`
	Region          string `gorm:"type:varchar(50)"`

	PublicURL    string `gorm:"type:text;not null"`
	ThumbnailURL string `gorm:"type:text"`

	FileName    string
	MimeType    string `gorm:"type:varchar(50)"`
	Type        string `gorm:"index"`
	SizeInBytes int64  `gorm:"index"`

	Width    int    `gorm:"default:0"`
	Height   int    `gorm:"default:0"`
	BlurHash string `gorm:"type:varchar(100)"`
	AltText  string `gorm:"type:varchar(255)"`

	Visibility VisibilityStatus `gorm:"type:varchar(20);default:'public'"`

	UploadedBy uuid.UUID `gorm:"type:uuid;index"`
	Uploader   *User     `gorm:"foreignKey:UploadedBy"`
	UploadedAt time.Time `gorm:"autoCreateTime"`

	IsEncrypted  bool   `gorm:"default:false"`
	Checksum     string `gorm:"type:varchar(64)"`
	LastAccessed time.Time
	UsageCount   int `gorm:"default:0"`
}
