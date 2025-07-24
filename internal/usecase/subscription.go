package usecase

import (
	"context"

	"github.com/effective-mobile/subscriptions/internal/domain"
	"github.com/google/uuid"
)

type subscriptionUseCase struct {
	repo domain.SubscriptionRepository
}

func NewSubscriptionUseCase(repo domain.SubscriptionRepository) SubscriptionUseCase {
	return &subscriptionUseCase{repo: repo}
}

func (u *subscriptionUseCase) Create(ctx context.Context, sub *domain.Subscription) error {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}
	return u.repo.Create(sub)
}

func (u *subscriptionUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return u.repo.GetByID(id)
}

func (u *subscriptionUseCase) Update(ctx context.Context, sub *domain.Subscription) error {
	return u.repo.Update(sub)
}

func (u *subscriptionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(id)
}

func (u *subscriptionUseCase) List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	return u.repo.List(filter)
}

func (u *subscriptionUseCase) Sum(ctx context.Context, filter domain.SubscriptionSumFilter) (int, error) {
	return u.repo.Sum(filter)
}

// SubscriptionUseCase описывает бизнес-операции над подписками
// (application layer)
type SubscriptionUseCase interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter domain.SubscriptionFilter) ([]domain.Subscription, error)
	Sum(ctx context.Context, filter domain.SubscriptionSumFilter) (int, error)
}
