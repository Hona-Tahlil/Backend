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
	unitOfWork     ports.UnitOfWork
	userService    usecase.UserService
	requestService usecase.RequestService
}

func NewCommentService(unitOfWork ports.UnitOfWork, userService usecase.UserService, requestService usecase.RequestService) *CommentService {
	return &CommentService{
		unitOfWork:     unitOfWork,
		userService:    userService,
		requestService: requestService,
	}
}

func (cs *CommentService) CreateComment(info comment.CreateCommentRequest) error {
	foundRequest, err := cs.requestService.FindRequestByID(info.RequestID)
	if err != nil {
		return err
	}

	err = cs.requestService.EnsureRequestIsFinished(foundRequest)
	if err != nil {
		return err
	}
	if foundRequest.Comment != nil {
		var ce exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Comment, bootstrap.Run().Constants.ErrorTags.AlreadyExist)
		return &ce
	}

	newComment := &entities.Comment{
		UserID:      info.UserID,
		PetSitterID: foundRequest.PetSitterID,
		RequestID:   foundRequest.ID,
		Text:        info.Text,
		Rating:      info.Rating,
	}

	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	err = commentRepo.CreateComment(newComment)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) EditComment(info comment.EditCommentRequest) error {
	foundComment, err := cs.FindUserCommentByID(info.CommentID, info.UserID)
	if err != nil {
		return err
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
	foundComment, err := cs.FindUserCommentByID(info.CommentID, info.UserID)
	if err != nil {
		return err
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
		averageRating += float32(c.Rating)
		r = append(r, comment.CommentResponse{
			UserID:        user.ID,
			UserFirstName: user.FirstName,
			UserLastName:  user.LastName,
			Text:          c.Text,
			Rating:        c.Rating,
			UpdatedAt:     c.UpdatedAt,
		})
	}
	averageRating /= float32(len(r))
	return &comment.AllCommentsResponse{
		CommentCount:  uint(len(r)),
		AverageRating: averageRating,
		Comments:      r,
	}, nil
}

func (cs *CommentService) FindUserCommentByID(commentID, userID uint) (*entities.Comment, error) {
	foundComment, err := cs.FindCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	if foundComment.UserID != userID {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Comment)
	}
	return foundComment, nil
}
