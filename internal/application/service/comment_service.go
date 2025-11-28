package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type CommentService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
}

func NewCommentService(unitOfWork ports.UnitOfWork, userService usecase.UserService) *CommentService {
	return &CommentService{
		unitOfWork:  unitOfWork,
		userService: userService,
	}
}

func (cs *CommentService) CreateComment(info comment.CreateCommentRequest) error {
	// TODO: Check if the request exists

	// TODO: check if it's status is finished

	newComment := &entities.Comment{
		UserID: info.UserID,
		// PetSitterID: ,
		// RequestID: ,
		Text:   info.Text,
		Rating: info.Rating,
	}

	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	err := commentRepo.CreateComment(newComment)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) EditComment(info comment.EditCommentRequest) error {
	foundComment, err := cs.FindCommentByID(info.CommentID)
	if err != nil {
		return err
	}

	if foundComment.UserID != info.UserID {
		return exceptions.NewAccessDeniedError("you can't edit this comment")
	}

	foundComment.Text = info.Text
	foundComment.Rating = info.Rating

	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	err = commentRepo.EditComment(foundComment)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CommentService) FindCommentByID(id uint) (*entities.Comment, error) {
	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	foundComment, err := commentRepo.FindCommentByID(id)
	if err != nil {
		return nil, err
	}
	if foundComment == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Comment)
	}
	return foundComment, nil
}

func (cs *CommentService) DeleteComment(info comment.DeleteCommentRequest) error {
	foundComment, err := cs.FindCommentByID(info.CommentID)
	if err != nil {
		return err
	}

	if foundComment.UserID != info.UserID {
		return exceptions.NewAccessDeniedError("you can't edit this comment")
	}

	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	err = commentRepo.DeleteComment(foundComment)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CommentService) GetAllPetSitterComments(info comment.GetAllPetSitterCommentsRequest) (*comment.AllCommentsResponse, error) {
	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	comments, err := commentRepo.GetAllPetSitterCommentsByID(info.PetSitterID)
	if err != nil {
		return nil, err
	}
	r := make([]comment.CommentResponse, 0)
	var averageRating float32 = 0
	for _, c := range comments {
		user, err := cs.userService.FindUserByID(c.UserID)
		if err != nil {
			return nil, err
		}
		name := user.FirstName + " " + user.LastName
		averageRating += float32(c.Rating)
		r = append(r, comment.CommentResponse{
			UserName: name,
			Text:     c.Text,
			Rating:   c.Rating,
		})
	}
	averageRating /= float32(len(r))
	return &comment.AllCommentsResponse{
		CommentCount:  uint(len(r)),
		AverageRating: averageRating,
		Comments:      r,
	}, nil
}
