package api

import (
	"net/http"

	"github.com/NatdanaiKhe/simplebank/api/dto"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	pg "github.com/lib/pq"
)

func (server *Server) createUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, err)
		return
	}

	params := service.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
	}
	user, err := server.services.UserService.CreateUser(c, params)

	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				errorResponse(c, service.ErrDuplicateUsername)
				return
			}
		}

		errorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
