package domainpostgres

import "hona/backend/internal/domain/entities"

type CommentRepository interface {
	CreateComment(comment *entities.Comment) error
	FindCommentByID(id uint) (*entities.Comment, error)
	EditComment(comment *entities.Comment) error
	DeleteComment(comment *entities.Comment) error
	GetAllPetSitterCommentsByID(id uint) ([]entities.Comment, error)
}
