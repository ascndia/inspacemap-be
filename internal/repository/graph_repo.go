package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
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