package petsitter

import (
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterCommentController struct {
	commentService usecase.CommentService
}

func NewPetSitterCommentController(commentService usecase.CommentService) *PetSitterCommentController {
	return &PetSitterCommentController{
		commentService: commentService,
	}
}

func (cc *PetSitterCommentController) GetAllComments(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	res, err := cc.commentService.GetAllPetSitterCommentsForPetSitter(UserID)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
