package comment

type AllCommentsResponse struct {
	CommentCount  uint              `json:"commentCount"`
	AverageRating float32           `json:"averageRating"`
	Comments      []CommentResponse `json:"comments"`
}

type CommentResponse struct {
	UserName string `json:"userName"`
	// date
	Text   *string `json:"text"`
	Rating uint    `json:"rating"`
}
