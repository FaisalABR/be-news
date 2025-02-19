package service

import (
	"bwa-news/config"
	"bwa-news/internal/adapter/cloudflare"
	"bwa-news/internal/adapter/repository"
	"bwa-news/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	DeleteContent(ctx context.Context, id int64) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepository repository.ContentRepository
	cfg               *config.Config
	r2                cloudflare.CloudflareR2Adapter
}

func (c *contentService) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error) {
	results, err := c.contentRepository.GetContents(ctx, query)
	if err != nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepository.GetContentByID(ctx, id)
	if err != nil {
		code = "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	err := c.contentRepository.DeleteContent(ctx, id)
	if err != nil {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.UpdateContent(ctx, req)
	if err != nil {
		code = "[SERVICE] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.CreateContent(ctx, req)
	if err != nil {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	urlImage, err := c.r2.UploadImage(&req)
	if err != nil {
		code = "[SERVICE] UploadImageR2 - 1"
		log.Errorw(code, err)
		return "", err
	}

	return urlImage, nil
}

func NewContentService(contentRepository repository.ContentRepository, cfg *config.Config,
	r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepository: contentRepository,
		cfg:               cfg,
		r2:                r2,
	}
}
