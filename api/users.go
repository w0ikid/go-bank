package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/w0ikid/go-bank/db/sqlc"
	"github.com/w0ikid/go-bank/util"
)

type createUserParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string `json:"username"`
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	PasswordChangedAt string `json:"password_changed_at"`
	CreatedAt         string `json:"created_at"`
}

func (server *Server) createUsers(c *gin.Context) {
	var req createUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashpassword, err := util.HashPassword(req.Password)
	if err != nil {
		return 
	}
		
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashpassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.Time.String(),
		CreatedAt:         user.CreatedAt.Time.String(),
	}

	c.JSON(http.StatusOK, response)
}