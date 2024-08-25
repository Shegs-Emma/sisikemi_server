package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

// =================================================================== Create a single media

type createMediaRequest struct {
	// MediaRef string `json:"media_ref" binding:"required"`
	Url string `json:"url" binding:"required"`
	AwsId string `json:"aws_id" binding:"required"`
}

func (server *Server) createMedia(ctx *gin.Context) {
	var req createMediaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateMediaParams {
		MediaRef: util.RandomString(10),
		Url: req.Url,
		AwsID: req.AwsId,
	}

	media, err := server.store.CreateMedia(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, media)
}

// =================================================================== Update media 
type updateMediaRequest struct {
	// MediaRef string `json:"media_ref" binding:"required"`
	ID int64 `json:"id" binding:"required"`
	Url string `json:"url" binding:"required"`
	AwsId string `json:"aws_id" binding:"required"`
}

func (server *Server) updateMedia(ctx *gin.Context) {
	var req updateMediaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	media, err := server.store.GetMedia(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateMediaParams {
		ID: media.ID,
		Url: req.Url,
		AwsID: req.AwsId,
	}

	mediaUpdate, err := server.store.UpdateMedia(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mediaUpdate)
}

// =================================================================== Get a single media 

type getMediaRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getMedia (ctx *gin.Context) {
	var req getMediaRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	media, err := server.store.GetMedia(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, media)
}


// =================================================================== Get  medium 
type listMediaRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listMedia (ctx *gin.Context) {
	var req listMediaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListMediaParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	media, err := server.store.ListMedia(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, media)
}

// =================================================================== Delete media 
type deleteMediaRequest struct {
	id int64 `uri:"id" binding:"required"`
}

func (server *Server) deleteMedia (ctx *gin.Context) {
	var req deleteMediaRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Print("req", req.id)

	err := server.store.DeleteMedia(ctx, req.id)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{"message": "Successfully Deleted"})
}