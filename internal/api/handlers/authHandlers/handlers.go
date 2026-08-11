package authhandlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"todoapp/internal/api/helpers"
	"todoapp/internal/api/models"
	"todoapp/internal/metrics"
)

func (u *AuthHandler) SignUp(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.SignUp")
	defer span.End()

	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "SignUp")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignUp: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if user.Name == "" || user.Pass == "" || user.Email == "" {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "SignUp")
		span.RecordError(fmt.Errorf("name, email and password are required"))
		u.logger.Error("authhandlers.SignUp: empty required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userId, err := u.userService.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.metrics.RecordError(string(metrics.ErrInternal), "SignUp")
			span.RecordError(err)
			u.logger.Error("authhandlers.SignUp: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.metrics.RecordError(string(metrics.ErrInternal), "SignUp")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignUp: add user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	token, err := u.jwtService.GenerateToken(userId, user.Name, user.Email)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "SignUp")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignUp: generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	u.logger.Info("authhandlers.SignUp: success", "user", user.Email)

	c.JSON(http.StatusCreated, gin.H{"user_id": userId})
}

func (u *AuthHandler) SignIn(c *gin.Context) {
	ctx, span := u.tracer.Start(c.Request.Context(), "handler.SignIn")
	defer span.End()

	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "SignIn")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignIn: bind json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if user.Pass == "" || user.Email == "" {
		u.metrics.RecordError(string(metrics.ErrBadRequest), "SignIn")
		span.RecordError(fmt.Errorf("email and password are required"))
		u.logger.Error("authhandlers.SignIn: empty required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	foundUser, err := u.userService.CheckAndGetUser(ctx, user)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			u.metrics.RecordError(string(metrics.ErrInternal), "SignIn")
			span.RecordError(err)
			u.logger.Error("authhandlers.SignIn: timeout", "error", err)
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			return
		}
		u.metrics.RecordError(string(metrics.ErrForbidden), "SignIn")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignIn: check user", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token, err := u.jwtService.GenerateToken(foundUser.Id, foundUser.Name, foundUser.Email)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrInternal), "SignIn")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignIn: generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	u.logger.Info("authhandlers.SignIn: success", "user", user.Email)

	c.JSON(http.StatusOK, gin.H{"user": foundUser})
}

func (u *AuthHandler) SignOut(c *gin.Context) {
	_, span := u.tracer.Start(c.Request.Context(), "handler.SignOut")
	defer span.End()

	userId, err := helpers.GetUserIdFromContext(c)
	if err != nil {
		u.metrics.RecordError(string(metrics.ErrForbidden), "SignOut")
		span.RecordError(err)
		u.logger.Error("authhandlers.SignOut: get user id", "error", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.SetCookie("token", "1", -1, "/", "", false, false)

	u.logger.Info("authhandlers.SignOut: success", "user", userId)

	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}