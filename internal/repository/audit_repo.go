package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type auditRepo struct {
	BaseRepository[entity.AuditLog, uint]
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditLogRepository {
	return &auditRepo{
		BaseRepository: NewBaseRepository[entity.AuditLog, uint](db),
		db:             db,
	}
}

func (r *auditRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("organization_id = ?", orgID).
		Order("id desc"). 
		Limit(100).       
		Find(&logs).Error
	return logs, err
}

func (r *auditRepo) FilterAuditLogs(ctx context.Context, filter models.AuditLogFilter) ([]entity.AuditLog, error) {
	var logs []entity.AuditLog
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditRepo) PagedAuditLogs(ctx context.Context, q models.AuditLogQuery) ([]entity.AuditLog, int64, error) {
	var logs []entity.AuditLog
	var total int64

	db := r.buildFilterQuery(ctx, q.AuditLogFilter)

	if err := db.Model(&entity.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.Sort != nil {
		db = db.Order(*q.Sort)
	} else {
		db = db.Order("id desc") 
	}

	limit := 20
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	offset := 0
	if q.Offset != nil && *q.Offset >= 0 {
		offset = *q.Offset
	}

	err := db.Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *auditRepo) CursorAuditLogs(ctx context.Context, q models.AuditLogQueryCursor) ([]entity.AuditLog, string, error) {
	var logs []entity.AuditLog
	db := r.buildFilterQuery(ctx, q.AuditLogFilter)
	
	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := strconv.ParseUint(*q.Cursor, 10, 64); err == nil {
			
			db = db.Where("id < ?", cursorID)
		}
	}

	limit := 20
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	
	err := db.Order("id desc").
		Limit(limit + 1). 
		Find(&logs).Error

	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(logs) > limit {
		logs = logs[:limit]
		
		nextCursor = strconv.FormatUint(uint64(logs[len(logs)-1].ID), 10)
	}

	return logs, nextCursor, nil
}

func (r *auditRepo) buildFilterQuery(ctx context.Context, f models.AuditLogFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Preload("User")

	if f.OrganizationID != "" {
		db = db.Where("organization_id = ?", f.OrganizationID)
	}
	if f.Entity != "" {
		db = db.Where("entity = ?", f.Entity)
	}
	if f.EntityID != "" {
		db = db.Where("entity_id = ?", f.EntityID)
	}
	if f.Action != "" {
		db = db.Where("action = ?", f.Action)
	}
	if f.UserID != "" {
		db = db.Where("user_id = ?", f.UserID)
	}
	if f.IPAddress != "" {
		db = db.Where("ip_address = ?", f.IPAddress)
	}
	
	if f.FromDate != "" {
		if t, err := time.Parse("2006-01-02", f.FromDate); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if f.ToDate != "" {
		if t, err := time.Parse("2006-01-02", f.ToDate); err == nil {
			
			t = t.Add(24 * time.Hour)
			db = db.Where("created_at < ?", t)
		}
	}

	return db
}