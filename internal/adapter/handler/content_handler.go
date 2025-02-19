package handler

import (
	"bwa-news/internal/adapter/handler/request"
	"bwa-news/internal/adapter/handler/response"
	"bwa-news/internal/core/domain/entity"
	"bwa-news/internal/core/service"
	"bwa-news/lib/conv"
	validatorLib "bwa-news/lib/validator"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error

	// FE
	GetContentsWithQuery(c *fiber.Ctx) error
	GetContentDetails(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

func (ch *contentHandler) GetContentDetails(c *fiber.Ctx) error {
	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] GetContentDetails - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContentDetails - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Status:       result.Status,
		Tags:         result.Tags,
		CreatedAt:    result.CreatedAt.Local().Format("02 January 2006"),
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CategoryID,
		CategoryName: result.Category.Title,
		Author:       result.User.Username,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Data fetched successfully"
	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) GetContentsWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[HANDLER] GetContentsWithQuery - 1"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	limit := 10
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[HANDLER] GetContentsWithQuery - 2"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "DESC"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("content_id") != "" {
		categoryID, err = conv.StringToInt(c.Query("content_id"))
		if err != nil {
			code = "[HANDLER] GetContactByID - 2"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
	}

	query := entity.QueryString{
		Page:       page,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		CategoryID: categoryID,
	}

	results, err := ch.contentService.GetContents(c.Context(), query)
	if err != nil {
		code = "[HANDLER] GetContentsWithQuery - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	contentResps := []response.SuccessContentResponse{}

	for _, result := range results {
		contentResps = append(contentResps, response.SuccessContentResponse{
			ID:           result.ID,
			Title:        result.Title,
			Excerpt:      result.Excerpt,
			Description:  result.Description,
			Image:        result.Image,
			Tags:         result.Tags,
			Status:       result.Status,
			CategoryID:   result.CategoryID,
			CreatedByID:  result.CreatedByID,
			CreatedAt:    result.CreatedAt.Local().Format("02 January 2006"),
			CategoryName: result.Category.Title,
			Author:       result.User.Username,
		})
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content fetched successfully"
	defaultSuccessResponse.Data = contentResps

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)

	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[HANDLER] GetContents - 2"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	limit := 10
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[HANDLER] GetContents - 3"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "DESC"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	query := entity.QueryString{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		OrderType: orderType,
		Search:    search,
	}

	results, err := ch.contentService.GetContents(c.Context(), query)
	if err != nil {
		code = "[HANDLER] GetContents - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	contentResps := []response.SuccessContentResponse{}

	for _, result := range results {
		contentResps = append(contentResps, response.SuccessContentResponse{
			ID:           result.ID,
			Title:        result.Title,
			Excerpt:      result.Excerpt,
			Description:  result.Description,
			Image:        result.Image,
			Tags:         result.Tags,
			Status:       result.Status,
			CategoryID:   result.CategoryID,
			CreatedByID:  result.CreatedByID,
			CreatedAt:    result.CreatedAt.Local().Format("02 January 2006"),
			CategoryName: result.Category.Title,
			Author:       result.User.Username,
		})
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content fetched successfully"
	defaultSuccessResponse.Data = contentResps

	return c.JSON(defaultSuccessResponse)
}
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetContactByID - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] GetContactByID - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContactByID - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Status:       result.Status,
		Tags:         result.Tags,
		CreatedAt:    result.CreatedAt.Local().Format("02 January 2006"),
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CategoryID,
		CategoryName: result.Category.Title,
		Author:       result.User.Username,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Data fetched successfully"
	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), id)
	if err != nil {
		code = "[HANDLER] DeleteContent - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultSuccessResponse.Meta.Status = false
	defaultSuccessResponse.Meta.Message = "Content deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}
func (ch *contentHandler) UpdateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] UpdateContent - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] UpdateContent - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] UpdateContent - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	idString := c.Params("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		code = "[HANDLER] UpdateContent - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          id,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = ch.contentService.UpdateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] UpdateContent - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content updated successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (ch *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JWTData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] UploadImageR2 - 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	var req request.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Invalid body request"

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code = "[HANDLER] UploadImageR2 - 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(userID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := ch.contentService.UploadImageR2(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code = "[HANDLER] UploadImageR2 - 5"
			log.Errorw(code, err)
			errResp.Meta.Status = false
			errResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
	}

	urlImageResp := map[string]interface{}{
		"urlImage": imageUrl,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = urlImageResp

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
