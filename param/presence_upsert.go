package param

type UpsertPresenceRequest struct {
	UserID    uint  `json:"userId"`
	Timestamp int64 `json:"timestamp"`
}
type UpsertPresenceResponse struct{}
