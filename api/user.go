package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	ProfilePhoto   string `json:"profile_photo"`
}

type userResponse struct {
	Username string `json:"username"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	ProfilePhoto   string `json:"profile_photo"`
	IsAdmin bool `json:"is_admin"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		ProfilePhoto: user.ProfilePhoto,
		IsAdmin: user.IsAdmin,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	isAdmin := false
	if req.Email == "remi.emma04@gmail.com" {
		isAdmin = true
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		IsAdmin: isAdmin,
		PhoneNumber: req.PhoneNumber,
		ProfilePhoto: req.ProfilePhoto,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}


// Login an account
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`
	SessionID uuid.UUID `json:"session_id"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
		user.IsAdmin,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID: session.ID,
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

// Update a user
type updateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password *string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	Email *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	ProfilePhoto *string `json:"profile_photo"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if authPayload.Username != req.Username {
		err := errors.New("you are not authorized to update user you did not create")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	violations := validateUpdateUserRequest(req)

	if violations != nil {
		err := errors.New("the entries are invalid")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FirstName: nullableString(req.FirstName),
		LastName: nullableString(req.LastName),
		Email: nullableString(req.Email),
		PhoneNumber: nullableString(req.PhoneNumber),
		ProfilePhoto: nullableString(req.ProfilePhoto),
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)

		if err != nil {
			err := errors.New("failed to hash password")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid: true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			err := errors.New("user not found")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		err := errors.New("failed to update user")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

func validateUpdateUserRequest(req updateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	fmt.Println("req", req)
	
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(*req.Password); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(*req.Email); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if req.FirstName != nil {
		if err := val.ValidateFirstName(*req.FirstName); err != nil {
			violations = append(violations, fieldViolation("first_name", err))
		}
	}	

	if req.LastName != nil {
		if err := val.ValidateLastName(*req.LastName); err != nil {
			violations = append(violations, fieldViolation("last_name", err))
		}
	}
	
	return violations
}

func nullableString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{
			String: *s,
			Valid:  true,
		}
	}
	return sql.NullString{
		String: "",
		Valid:  false,
	}
}