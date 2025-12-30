package repository

import (
	"context"
	"time"

	"github.com/radmickey/money-control/backend/services/insights/models"
	"gorm.io/gorm"
)

// SnapshotRepository handles snapshot operations
type SnapshotRepository struct {
	db *gorm.DB
}

// NewSnapshotRepository creates a new snapshot repository
func NewSnapshotRepository(db *gorm.DB) *SnapshotRepository {
	return &SnapshotRepository{db: db}
}

// Create creates a new snapshot
func (r *SnapshotRepository) Create(ctx context.Context, snapshot *models.Snapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

// GetByDate gets a snapshot for a specific date
func (r *SnapshotRepository) GetByDate(ctx context.Context, userID string, date time.Time) (*models.Snapshot, error) {
	var snapshot models.Snapshot
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&snapshot).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// GetLatest gets the latest snapshot for a user
func (r *SnapshotRepository) GetLatest(ctx context.Context, userID string) (*models.Snapshot, error) {
	var snapshot models.Snapshot
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&snapshot).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// GetByDateRange gets snapshots within a date range
func (r *SnapshotRepository) GetByDateRange(ctx context.Context, userID string, startDate, endDate time.Time, limit int) ([]models.Snapshot, error) {
	var snapshots []models.Snapshot
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	query = query.Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&snapshots).Error; err != nil {
		return nil, err
	}
	return snapshots, nil
}

// GetNetWorthHistory gets net worth history
func (r *SnapshotRepository) GetNetWorthHistory(ctx context.Context, userID string, startDate, endDate time.Time) ([]models.TrendPoint, error) {
	var results []struct {
		Date         time.Time
		TotalNetWorth float64
	}

	query := r.db.WithContext(ctx).
		Model(&models.Snapshot{}).
		Select("date, total_net_worth").
		Where("user_id = ?", userID)

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Order("date ASC").Scan(&results).Error; err != nil {
		return nil, err
	}

	points := make([]models.TrendPoint, len(results))
	for i, r := range results {
		points[i] = models.TrendPoint{
			Date:  r.Date,
			Value: r.TotalNetWorth,
		}
	}
	return points, nil
}

// Upsert creates or updates a snapshot for a date
func (r *SnapshotRepository) Upsert(ctx context.Context, snapshot *models.Snapshot) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", snapshot.UserID, snapshot.Date.Format("2006-01-02")).
		Assign(snapshot).
		FirstOrCreate(snapshot).Error
}

// DeleteOld deletes snapshots older than a date
func (r *SnapshotRepository) DeleteOld(ctx context.Context, userID string, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND date < ?", userID, before).
		Delete(&models.Snapshot{}).Error
}

