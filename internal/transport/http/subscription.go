package http

import (
	"net/http"
	"time"

	"github.com/effective-mobile/subscriptions/internal/domain"
	"github.com/effective-mobile/subscriptions/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	uc usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(uc usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{uc: uc}
}

// SubscriptionCreateRequest — схема запроса на создание подписки
// swagger:model
// start_date и end_date — формат "MM-YYYY"
type SubscriptionCreateRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"required,gte=0"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SumRequest struct {
	UserID      string `query:"user_id"`
	ServiceName string `query:"service_name"`
	From        string `query:"from"` // MM-YYYY
	To          string `query:"to"`   // MM-YYYY
}

type SumResponse struct {
	Sum int `json:"sum"`
}

func (h *SubscriptionHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/subscriptions", h.List)
	r.Get("/subscriptions/sum", h.Sum)
	r.Post("/subscriptions", h.Create)
	r.Get("/subscriptions/:id", h.GetByID)
	r.Put("/subscriptions/:id", h.Update)
	r.Delete("/subscriptions/:id", h.Delete)
}

func (h *SubscriptionHandler) Create(c *fiber.Ctx) error {
	var req SubscriptionCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid start_date")
	}
	var endDate *time.Time
	if req.EndDate != "" {
		e, err := parseMonthYear(req.EndDate)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "invalid end_date")
		}
		endDate = &e
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user_id")
	}
	sub := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	if err := h.uc.Create(c.Context(), sub); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(toResponse(sub))
}

func (h *SubscriptionHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	sub, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, err.Error())
	}
	return c.JSON(toResponse(sub))
}

func (h *SubscriptionHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	var req SubscriptionCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid start_date")
	}
	var endDate *time.Time
	if req.EndDate != "" {
		e, err := parseMonthYear(req.EndDate)
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, "invalid end_date")
		}
		endDate = &e
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid user_id")
	}
	sub := &domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
		UpdatedAt:   time.Now(),
	}
	if err := h.uc.Update(c.Context(), sub); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(toResponse(sub))
}

func (h *SubscriptionHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	if err := h.uc.Delete(c.Context(), id); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(c *fiber.Ctx) error {
	var filter domain.SubscriptionFilter
	if userID := c.Query("user_id"); userID != "" {
		u, err := uuid.Parse(userID)
		if err == nil {
			filter.UserID = &u
		}
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}
	if from := c.Query("from"); from != "" {
		if t, err := parseMonthYear(from); err == nil {
			filter.From = &t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := parseMonthYear(to); err == nil {
			filter.To = &t
		}
	}
	subs, err := h.uc.List(c.Context(), filter)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	resp := make([]SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		resp = append(resp, toResponse(&s))
	}
	return c.JSON(resp)
}

func (h *SubscriptionHandler) Sum(c *fiber.Ctx) error {
	var req SumRequest
	if err := c.QueryParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid query params")
	}
	var (
		userID      *uuid.UUID
		serviceName *string
		from, to    time.Time
		err         error
	)
	if req.UserID != "" {
		u, err := uuid.Parse(req.UserID)
		if err == nil {
			userID = &u
		}
	}
	if req.ServiceName != "" {
		serviceName = &req.ServiceName
	}
	if req.From == "" || req.To == "" {
		return fiber.NewError(http.StatusBadRequest, "from and to are required")
	}
	from, err = parseMonthYear(req.From)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid from")
	}
	to, err = parseMonthYear(req.To)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid to")
	}
	filter := domain.SubscriptionSumFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
	}
	sum, err := h.uc.Sum(c.Context(), filter)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(SumResponse{Sum: sum})
}

func toResponse(sub *domain.Subscription) SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sub.UpdatedAt.Format(time.RFC3339),
	}
	if sub.EndDate != nil {
		resp.EndDate = sub.EndDate.Format("01-2006")
	}
	return resp
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}
