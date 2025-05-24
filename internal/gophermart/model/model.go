package model

import "time"

type OrderEntity struct {
	ID          int64
	UserLogin   string
	OrderNumber int64
	Accrual     *float32
	Status      string
	CreatedAt   time.Time
}

type WithdrawHistoryEntity struct {
	ID          int64
	UserLogin   string
	OrderNumber int64
	Sum         int64
	CreatedAt   time.Time
}

type BalanceEntity struct {
	ID           int64
	UserLogin    string
	WithdrawnSum float32
	CurrentSum   float32
	CreatedAt    time.Time
}

type WithdrawRequest struct {
	Order string  `json:"order" validate:"required,min=1,max=255"`
	Sum   float32 `json:"sum" validate:"required"`
}
type RegisterUserRequest struct {
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}

type LoginUserRequest struct {
	Login    string `json:"login" validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=1,max=255"`
}

type OrderItemResponse struct {
	OrderNumber string  `json:"number" validate:"required,min=1,max=255"`
	Status      string  `json:"status" validate:"required,min=1,max=255"`
	Accrual     float32 `json:"accrual" validate:"required,min=1,max=255"`
	UploadedAt  string  `json:"uploaded_at" validate:"required,min=1,max=255"`
}

type OrderAccrualResponse struct {
	OrderNumber string  `json:"order" validate:"required,min=1,max=255"`
	Status      string  `json:"status" validate:"required,min=1,max=255"`
	Accrual     float32 `json:"accrual" validate:"required,min=1,max=255"`
}

type WithdrawHistoryItemResponse struct {
	OrderNumber string  `json:"order" validate:"required,min=1,max=255"`
	Sum         float32 `json:"sum" validate:"required"`
	ProcessedAt string  `json:"processed_at" validate:"required,min=1,max=255"`
}

type GetBalanceResponse struct {
	Current   float32 `json:"current" validate:"required"`
	Withdrawn float32 `json:"withdrawn" validate:"required"`
}

const (
	NewOrderStatus        = "NEW"
	ProcessingOrderStatus = "PROCESSING"
	InvalidOrderStatus    = "INVALID"
	ProcessedOrderStatus  = "PROCESSED"
)
