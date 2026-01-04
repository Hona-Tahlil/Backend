package usecase

import (
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/domain/entities"
)

type CommentService interface {
	CreateComment(info comment.CreateCommentRequest) error
	EditComment(info comment.EditCommentRequest) error
	FindUserCommentByID(commentID uint, userID uint) (*entities.Comment, error)
	DeleteComment(info comment.DeleteCommentRequest) error
	GetAllPetSitterComments(info comment.GetAllPetSitterCommentsRequest) (*comment.AllCommentsResponse, error)
	GetAllPetSitterCommentsForPetSitter(userID uint) (*comment.AllCommentsResponse, error)
	FindCommentByID(id uint) (*entities.Comment, error)
}
