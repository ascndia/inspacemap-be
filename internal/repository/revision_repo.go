package repository

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type revisionRepo struct {
	BaseRepository[entity.GraphRevision, uuid.UUID]
	db *gorm.DB
}

func NewGraphRevisionRepository(db *gorm.DB) GraphRevisionRepository {
	return &revisionRepo{
		BaseRepository: NewBaseRepository[entity.GraphRevision, uuid.UUID](db),
		db:             db,
	}
}

func (r *revisionRepo) GetDraftByVenueID(ctx context.Context, venueID uuid.UUID) (*entity.GraphRevision, error) {
	var revision entity.GraphRevision
	// Cari revisi dengan status 'draft' milik venue ini
	// Preload struktur lengkap untuk Editor
	err := r.db.WithContext(ctx).
		Preload("Floors").
		Preload("Floors.MapImage").
		Preload("Floors.Nodes").
		Preload("Floors.Nodes.Panorama").
		Preload("Floors.Nodes.Area").
		Preload("Floors.Nodes.OutgoingEdges").
		Where("venue_id = ? AND status = ?", venueID, entity.StatusDraft).
		First(&revision).Error

	return &revision, err
}

func (r *revisionRepo) GetLiveByVenueID(ctx context.Context, venueID uuid.UUID) (*entity.GraphRevision, error) {
	var venue entity.Venue
	if err := r.db.WithContext(ctx).Select("live_revision_id").First(&venue, "id = ?", venueID).Error; err != nil {
		return nil, err
	}

	if venue.LiveRevisionID == uuid.Nil {
		return nil, errors.New("no live revision found")
	}

	var revision entity.GraphRevision
	err := r.db.WithContext(ctx).
		Preload("Floors"). // Preload secukupnya
		First(&revision, "id = ?", venue.LiveRevisionID).Error

	return &revision, err
}

func (r *revisionRepo) GetDraftByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND status = ?", orgID, entity.StatusDraft).
		Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) GetLiveByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND status = ?", orgID, entity.StatusPublished).
		Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) GetDraftByFloorID(ctx context.Context, floorID uuid.UUID) (*entity.GraphRevision, error) {
	var rev entity.GraphRevision
	// Join untuk memastikan Floor tersebut milik revisi yang statusnya DRAFT
	err := r.db.WithContext(ctx).
		Joins("JOIN floors ON floors.graph_revision_id = graph_revisions.id").
		Where("floors.id = ? AND graph_revisions.status = ?", floorID, entity.StatusDraft).
		First(&rev).Error

	if err != nil {
		return nil, err // Error jika tidak ketemu (artinya Floor itu milik Live atau tidak ada)
	}
	return &rev, nil
}

func (r *revisionRepo) GetLiveByFloorID(ctx context.Context, floorID uuid.UUID) (*entity.GraphRevision, error) {
	var floor entity.Floor
	if err := r.db.WithContext(ctx).Select("graph_revision_id").First(&floor, "id = ?", floorID).Error; err != nil {
		return nil, err
	}
	var rev entity.GraphRevision
	if err := r.db.WithContext(ctx).First(&rev, "id = ? AND status = ?", floor.GraphRevisionID, entity.StatusPublished).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *revisionRepo) GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	err := r.db.WithContext(ctx).Where("venue_id = ?", venueID).Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) FilterGraphRevisions(ctx context.Context, filter models.FilterGraphRevision) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) PagedGraphRevisions(ctx context.Context, q models.QueryGraphRevision) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision

	db := r.buildFilterQuery(ctx, q.FilterGraphRevision)

	if q.Sort != nil {
		db = db.Order(*q.Sort)
	} else {
		db = db.Order("created_at desc")
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	offset := 0
	if q.Offset != nil && *q.Offset >= 0 {
		offset = *q.Offset
	}

	err := db.Limit(limit).Offset(offset).Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) CursorGraphRevisions(ctx context.Context, q models.CursorGraphRevisionQuery) ([]entity.GraphRevision, error) {
	var revs []entity.GraphRevision
	db := r.buildFilterQuery(ctx, q.FilterGraphRevision)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorRev entity.GraphRevision
			if err := r.db.Select("created_at").First(&cursorRev, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorRev.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").Limit(limit).Find(&revs).Error
	return revs, err
}

func (r *revisionRepo) CreateDraft(ctx context.Context, venueID uuid.UUID, userID uuid.UUID) (*entity.GraphRevision, error) {
	// 1. Ambil organization_id dari venue
	var venue entity.Venue
	if err := r.db.WithContext(ctx).Select("organization_id").First(&venue, "id = ?", venueID).Error; err != nil {
		return nil, err
	}

	// 2. Siapkan Draft Baru
	newDraft := entity.GraphRevision{
		OrganizationID: venue.OrganizationID,
		CreatedByID:    userID,
		VenueID:        venueID,
		Status:         entity.StatusDraft,
		Note:           "Auto-generated draft",
	}

	// Transaction: Simpan Revisi
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newDraft).Error; err != nil {
			return err
		}
		return nil
	})

	return &newDraft, err
}

