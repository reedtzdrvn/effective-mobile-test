package domain

import (
	"time"
	"github.com/google/uuid"
)

// Subscription — бизнес-сущность подписки
// EndDate может быть nil (опционально)
type Subscription struct {
	ID           uuid.UUID  // внутренний идентификатор (можно использовать для PK)
	ServiceName  string     // название сервиса
	Price        int        // стоимость в рублях
	UserID       uuid.UUID  // id пользователя
	StartDate    time.Time  // месяц и год начала (день = 1)
	EndDate      *time.Time // месяц и год окончания (опционально)
	CreatedAt    time.Time  // дата создания
	UpdatedAt    time.Time  // дата обновления
}

// SubscriptionRepository — интерфейс для работы с подписками
// Реализация будет в internal/repository
//
type SubscriptionRepository interface {
	Create(sub *Subscription) error
	GetByID(id uuid.UUID) (*Subscription, error)
	Update(sub *Subscription) error
	Delete(id uuid.UUID) error
	List(filter SubscriptionFilter) ([]Subscription, error)
	Sum(filter SubscriptionSumFilter) (int, error)
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time // период: с
	To          *time.Time // период: по
}

type SubscriptionSumFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        time.Time // период: с (обязателен)
	To          time.Time // период: по (обязателен)
}



