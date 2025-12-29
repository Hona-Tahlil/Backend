package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (cr *CommentRepository) CreateComment(comment *entities.Comment) error {
	return cr.db.Create(comment).Error
}

func (cr *CommentRepository) FindCommentByID(id uint) (*entities.Comment, error) {
	var foundComment entities.Comment

	if result := cr.db.First(&foundComment, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundComment, nil
}

func (cr *CommentRepository) EditComment(comment *entities.Comment) error {
	return cr.db.Save(comment).Error
}

func (cr *CommentRepository) DeleteComment(comment *entities.Comment) error {
	return cr.db.Delete(comment).Error
}

func (cr *CommentRepository) GetAllPetSitterCommentsByID(id uint) ([]entities.Comment, error) {
	var comments []entities.Comment

	if err := cr.db.Find(&comments, "pet_sitter_id = ?", id).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
