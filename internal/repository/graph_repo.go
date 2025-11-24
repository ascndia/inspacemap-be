package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type graphRepo struct {
	db *gorm.DB
}

func NewGraphRepository(db *gorm.DB) GraphRepository {
	return &graphRepo{db: db}
}

func (r *graphRepo) CreateNode(ctx context.Context, node *entity.GraphNode) error {
	return r.db.WithContext(ctx).Create(node).Error
}

func (r *graphRepo) UpdateNodePosition(ctx context.Context, id uuid.UUID, x, y float64) error {
	return r.db.WithContext(ctx).Model(&entity.GraphNode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"x": x, "y": y}).Error
}

func (r *graphRepo) UpdateNodeCalibration(ctx context.Context, id uuid.UUID, offset float64) error {
	return r.db.WithContext(ctx).Model(&entity.GraphNode{}).
		Where("id = ?", id).
		Update("rotation_offset", offset).Error
}

func (r *graphRepo) UpdateNode(ctx context.Context, id uuid.UUID, req models.UpdateNodeRequest) error {
	updates := make(map[string]interface{})
	
	if req.X != nil {
		updates["x"] = *req.X
	}
	if req.Y != nil {
		updates["y"] = *req.Y
	}
	if req.PanoramaAssetID != nil {
		updates["panorama_asset_id"] = *req.PanoramaAssetID
	}
	if req.Label != nil {
		updates["label"] = *req.Label
	}
	if req.RotationOffset != nil {
		updates["rotation_offset"] = *req.RotationOffset
	}
	
	if len(updates) == 0 {
		return nil // No updates to make
	}
	
	return r.db.WithContext(ctx).Model(&entity.GraphNode{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *graphRepo) DeleteNode(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.GraphNode{}, "id = ?", id).Error
}

func (r *graphRepo) ConnectNodes(ctx context.Context, edge *entity.GraphEdge) error {
	var nodeA, nodeB entity.GraphNode
	if err := r.db.WithContext(ctx).Select("id, x, y, floor_id").First(&nodeA, "id = ?", edge.FromNodeID).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Select("id, x, y, floor_id").First(&nodeB, "id = ?", edge.ToNodeID).Error; err != nil {
		return err
	}

	dx := nodeB.X - nodeA.X
	dy := nodeB.Y - nodeA.Y 

	dist := math.Sqrt(dx*dx + dy*dy)

	headingRad := math.Atan2(dx, -dy)
	headingDeg := headingRad * (180 / math.Pi)
	if headingDeg < 0 {
		headingDeg += 360
	}

	edge.Distance = dist
	edge.Heading = headingDeg
	
	return r.db.WithContext(ctx).Create(edge).Error
}

func (r *graphRepo) DeleteEdge(ctx context.Context, fromID, toID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("from_node_id = ? AND to_node_id = ?", fromID, toID).
		Delete(&entity.GraphEdge{}).Error
}

func (r *graphRepo) CountNodesByFloorID(ctx context.Context, floorID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.GraphNode{}).
		Where("floor_id = ?", floorID).
		Count(&count).Error
	return int(count), err
}

func (r *graphRepo) CountAreasByFloorID(ctx context.Context, floorID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.GraphNode{}).
		Where("floor_id = ? AND area_id IS NOT NULL", floorID).
		Count(&count).Error
	return int(count), err
}

func (r *graphRepo) GetNodesByFloorID(ctx context.Context, floorID uuid.UUID) ([]entity.GraphNode, error) {
	var nodes []entity.GraphNode
	err := r.db.WithContext(ctx).Preload("Panorama").Preload("Area").Where("floor_id = ?", floorID).Find(&nodes).Error
	return nodes, err
}

func (r *graphRepo) GetEdgesFromNode(ctx context.Context, nodeID uuid.UUID) ([]entity.GraphEdge, error) {
	var edges []entity.GraphEdge
	err := r.db.WithContext(ctx).Preload("ToNode").Where("from_node_id = ?", nodeID).Find(&edges).Error
	return edges, err
}

func (r *graphRepo) CreateEdge(ctx context.Context, edge *entity.GraphEdge) error {
	return r.db.WithContext(ctx).Create(edge).Error
}