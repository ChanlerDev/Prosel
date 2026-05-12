package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
)

type SubscriberModel struct {
	ID               string `gorm:"primaryKey;size:36"`
	Email            string `gorm:"size:255;not null;uniqueIndex"`
	Name             string `gorm:"size:100"`
	Status           string `gorm:"size:20;not null;index"`
	VerifyToken      string `gorm:"size:100;not null;uniqueIndex"`
	UnsubscribeToken string `gorm:"size:100;not null;uniqueIndex"`
	VerifiedAt       *time.Time
	UnsubscribedAt   *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (SubscriberModel) TableName() string { return "subscribers" }

type EmailDeliveryModel struct {
	ID           string  `gorm:"primaryKey;size:36"`
	SubscriberID *string `gorm:"size:36"`
	Subject      string  `gorm:"size:255;not null"`
	RefType      string  `gorm:"size:20;index"`
	RefID        string  `gorm:"size:36;index"`
	Status       string  `gorm:"size:20;not null;index"`
	ErrorMessage string
	SentAt       *time.Time
	CreatedAt    time.Time
}

func (EmailDeliveryModel) TableName() string { return "email_deliveries" }

type SubscriberRepository struct{ db *gorm.DB }

func NewSubscriberRepository(db *gorm.DB) *SubscriberRepository { return &SubscriberRepository{db: db} }

func (r *SubscriberRepository) Create(ctx context.Context, subscriber *domain.Subscriber) error {
	return r.db.WithContext(ctx).Create(toSubscriberModel(subscriber)).Error
}

func (r *SubscriberRepository) GetByEmail(ctx context.Context, email string) (*domain.Subscriber, error) {
	var model SubscriberModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	return subscriberFromModel(model, err)
}

func (r *SubscriberRepository) Verify(ctx context.Context, token string) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&SubscriberModel{}).Where("verify_token = ?", token).Updates(map[string]any{"status": string(domain.SubscriberActive), "verified_at": &now, "unsubscribed_at": nil, "updated_at": now})
	return subscriberResultError(result)
}

func (r *SubscriberRepository) Unsubscribe(ctx context.Context, token string) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&SubscriberModel{}).Where("unsubscribe_token = ?", token).Updates(map[string]any{"status": string(domain.SubscriberUnsubscribed), "unsubscribed_at": &now, "updated_at": now})
	return subscriberResultError(result)
}

func (r *SubscriberRepository) List(ctx context.Context, filter domain.SubscriberFilter) ([]domain.Subscriber, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := applySubscriberFilter(r.db.WithContext(ctx).Model(&SubscriberModel{}), filter)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	var models []SubscriberModel
	if err := query.Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	return subscribersFromModels(models), domain.NewPagination(page, perPage, total), nil
}

func (r *SubscriberRepository) ListActive(ctx context.Context) ([]domain.Subscriber, error) {
	var models []SubscriberModel
	if err := r.db.WithContext(ctx).Where("status = ?", string(domain.SubscriberActive)).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return subscribersFromModels(models), nil
}

func (r *SubscriberRepository) CreateDelivery(ctx context.Context, delivery *domain.EmailDelivery) error {
	return r.db.WithContext(ctx).Create(toEmailDeliveryModel(delivery)).Error
}

func (r *SubscriberRepository) UpdateDeliveryStatus(ctx context.Context, id string, status domain.EmailDeliveryStatus, errorMessage string, sentAt *time.Time) error {
	result := r.db.WithContext(ctx).Model(&EmailDeliveryModel{}).Where("id = ?", id).Updates(map[string]any{"status": string(status), "error_message": errorMessage, "sent_at": sentAt})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrSubscriberNotFound
	}
	return nil
}

func applySubscriberFilter(query *gorm.DB, filter domain.SubscriberFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("email ILIKE ? OR name ILIKE ?", search, search)
	}
	return query
}

func subscriberResultError(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrSubscriberNotFound
	}
	return nil
}

func subscriberFromModel(model SubscriberModel, err error) (*domain.Subscriber, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSubscriberNotFound
	}
	if err != nil {
		return nil, err
	}
	return subscriberFromModelNoError(model), nil
}

func subscribersFromModels(models []SubscriberModel) []domain.Subscriber {
	subscribers := make([]domain.Subscriber, 0, len(models))
	for _, model := range models {
		subscribers = append(subscribers, *subscriberFromModelNoError(model))
	}
	return subscribers
}

func subscriberFromModelNoError(model SubscriberModel) *domain.Subscriber {
	return &domain.Subscriber{ID: model.ID, Email: model.Email, Name: model.Name, Status: domain.SubscriberStatus(model.Status), VerifyToken: model.VerifyToken, UnsubscribeToken: model.UnsubscribeToken, VerifiedAt: model.VerifiedAt, UnsubscribedAt: model.UnsubscribedAt, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toSubscriberModel(subscriber *domain.Subscriber) *SubscriberModel {
	return &SubscriberModel{ID: subscriber.ID, Email: subscriber.Email, Name: subscriber.Name, Status: string(subscriber.Status), VerifyToken: subscriber.VerifyToken, UnsubscribeToken: subscriber.UnsubscribeToken, VerifiedAt: subscriber.VerifiedAt, UnsubscribedAt: subscriber.UnsubscribedAt, CreatedAt: subscriber.CreatedAt, UpdatedAt: subscriber.UpdatedAt}
}

func toEmailDeliveryModel(delivery *domain.EmailDelivery) *EmailDeliveryModel {
	return &EmailDeliveryModel{ID: delivery.ID, SubscriberID: delivery.SubscriberID, Subject: delivery.Subject, RefType: delivery.RefType, RefID: delivery.RefID, Status: string(delivery.Status), ErrorMessage: delivery.ErrorMessage, SentAt: delivery.SentAt, CreatedAt: delivery.CreatedAt}
}
