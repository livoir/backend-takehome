package usecase

import (
	"app/domain"
	"app/pkg/common"
	"context"
)

type CommentUseCaseImpl struct {
	commentRepository domain.CommentRepository
	userRepository    domain.UserRepository
	postRepository    domain.PostRepository
	transactor        domain.Transactor
}

func NewCommentUseCaseImpl(commentRepository domain.CommentRepository, userRepository domain.UserRepository, postRepository domain.PostRepository, transactor domain.Transactor) domain.CommentUsecase {
	return &CommentUseCaseImpl{
		commentRepository: commentRepository,
		userRepository:    userRepository,
		postRepository:    postRepository,
		transactor:        transactor,
	}
}

// CreateComment implements domain.CommentUsecase.
func (uc *CommentUseCaseImpl) CreateComment(ctx context.Context, postID int64, req domain.CreateCommentRequestDTO) (*domain.CreateCommentResponseDTO, error) {
	tx, err := uc.transactor.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	req.Content = common.Sanitize(req.Content)
	user, err := uc.userRepository.FindByID(ctx, req.AuthorID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}
	post, err := uc.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, common.ErrPostNotFound
	}
	comment := &domain.Comment{
		Content:    req.Content,
		PostID:     postID,
		AuthorName: user.Name,
	}
	err = uc.commentRepository.Create(ctx, tx, comment)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	response := &domain.CreateCommentResponseDTO{
		ID:         comment.ID,
		Content:    comment.Content,
		PostID:     comment.PostID,
		AuthorName: comment.AuthorName,
		CreatedAt:  comment.CreatedAt,
	}
	return response, nil
}

// FindCommentsByPostID implements domain.CommentUsecase.
func (uc *CommentUseCaseImpl) FindCommentsByPostID(ctx context.Context, postID int64, param domain.SearchParam) ([]*domain.Comment, int64, error) {
	post, err := uc.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, 0, err
	}
	if post == nil {
		return nil, 0, common.ErrPostNotFound
	}
	comments, total, err := uc.commentRepository.FindByPostID(ctx, postID, param)
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}