// func (r *revisionRepo) PublishDraft(ctx context.Context, draftID uuid.UUID, note string) error {
// 	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		var draft entity.GraphRevision
// 		if err := tx.Preload("Floors.Nodes.OutgoingEdges").First(&draft, draftID).Error; err != nil {
// 			return err
// 		}

// 		newRevision := entity.GraphRevision{
// 			VenueID:        draft.VenueID,
// 			OrganizationID: draft.OrganizationID,
// 			Status:         entity.StatusPublished,
// 			Note:           note,
// 		}
// 		if err := tx.Create(&newRevision).Error; err != nil {
// 			return err
// 		}

// 		// 3. CLONING PROCESS
// 		nodeIDMap := make(map[uuid.UUID]uuid.UUID)
// 		var newStartNodeID *uuid.UUID

// 		for _, floor := range draft.Floors {
// 			newFloor := entity.Floor{
// 				GraphRevisionID: newRevision.ID,
// 				VenueID:         draft.VenueID,
// 				Name:            floor.Name,
// 				LevelIndex:      floor.LevelIndex,
// 				MapImageID:      floor.MapImageID,
// 				PixelsPerMeter:  floor.PixelsPerMeter,
// 				IsActive:        floor.IsActive,
// 			}
// 			if err := tx.Create(&newFloor).Error; err != nil {
// 				return err
// 			}

// 			// Clone Nodes
// 			for _, node := range floor.Nodes {
// 				newNode := entity.GraphNode{
// 					FloorID:         newFloor.ID,
// 					AreaID:          node.AreaID,
// 					X:               node.X,
// 					Y:               node.Y,
// 					PanoramaAssetID: node.PanoramaAssetID,
// 					RotationOffset:  node.RotationOffset,
// 					Label:           node.Label,
// 					Properties:      node.Properties,
// 					IsActive:        node.IsActive,
// 				}
// 				if err := tx.Create(&newNode).Error; err != nil {
// 					return err
// 				}

// 				nodeIDMap[node.ID] = newNode.ID

// 				// Cek Start Node
// 				if draft.StartNodeID != nil && node.ID == *draft.StartNodeID {
// 					id := newNode.ID
// 					newStartNodeID = &id
// 				}
// 			}
// 		}

// 		// 4. Reconstruct Edges
// 		for _, floor := range draft.Floors {
// 			for _, node := range floor.Nodes {
// 				for _, edge := range node.OutgoingEdges {
// 					newFrom, ok1 := nodeIDMap[edge.FromNodeID]
// 					newTo, ok2 := nodeIDMap[edge.ToNodeID]
// 					if ok1 && ok2 {
// 						newEdge := entity.GraphEdge{
// 							FromNodeID: newFrom,
// 							ToNodeID:   newTo,
// 							Heading:    edge.Heading,
// 							Distance:   edge.Distance,
// 							Type:       edge.Type,
// 							IsActive:   edge.IsActive,
// 						}
// 						if err := tx.Create(&newEdge).Error; err != nil {
// 							return err
// 						}
// 					}
// 				}
// 			}
// 		}

// 		// 5. Update Start Node di Revisi Baru
// 		if newStartNodeID != nil {
// 			if err := tx.Model(&newRevision).Update("start_node_id", *newStartNodeID).Error; err != nil {
// 				return err
// 			}
// 		}

