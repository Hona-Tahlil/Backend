package controllers

import (
	"hona/backend/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Text   string
	Params []string
}

type singleMessageResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type multipleMessageResponse struct {
	StatusCode int               `json:"statusCode"`
	Messages   map[string]string `json:"messages"`
	Data       interface{}       `json:"data"`
}

func Respond[T Message | []Message](ctx *gin.Context, statusCode int, messages T, data interface{}) {
	translator := GetTranslator(ctx, bootstrap.ProjectConfig.Constants.Context.Translator)

	switch msg := any(messages).(type) {
	case Message:
		if msg.Text == "" {
			msg.Text = http.StatusText(statusCode)
		}
		message, _ := translator.T(msg.Text, msg.Params...)
		ctx.JSON(statusCode, singleMessageResponse{
			StatusCode: statusCode,
			Message:    message,
			Data:       data,
		})
	case []Message:
		mms := multipleMessageResponse{
			StatusCode: statusCode,
			Messages:   map[string]string{},
			Data:       data,
		}
		for _, ms := range msg {
			translatedFieldValue, _ := translator.T(ms.Params[0])
			translatedTag, _ := translator.T(ms.Text, translatedFieldValue)
			mms.Messages[ms.Params[0]] = translatedTag
		}
		ctx.JSON(statusCode, mms)
	}
}
