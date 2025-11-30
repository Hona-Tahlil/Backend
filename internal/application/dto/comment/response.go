package comment

import "time"

type AllCommentsResponse struct {
	CommentCount  uint              `json:"commentCount"`
	AverageRating float32           `json:"averageRating"`
	Comments      []CommentResponse `json:"comments"`
}

type CommentResponse struct {
	UserName  string    `json:"userName"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      *string   `json:"text"`
	Rating    uint      `json:"rating"`
}
