package comment

type CreateCommentRequest struct {
	UserID    uint
	RequestID uint
	Text      *string
	Rating    uint
}

type EditCommentRequest struct {
	UserID    uint
	CommentID uint
	Text      *string
	Rating    uint
}

type DeleteCommentRequest struct {
	UserID    uint
	CommentID uint
}

type GetAllPetSitterCommentsRequest struct {
	PetSitterID uint
}
