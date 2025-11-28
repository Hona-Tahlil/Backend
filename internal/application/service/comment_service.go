package service

import "hona/backend/internal/application/dto/comment"

type CommentService struct {
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (cs *CommentService) CreateComment(info comment.CreateCommentRequest) error {

	return nil
}
