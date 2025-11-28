package user

import (
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService *service.CommentService
}

func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

func (cc *CommentController) CreateComment(ctx *gin.Context) {
	type Params struct {
		RequestID uint    `json:"requestID" validate:"required"`
		Text      *string `json:"text"`
		Rating    uint    `json:"rating" validate:"required,min=1,max=5"`
	}
	params := controllers.Receive[Params](ctx)

	info := comment.CreateCommentRequest{
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
