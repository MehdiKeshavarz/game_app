package param

import (
	"game_app/entity"
	"time"
)

type AddToWaitingListRequest struct {
	UserID   uint            `json:"userId"`
	Category entity.Category `json:"category"`
}

type AddToWaitingListResponse struct {
	TimeOut time.Duration `json:"timeout_in_nanoseconds"`
}
