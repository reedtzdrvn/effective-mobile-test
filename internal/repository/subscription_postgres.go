package repository

import (
	"context"
	"time"

	"github.com/effective-mobile/subscriptions/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionPostgres struct {
	pool *pgxpool.Pool
}

func NewSubscriptionPostgres(pool *pgxpool.Pool) *SubscriptionPostgres {
	return &SubscriptionPostgres{pool: pool}
}

func (r *SubscriptionPostgres) Create(sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())
		RETURNING created_at, updated_at`
	err := r.pool.QueryRow(context.Background(), query,
		sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&sub.CreatedAt, &sub.UpdatedAt)
	return err
}

func (r *SubscriptionPostgres) GetByID(id uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions WHERE id = $1`
	row := r.pool.QueryRow(context.Background(), query, id)
	var sub domain.Subscription
	var endDate *time.Time
	if err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
		return nil, err
	}
	sub.EndDate = endDate
	return &sub, nil
}

func (r *SubscriptionPostgres) Update(sub *domain.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5, updated_at=now() 
		WHERE id=$6
		RETURNING created_at`

	err := r.pool.QueryRow(context.Background(), query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID,
	).Scan(&sub.CreatedAt)

	return err
}

func (r *SubscriptionPostgres) Delete(id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id=$1`
	_, err := r.pool.Exec(context.Background(), query, id)
	return err
}

func (r *SubscriptionPostgres) List(filter domain.SubscriptionFilter) ([]domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if filter.UserID != nil {
		query += " AND user_id = $" + itoa(idx)
		args = append(args, *filter.UserID)
		idx++
	}
	if filter.ServiceName != nil {
		query += " AND service_name = $" + itoa(idx)
		args = append(args, *filter.ServiceName)
		idx++
	}
	if filter.From != nil {
		query += " AND start_date >= $" + itoa(idx)
		args = append(args, *filter.From)
		idx++
	}
	if filter.To != nil {
		query += " AND (end_date <= $" + itoa(idx) + " OR end_date IS NULL)"
		args = append(args, *filter.To)
		idx++
	}
	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var endDate *time.Time
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, err
		}
		sub.EndDate = endDate
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *SubscriptionPostgres) Sum(filter domain.SubscriptionSumFilter) (int, error) {
	// query := `SELECT COALESCE(SUM(price),0) FROM subscriptions WHERE start_date >= $1 AND (end_date <= $2 OR end_date IS NULL)`
	query := `
	SELECT COALESCE(SUM(price), 0) 
FROM subscriptions 
WHERE 
    start_date <= $2
    AND (end_date >= $1 OR end_date IS NULL)`

	args := []any{filter.From, filter.To}
	idx := 3

	if filter.UserID != nil {
		query += " AND user_id = $" + itoa(idx)
		args = append(args, *filter.UserID)
		idx++
	}
	if filter.ServiceName != nil {
		query += " AND service_name = $" + itoa(idx)
		args = append(args, *filter.ServiceName)
		idx++
	}
	var sum int
	row := r.pool.QueryRow(context.Background(), query, args...)
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}
	return sum, nil
}

// itoa — helper для конвертации int в string (без strconv для минимализма)
func itoa(i int) string {
	return string('0' + i)
}
