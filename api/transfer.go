package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/w0ikid/go-bank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,min=1"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req createTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Validate accounts exist
	fromAccount, err := server.store.GetAccount(c, req.FromAccountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "from account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	toAccount, err := server.store.GetAccount(c, req.ToAccountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "to account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Create transfer
	if fromAccount.Currency != toAccount.Currency {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currency mismatch"})
		return
	}
	if fromAccount.ID == toAccount.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot transfer to the same account"})
		return
	}

	transfer, err := server.store.CreateTransfer(c, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "transfer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transfer)
}
