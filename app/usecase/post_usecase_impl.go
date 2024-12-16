package usecase

import (
	"app/domain"
	"app/pkg/common"
	"context"
	"fmt"
	"time"
)

type PostUsecaseImpl struct {
	postRepository domain.PostRepository
	transactor     domain.Transactor
}

func NewPostUsecaseImpl(postRepository domain.PostRepository, transactor domain.Transactor) domain.PostUseCase {
	return &PostUsecaseImpl{
		postRepository: postRepository,
		transactor:     transactor,
	}
}

// Create implements domain.PostUseCase.
func (uc *PostUsecaseImpl) Create(ctx context.Context, post *domain.CreatePostRequestDTO) (*domain.CreatePostResponseDTO, error) {
	post.Content = common.Sanitize(post.Content)
	tx, err := uc.transactor.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	postModel := &domain.Post{
		Title:    post.Title,
		Content:  post.Content,
		AuthorID: post.AuthorID,
	}
	err = uc.postRepository.Create(ctx, tx, postModel)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	res := &domain.CreatePostResponseDTO{
		ID:        postModel.ID,
		Title:     postModel.Title,
		Content:   postModel.Content,
		AuthorID:  postModel.AuthorID,
		CreatedAt: postModel.CreatedAt,
	}
	return res, nil
}

// Delete implements domain.PostUseCase.
func (uc *PostUsecaseImpl) Delete(ctx context.Context, id int64, post *domain.DeletePostRequestDTO) error {
	tx, err := uc.transactor.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	postModel, err := uc.postRepository.SelectForUpdate(ctx, tx, id)
	fmt.Println(postModel, err)
	if err != nil {
		return err
	}
	if postModel == nil {
		return common.ErrPostNotFound
	}
	if postModel.AuthorID != post.AuthorID {
		return common.ErrPostOwnerMismatch
	}
	now := time.Now()
	postModel.DeletedAt = &now
	err = uc.postRepository.Update(ctx, tx, id, postModel)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetAll implements domain.PostUseCase.
func (uc *PostUsecaseImpl) GetAll(ctx context.Context, search domain.SearchParam) ([]domain.Post, error) {
	posts, err := uc.postRepository.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetByID implements domain.PostUseCase.
func (uc *PostUsecaseImpl) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	post, err := uc.postRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, common.ErrPostNotFound
	}
	return post, nil
}

// Update implements domain.PostUseCase.
func (uc *PostUsecaseImpl) Update(ctx context.Context, id int64, post *domain.UpdatePostRequestDTO) (*domain.UpdatePostResponseDTO, error) {
	post.Content = common.Sanitize(post.Content)
	tx, err := uc.transactor.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	postModel, err := uc.postRepository.SelectForUpdate(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if postModel == nil {
		return nil, common.ErrPostNotFound
	}
	if postModel.AuthorID != post.AuthorID {
		fmt.Println(postModel.AuthorID, post.AuthorID)
		return nil, common.ErrPostOwnerMismatch
	}
	now := time.Now()
	postModel.Title = post.Title
	postModel.Content = post.Content
	postModel.UpdatedAt = &now
	err = uc.postRepository.Update(ctx, tx, id, postModel)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	res := &domain.UpdatePostResponseDTO{
		ID:        postModel.ID,
		Title:     postModel.Title,
		Content:   postModel.Content,
		AuthorID:  postModel.AuthorID,
		CreatedAt: postModel.CreatedAt,
		UpdatedAt: postModel.UpdatedAt,
	}
	return res, nil
}
