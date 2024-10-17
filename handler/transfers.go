package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhynever/wallet-test/services"
	"github.com/xhynever/wallet-test/services/responses"
)

func (h *Handler) CreateTransfer(ctx *gin.Context) {
	var req services.TxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	result, err := h.services.Accounts.CreateTransfer(req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse(err))
			return
		}
		if err == fmt.Errorf("currency mismatch") {
			ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) Business(ctx *gin.Context) {
	var req services.TxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	// fmt.Println(req)
	if req.Amount == 0 || req.Currency == "" {
		err := fmt.Errorf("参数问题")
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	if req.FromAccountID != req.ToAccountID {
		err := fmt.Errorf("是否为存取款业务")
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	// 创建交易
	result, err := h.services.Accounts.CreateTransfer(req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse(err))
			return
		}
		if err == fmt.Errorf("currency mismatch") {
			ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// WithDraw
// func (h *Handler) WithDraw(ctx *gin.Context) {
// 	var req services.TransferRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
// 		return
// 	}
// 	fmt.Println(req)
// 	if req.Amount < 0 || req.Currency == "" {
// 		err := fmt.Errorf("参数问题")
// 		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
// 		return
// 	}
// 	if req.FromAccountID != req.ToAccountID {
// 		err := fmt.Errorf("是否为存取款业务")
// 		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
// 		return
// 	}
// 	// 创建交易
// 	result, err := h.services.Accounts.CreateTransfer(req)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, responses.ErrorResponse(err))
// 			return
// 		}
// 		if err == fmt.Errorf("currency mismatch") {
// 			ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
// 		}
// 		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, result)

// }
