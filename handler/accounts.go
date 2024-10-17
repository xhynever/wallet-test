package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xhynever/wallet-test/services"
	"github.com/xhynever/wallet-test/services/responses"
)

func (h *Handler) CreateAccount(ctx *gin.Context) {
	// var req services.GetAccountRequest
	// if err := ctx.ShouldBindUri(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
	// 	return
	// }
	var req services.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("参数解析错误")
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}

	if req.Owner == "" {
		err := fmt.Errorf("参数为空")
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	account, err := h.services.Accounts.CreateAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"account": account,
	})
}

func (h *Handler) GetAccount(ctx *gin.Context) {
	var req services.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	account, err := h.services.Accounts.GetAccount(req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (h *Handler) UpdateAccount(ctx *gin.Context) {
	var req services.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	account, err := h.services.Accounts.UpdateAccount(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (h *Handler) DeleteAccount(ctx *gin.Context) {
	var req services.DeleteRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	//  获得账户信息。删除账户余额。删除账户
	err := h.services.Accounts.DeleteAccount(req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("account %v deleted successfully", req.ID),
	})
}

func (h *Handler) ListAccounts(ctx *gin.Context) {
	var req services.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	accounts, err := h.services.Accounts.ListAccounts(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
func (h *Handler) ListEntries(ctx *gin.Context) {
	var req services.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse(err))
		return
	}
	fmt.Println(req)
	accounts, err := h.services.Accounts.ListEntries(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