// 		// 6. Update Venue Pointer Live
// 		if err := tx.Model(&entity.Venue{BaseEntity: entity.BaseEntity{ID: draft.VenueID}}).
// 			Update("live_revision_id", newRevision.ID).Error; err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// PublishDraft: Deep Copy Logic
func (r *revisionRepo) PublishDraft(ctx context.Context, revisionID uuid.UUID, note string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Ambil Draft
		var draft entity.GraphRevision
		if err := tx.Preload("Floors").
			Preload("Floors.MapImage").
			Preload("Floors.Nodes", func(db *gorm.DB) *gorm.DB { return db.Preload("Panorama") }).
			Preload("Floors.Nodes.OutgoingEdges", func(db *gorm.DB) *gorm.DB { return db.Where("is_active = ?", true) }).
			Preload("Floors.Areas", func(db *gorm.DB) *gorm.DB { return db.Preload("Gallery").Preload("CoverImage") }).
			First(&draft, "id = ?", revisionID).Error; err != nil {
			return err
		}

		if draft.Status != entity.StatusDraft {
			return errors.New("only drafts can be published")
		}

		// 2. Buat Revisi Baru (LIVE)
		newLiveRev := entity.GraphRevision{
			VenueID:        draft.VenueID,
			OrganizationID: draft.OrganizationID,
			CreatedByID:    draft.CreatedByID, // Copy dari draft
			Status:         entity.StatusPublished,
			Note:           note,
		}
		if err := tx.Create(&newLiveRev).Error; err != nil {
			return err
		}

		// 3. CLONING (Mapping UUID Lama -> UUID Baru)
		nodeIDMap := make(map[uuid.UUID]uuid.UUID)
		areaIDMap := make(map[uuid.UUID]uuid.UUID)
		var newStartNodeID *uuid.UUID

		for _, floor := range draft.Floors {
			// Clone Floor
			newFloor := entity.Floor{
				GraphRevisionID: newLiveRev.ID,
				VenueID:         draft.VenueID,
				Name:            floor.Name,
				LevelIndex:      floor.LevelIndex,
				MapImageID:      floor.MapImageID,
				PixelsPerMeter:  floor.PixelsPerMeter,
				IsActive:        floor.IsActive,
				MapWidth:        floor.MapWidth,
				MapHeight:       floor.MapHeight,
			}
			if err := tx.Create(&newFloor).Error; err != nil {
				return err
			}

			// Clone Areas first and keep mapping
			for _, area := range floor.Areas {
				newArea := entity.Area{
					GraphRevisionID: newLiveRev.ID,
					FloorID:         newFloor.ID,
					VenueID:         draft.VenueID,
					Name:            area.Name,
					Slug:            area.Slug,
					Label:           area.Label,
					Description:     area.Description,
					Latitude:        area.Latitude,
					Longitude:       area.Longitude,
					Boundary:        area.Boundary,
					LabelX:          area.LabelX,
					LabelY:          area.LabelY,
					Category:        area.Category,
					CoverImageID:    area.CoverImageID,
					StartNodeID:     area.StartNodeID,
				}
				if err := tx.Create(&newArea).Error; err != nil {
					return err
				}

				areaIDMap[area.ID] = newArea.ID

				// Clone area gallery items
				for _, g := range area.Gallery {
					newG := entity.AreaGalleryItem{
						AreaID:       newArea.ID,
						MediaAssetID: g.MediaAssetID,
						SortOrder:    g.SortOrder,
						Caption:      g.Caption,
						IsVisible:    g.IsVisible,
					}
					if err := tx.Create(&newG).Error; err != nil {
						return err
					}
				}
			}

			// Clone Nodes (map AreaIDs to new ones)
			for _, node := range floor.Nodes {
				// Check if panorama_asset_id exists
				if err := tx.Model(&entity.MediaAsset{}).Where("id = ?", node.PanoramaAssetID).Select("1").First(&struct{}{}).Error; err != nil {
					// Panorama asset not found, skip this node
					continue
				}

				newNode := entity.GraphNode{
					FloorID:         newFloor.ID,
					AreaID:          node.AreaID,
					X:               node.X,
					Y:               node.Y,
					PanoramaAssetID: node.PanoramaAssetID,
					RotationOffset:  node.RotationOffset,
					Label:           node.Label,
					Properties:      node.Properties,
					IsActive:        node.IsActive,
				}
				// Remap AreaID if needed
				if node.AreaID != nil {
					if mapped, ok := areaIDMap[*node.AreaID]; ok {
						newNode.AreaID = &mapped
					} else {
						// area not present in this floor in draft (should be nil)
						newNode.AreaID = nil
					}
				}

				if err := tx.Create(&newNode).Error; err != nil {
					return err
				}

				nodeIDMap[node.ID] = newNode.ID

				// Cek Start Node
				if draft.StartNodeID != nil && node.ID == *draft.StartNodeID {
					id := newNode.ID
					newStartNodeID = &id
				}
			}
		}

		// 4. Reconstruct Edges
		for _, floor := range draft.Floors {
			for _, node := range floor.Nodes {
				for _, edge := range node.OutgoingEdges {
					newFrom, ok1 := nodeIDMap[edge.FromNodeID]
					newTo, ok2 := nodeIDMap[edge.ToNodeID]
					if ok1 && ok2 {
						newEdge := entity.GraphEdge{
							FromNodeID: newFrom,
							ToNodeID:   newTo,
							Heading:    edge.Heading,
							Distance:   edge.Distance,
							Type:       edge.Type,
							IsActive:   edge.IsActive,
						}
						if err := tx.Create(&newEdge).Error; err != nil {
							return err
						}
					}
				}
			}
		}

		// 4.5. Map Area StartNodeID to new node IDs
		for _, floor := range draft.Floors {
			for _, area := range floor.Areas {
				if area.StartNodeID != nil {
					if newNodeID, ok := nodeIDMap[*area.StartNodeID]; ok {
						// Update area StartNodeID to new node ID
						if err := tx.Model(&entity.Area{}).Where("id = ?", areaIDMap[area.ID]).Update("start_node_id", newNodeID).Error; err != nil {
							return err
						}
					}
				}
			}
		}

		// 5. Finalisasi Pointer
		if newStartNodeID != nil {
			tx.Model(&newLiveRev).Update("start_node_id", *newStartNodeID)
		}

		// Ambil live_revision_id lama sebelum update
		var venue entity.Venue
		if err := tx.Select("live_revision_id").Where("id = ?", draft.VenueID).First(&venue).Error; err != nil {
			return err
		}

		oldLiveID := venue.LiveRevisionID

		// Update Venue agar menunjuk ke Live Revision yang baru
		if err := tx.Model(&entity.Venue{BaseEntity: entity.BaseEntity{ID: draft.VenueID}}).Update("live_revision_id", newLiveRev.ID).Error; err != nil {
			return err
		}

		// Archive live revision lama jika ada dan berbeda dari yang baru
		if oldLiveID != uuid.Nil && oldLiveID != newLiveRev.ID {
			if err := tx.Model(&entity.GraphRevision{}).Where("id = ?", oldLiveID).Update("status", entity.StatusArchived).Error; err != nil {
				return err
			}
		}

		// Set live revision as immutable
		if err := tx.Model(&newLiveRev).Update("is_immutable", true).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *revisionRepo) UpdateRevision(ctx context.Context, revisionID uuid.UUID, note string) error {
	return r.db.WithContext(ctx).Model(&entity.GraphRevision{}).
		Where("id = ?", revisionID).
		Update("note", note).Error
}

// --- QUERY BUILDER ---
func (r *revisionRepo) buildFilterQuery(ctx context.Context, f models.FilterGraphRevision) *gorm.DB {
	db := r.db.WithContext(ctx)

	if f.OrganizationID != nil {
		db = db.Where("organization_id = ?", *f.OrganizationID)
	}
	if f.VenueID != nil {
		db = db.Where("venue_id = ?", *f.VenueID)
	}
	if f.CreatedByID != nil {
		db = db.Where("created_by_id = ?", *f.CreatedByID)
	}
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}
	if f.Note != nil {
		db = db.Where("note ILIKE ?", "%"+*f.Note+"%")
	}

	// Time Filters
	if f.CreatedAfter != nil {
		if t, err := time.Parse(time.RFC3339, *f.CreatedAfter); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if f.CreatedBefore != nil {
		if t, err := time.Parse(time.RFC3339, *f.CreatedBefore); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}

	// Filter by Floor (Join required)
	if f.FloorID != nil {
		db = db.Joins("JOIN floors ON floors.graph_revision_id = graph_revisions.id").
			Where("floors.id = ?", *f.FloorID)
	}

	return db
}
