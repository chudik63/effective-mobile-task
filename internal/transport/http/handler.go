package http

import "github.com/gin-gonic/gin"

type Service interface {
}

type Handler struct {
}

func NewHandler(router *gin.Engine, service Service) {

}
