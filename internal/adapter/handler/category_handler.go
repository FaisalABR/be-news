package handler

import (
	"bwa-news/internal/adapter/handler/request"
	"bwa-news/internal/adapter/handler/response"
	"bwa-news/internal/core/domain/entity"
	"bwa-news/internal/core/service"
	validatorLib "bwa-news/lib/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var errResp response.ErrorResponseDefault
var defaultSuccessResponse response.DefaultSuccessResponse

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategoryByID(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error

	// FE
	GetFeCategories(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

func (ch *categoryHandler) GetFeCategories(c *fiber.Ctx) error {
	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetFeCategories - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Username,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Categories fetched successfully"
	defaultSuccessResponse.Data = categoryResponses

	return c.JSON(defaultSuccessResponse)
}

func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Username,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Categories fetched successfully"
	defaultSuccessResponse.Data = categoryResponses

	return c.JSON(defaultSuccessResponse)

}
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {

	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	results, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID:            results.ID,
		Title:         results.Title,
		Slug:          results.Slug,
		CreatedByName: results.User.Username,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Category fetched successfully"
	defaultSuccessResponse.Data = categoryResponse

	return c.JSON(defaultSuccessResponse)

}
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Pagination = nil
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Category created successfully"

	return c.JSON(defaultSuccessResponse)

}
func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] EditCategoryByID - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditCategoryByID - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid body request"
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	err = validatorLib.ValidateStruct(req)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:    id,
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	if err = ch.categoryService.EditCategoryByID(c.Context(), reqEntity); err != nil {
		code = "[HANDLER] EditCategoryByID - 5"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Category updated successfully"
	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Pagination = nil

	return c.JSON(defaultSuccessResponse)
}
func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] DeleteCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errResp)
	}

	err = ch.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultSuccessResponse.Data = nil
	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Deleted category successfully"
	defaultSuccessResponse.Pagination = nil

	return c.JSON(defaultSuccessResponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}
