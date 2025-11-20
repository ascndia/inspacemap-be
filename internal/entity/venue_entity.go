package entity

import (
	"github.com/google/uuid"
)


type Venue struct {
	BaseEntity
	OrganizationID uuid.UUID   `gorm:"index;not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	Name           string `gorm:"type:varchar(100);not null"`
	Slug           string `gorm:"type:varchar(100);index"`
	Description    string `gorm:"type:text"`
	Address        string `gorm:"type:text"`    
	City           string `gorm:"type:varchar(100)"`
	Province       string `gorm:"type:varchar(100)"`
	PostalCode     string `gorm:"type:varchar(20)"`
	Visibility     VisibilityStatus `gorm:"type:varchar(20);default:'private'"`
	Latitude       float64 `gorm:"type:decimal(10,8);index"` 
	Longitude      float64 `gorm:"type:decimal(11,8);index"` 
	LiveRevisionID  uuid.UUID  `gorm:"index;not null"`
	LiveRevision    *GraphRevision `gorm:"foreignKey:LiveRevisionID"`
	DraftRevisionID *uuid.UUID `gorm:"index"`
	DraftRevision   *GraphRevision `gorm:"foreignKey:DraftRevisionID"`
	CoverImageID    *uuid.UUID  
	CoverImage      *MediaAsset `gorm:"foreignKey:CoverImageID"`
	Gallery         []VenueGalleryItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Revisions []GraphRevision `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}


type VenueGalleryItem struct {
	BaseEntity
	VenueID      uuid.UUID  `gorm:"type:uuid;index:idx_venue_media,unique;not null"`
	MediaAssetID uuid.UUID  `gorm:"type:uuid;index:idx_venue_media,unique;not null"`
	SortOrder    int        `gorm:"default:0"` 
	Caption      string     `gorm:"type:varchar(255)"` 
	IsVisible    bool       `gorm:"default:true"`
	IsFeatured   bool       `gorm:"default:false"`     
	Venue        Venue      `gorm:"foreignKey:VenueID"`
	MediaAsset   MediaAsset `gorm:"foreignKey:MediaAssetID"`
}

