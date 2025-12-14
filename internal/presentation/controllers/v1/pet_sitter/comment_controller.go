package petsitter

import (
	"hona/backend/internal/application/dto/comment"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterCommentController struct {
	commentService *service.CommentService
}

func NewPetSitterCommentController(commentService *service.CommentService) *PetSitterCommentController {
	return &PetSitterCommentController{
		commentService: commentService,
	}
}

func (cc *PetSitterCommentController) GetAllComments(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)

	info := comment.GetAllPetSitterCommentsRequest{
		PetSitterID: UserID,
	}

	res, err := cc.commentService.GetAllPetSitterComments(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
