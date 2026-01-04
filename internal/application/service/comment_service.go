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
	unitOfWork       ports.UnitOfWork
	userService      usecase.UserService
	requestService   usecase.RequestService
	petSitterService usecase.PetSitterService
}

func NewCommentService(unitOfWork ports.UnitOfWork, userService usecase.UserService, requestService usecase.RequestService, petSitterService usecase.PetSitterService) *CommentService {
	return &CommentService{
		unitOfWork:       unitOfWork,
		userService:      userService,
		requestService:   requestService,
		petSitterService: petSitterService,
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
	return cs.updatePetSitterRatingOnCreate(foundRequest.PetSitterID, info.Rating)
}

func (cs *CommentService) EditComment(info comment.EditCommentRequest) error {
	foundComment, err := cs.FindUserCommentByID(info.CommentID, info.UserID)
	if err != nil {
		return err
	}

	oldRating := foundComment.Rating
	foundComment.Text = info.Text
	foundComment.Rating = info.Rating

	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	err = commentRepo.EditComment(foundComment)
	if err != nil {
		return err
	}

	return cs.updatePetSitterRatingOnEdit(foundComment.PetSitterID, oldRating, info.Rating)
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

	return cs.updatePetSitterRatingOnDelete(foundComment.PetSitterID, foundComment.Rating)
}

func (cs *CommentService) GetAllPetSitterComments(info comment.GetAllPetSitterCommentsRequest) (*comment.AllCommentsResponse, error) {
	commentRepo := cs.unitOfWork.Factory().CommentRepository()
	comments, err := commentRepo.GetAllPetSitterCommentsByID(info.PetSitterID)
	if err != nil {
		return nil, err
	}
	r := make([]comment.CommentResponse, len(comments))
	var averageRating float32 = 0
	for i := range comments {
		averageRating += float32(comments[i].Rating)
		response, err := buildCommentResponse(cs.userService, nil, &comments[i])
		if err != nil {
			return nil, err
		}
		if response == nil {
			continue
		}
		r[i] = *response
	}
	if len(r) != 0 {
		averageRating /= float32(len(r))
	}
	return &comment.AllCommentsResponse{
		CommentCount:  uint(len(r)),
		AverageRating: averageRating,
		Comments:      r,
	}, nil
}

func (cs *CommentService) GetAllPetSitterCommentsForPetSitter(userID uint) (*comment.AllCommentsResponse, error) {
	petSitter, err := cs.petSitterService.GetPetSitterByUserID(userID)
	if err != nil {
		return nil, err
	}

	info := comment.GetAllPetSitterCommentsRequest{
		PetSitterID: petSitter.ID,
	}
	return cs.GetAllPetSitterComments(info)
}

func buildCommentResponse(userService usecase.UserService, fallbackUser *entities.User, commentEntity *entities.Comment) (*comment.CommentResponse, error) {
	if commentEntity == nil {
		return nil, nil
	}

	user := fallbackUser
	if user == nil || user.ID != commentEntity.UserID {
		foundUser, err := userService.FindUserByID(commentEntity.UserID)
		if err != nil {
			return nil, err
		}
		user = foundUser
	}

	return &comment.CommentResponse{
		UserID:        user.ID,
		UserFirstName: user.FirstName,
		UserLastName:  user.LastName,
		Text:          commentEntity.Text,
		Rating:        commentEntity.Rating,
		UpdatedAt:     commentEntity.UpdatedAt,
	}, nil
}

func (cs *CommentService) updatePetSitterRatingOnCreate(petSitterID uint, rating uint) error {
	petSitter, err := cs.petSitterService.GetPetSitterByID(petSitterID)
	if err != nil {
		return err
	}

	newCount := petSitter.CommentsCount + 1
	total := petSitter.Rating*float32(petSitter.CommentsCount) + float32(rating)
	averageRating := total / float32(newCount)
	return cs.petSitterService.UpdateRatingAndCommentsCount(petSitterID, averageRating, newCount)
}

func (cs *CommentService) updatePetSitterRatingOnEdit(petSitterID uint, oldRating uint, newRating uint) error {
	petSitter, err := cs.petSitterService.GetPetSitterByID(petSitterID)
	if err != nil {
		return err
	}

	if petSitter.CommentsCount == 0 {
		return cs.petSitterService.UpdateRatingAndCommentsCount(petSitterID, 0, 0)
	}

	total := petSitter.Rating*float32(petSitter.CommentsCount) - float32(oldRating) + float32(newRating)
	averageRating := total / float32(petSitter.CommentsCount)
	return cs.petSitterService.UpdateRatingAndCommentsCount(petSitterID, averageRating, petSitter.CommentsCount)
}

func (cs *CommentService) updatePetSitterRatingOnDelete(petSitterID uint, rating uint) error {
	petSitter, err := cs.petSitterService.GetPetSitterByID(petSitterID)
	if err != nil {
		return err
	}

	if petSitter.CommentsCount <= 1 {
		return cs.petSitterService.UpdateRatingAndCommentsCount(petSitterID, 0, 0)
	}

	newCount := petSitter.CommentsCount - 1
	total := petSitter.Rating*float32(petSitter.CommentsCount) - float32(rating)
	averageRating := total / float32(newCount)
	return cs.petSitterService.UpdateRatingAndCommentsCount(petSitterID, averageRating, newCount)
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
