package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type dayClosure struct {
	UserId   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Date     time.Time `gorm:"type:date;primaryKey"`
	ClosedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (dayClosure) TableName() string { return "day_closures" }

type dayClosureRepository struct {
	db *gorm.DB
}

var _ ports.DayClosureRepository = (*dayClosureRepository)(nil)

func NewDayClosureRepository(db *gorm.DB) *dayClosureRepository {
	return &dayClosureRepository{db: db}
}

func (r *dayClosureRepository) IsClosed(userId uuid.UUID, date time.Time) (bool, error) {
	var count int64
	err := r.db.
		Model(&dayClosure{}).
		Where("user_id = ? AND date = ?", userId, date).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking day closure", err)
	}
	return count > 0, nil
}

func (r *dayClosureRepository) Close(userId uuid.UUID, date time.Time) error {
	record := &dayClosure{
		UserId: userId,
		Date:   date,
	}
	err := r.db.
		Where("user_id = ? AND date = ?", userId, date).
		FirstOrCreate(record).
		Error
	if err != nil {
		return cerr.NewInternalError("closing day", err)
	}
	return nil
}

func (r *dayClosureRepository) Open(userId uuid.UUID, date time.Time) error {
	err := r.db.
		Where("user_id = ? AND date = ?", userId, date).
		Delete(&dayClosure{}).
		Error
	if err != nil {
		return cerr.NewInternalError("opening day", err)
	}
	return nil
}

func (r *dayClosureRepository) GetClosedDates(userId uuid.UUID, from, to time.Time) (map[string]bool, error) {
	var results []struct {
		Date string
	}
	err := r.db.
		Table("day_closures").
		Select("date::text as date").
		Where("user_id = ? AND date >= ? AND date <= ?", userId, from, to).
		Scan(&results).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching closed dates", err)
	}
	m := make(map[string]bool, len(results))
	for _, r := range results {
		m[r.Date] = true
	}
	return m, nil
}
