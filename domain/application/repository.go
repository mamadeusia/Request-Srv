package application

import (
	"context"

	"github.com/mamadeusia/RequestSrv/entity"
)

type Repository interface {
	CreateRequest(ctx context.Context, request entity.Request) error
	GetRequestByID(ctx context.Context, id string) (*entity.Request, error)
	GetRequestByAdminID(ctx context.Context, adminID int64, limit, offset int32) ([]entity.Request, error)
	GetOrphanRequests(ctx context.Context, limit, offset int32) ([]entity.Request, error)
	CountRequestByAdminID(ctx context.Context, adminID int64) (int64, error)
	GetRequestByRequesterID(ctx context.Context, requesterID int64, limit, offset int32) ([]entity.Request, error)
	UpdateRequestStatus(ctx context.Context, requestID string, status entity.Status) error
	UpdateRequestCollaboratorsAdminID(ctx context.Context, requestID string, adminID int64) error
	UpdateRequestCollaboratorsValidators(ctx context.Context, requestID string, validators []int64) error
	UpdateRequesterQuestionANDAnswers(ctx context.Context, requestID string, questionAnswers entity.QuestionAnswerSlice) error
	GetRequestQuestionAnswer(ctx context.Context, requestID string) (*entity.QuestionAnswerSlice, error)
}
