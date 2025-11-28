package user

import (
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type UserCommentController struct {
	commentService *service.CommentService
}

func NewUserCommentController(commentService *service.CommentService) *UserCommentController {
	return &UserCommentController{
		commentService: commentService,
	}
}

func (cc *UserCommentController) CreateComment(ctx *gin.Context) {
	type Params struct {
		RequestID uint    `json:"requestID" validate:"required"`
		Text      *string `json:"text"`
		Rating    uint    `json:"rating" validate:"required,min=1,max=5"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)

	info := comment.CreateCommentRequest{
		UserID:    UserID,
		RequestID: params.RequestID,
		Text:      params.Text,
		Rating:    params.Rating,
	}

	if err := cc.commentService.CreateComment(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (cc *UserCommentController) EditComment(ctx *gin.Context) {
	type Params struct {
		CommentID uint    `json:"commentID" validate:"required"`
		Text      *string `json:"text"`
		Rating    uint    `json:"rating" validate:"required,min=1,max=5"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)

	info := comment.EditCommentRequest{
		UserID:    UserID,
		CommentID: params.CommentID,
		Text:      params.Text,
		Rating:    params.CommentID,
	}

	if err := cc.commentService.EditComment(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (cc *UserCommentController) DeleteComment(ctx *gin.Context) {
	type Params struct {
		CommentID uint `uri:"commentID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)
	info := comment.DeleteCommentRequest{
		UserID:    UserID,
		CommentID: params.CommentID,
	}

	if err := cc.commentService.DeleteComment(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (cc *UserCommentController) GetAllPetSitterComments(ctx *gin.Context) {
	type Params struct {
		PetSitterID uint `uri:"petSitterID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)

	info := comment.GetAllPetSitterCommentsRequest{
		PetSitterID: params.PetSitterID,
	}

	res, err := cc.commentService.GetAllPetSitterComments(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

// TODO: comment info in request full data
