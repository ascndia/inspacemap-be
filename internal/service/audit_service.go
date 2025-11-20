package service

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"log"

	"github.com/google/uuid"
)

type auditService struct {
	auditRepo repository.AuditLogRepository
}

func NewAuditService(repo repository.AuditLogRepository) AuditService {
	return &auditService{
		auditRepo: repo,
	}
}

// 1. GetActivityLogs (Read)
func (s *auditService) GetActivityLogs(ctx context.Context, orgID uuid.UUID, req models.AuditLogQueryCursor) (*models.AuditListResponse, error) {
	// Panggil Repository dengan Filter
	logs, nextCursor, err := s.auditRepo.CursorAuditLogs(ctx, req)
	if err != nil {
		return nil, err
	}

	// Mapping Entity -> DTO
	var data []models.AuditLogResponse
	for _, item := range logs {
		actorName := "System"
		actorEmail := ""

		// Jika user ada (terhapus pun tetap ke-load kalau soft delete,
		// tapi kalau hard delete mungkin nil, jadi perlu cek)
		if item.User != nil {
			actorName = item.User.FullName
			actorEmail = item.User.Email
		}

		data = append(data, models.AuditLogResponse{
			ID:             item.ID,
			CreatedAt:      item.CreatedAt,
			OrganizationID: item.OrganizationID,
			UserID:         &item.UserID,
			ActorName:      actorName,
			ActorEmail:     actorEmail,
			Action:         item.Action,
			Entity:         item.Entity,
			EntityID:       item.EntityID,
			Details:        item.Details, // JSONMap otomatis jadi map[string]interface{}
			IPAddress:      item.IPAddress,
		})
	}

	return &models.AuditListResponse{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
	}, nil
}

// 2. LogActivity (Write - Async)
func (s *auditService) LogActivity(ctx context.Context, req models.CreateAuditLogRequest) {
	// Jalankan di Goroutine agar Non-Blocking
	// Kita copy context atau buat context background baru agar tidak cancel saat HTTP request selesai
	go func(payload models.CreateAuditLogRequest) {
		// Buat Entity
		logEntry := entity.AuditLog{
			OrganizationID: payload.OrganizationID,
			UserID:         payload.UserID,
			Action:         payload.Action,
			Entity:         payload.Entity,
			EntityID:       payload.EntityID,
			Details:        payload.Details,
			IPAddress:      payload.IPAddress,
			UserAgent:      payload.UserAgent,
			// CreatedAt otomatis
		}

		// Simpan (Gunakan context.Background karena request asli mungkin sudah selesai)
		if err := s.auditRepo.Create(context.Background(), &logEntry); err != nil {
			// Jika gagal log, kita cuma bisa print error ke console server
			// Tidak boleh membatalkan transaksi bisnis utama
			log.Printf("[AUDIT ERROR] Failed to write log: %v", err)
		}
	}(req)
}
