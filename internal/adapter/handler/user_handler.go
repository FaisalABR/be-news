package handler

import (
	"bwa-news/internal/adapter/handler/request"
	"bwa-news/internal/adapter/handler/response"
	"bwa-news/internal/core/domain/entity"
	"bwa-news/internal/core/service"
	validatorLib "bwa-news/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	GetUserByID(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)

	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetUserByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(userID))
	if err != nil {
		code = "[HANDLER] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	resp := response.UserResponse{
		ID:       int64(userID),
		Username: user.Username,
		Email:    user.Email,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Data fetched successfully"
	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}

func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)

	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var reqPassword request.UpdatePasswordRequest
	if err = c.BodyParser(&reqPassword); err != nil {
		code = "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(reqPassword); err != nil {
		code = "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.UserEntity{
		Password: reqPassword.NewPassword,
	}

	err = u.userService.UpdatePassword(c.Context(), reqEntity.Password, int64(userID))
	if err != nil {
		code = "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Password updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}
