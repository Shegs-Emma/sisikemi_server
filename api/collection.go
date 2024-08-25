package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
)

type createCollectionRequest struct {
	CollectionName string `json:"collection_name" binding:"required"`
}

func (server *Server) createCollection(ctx *gin.Context) {
	var req createCollectionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	collection, err := server.store.CreateCollection(ctx, req.CollectionName)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !authPayload.IsAdmin {
		err := errors.New("you are not authorized to create")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// =================== Get collection =================
type getCollectionRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getCollection (ctx *gin.Context) {
	var req getCollectionRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	collection, err := server.store.GetCollection(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, collection)
}

// ==== ============ List collections ==============
type listCollectionRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCollections (ctx *gin.Context) {
	var req listCollectionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCollectionParams {
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	collection, err := server.store.ListCollection(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collection)
}