package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type API struct {
	Router *gin.Engine
	I18n   *i18n.Bundle
}
