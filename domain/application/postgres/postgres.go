package postgres

import (
	"context"
	"errors"
	"fmt"

	postgresClient "github.com/mamadeusia/RequestSrv/client/postgres"
	"github.com/mamadeusia/RequestSrv/entity"

	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
)

type PostgresRepository struct {
	Postgres *postgresClient.PostgresClient
}

func NewPostgresRepository(postgres *postgresClient.PostgresClient) *PostgresRepository {
	return &PostgresRepository{
		Postgres: postgres,
	}
}

// CreateRequest - responsible for calling sqlc
func (pr *PostgresRepository) CreateRequest(ctx context.Context, request entity.Request) error {

	// create in db just if the status is filled
	// it's a duplicate we have it in handler too!!
	if request.Status != entity.StatusFilled {
		logger.Info("DOMAIN::CreateRequest, has failed with error : request status is not statusFilled!")
		return errors.New("request status is not statusFilled")
	}
	params := postgresClient.CreateRequestParams{
		ID:          request.ID,
		RequesterID: request.RequesterID,
		FullName:    request.FullName,
		Age:         request.Age,
		LocationLat: request.LocationLat,
		LocationLon: request.LocationLon,
		Status:      postgresClient.RequestStatusFilled,
		Photo:       request.Photo,
		Msgs:        request.Msgs,
	}
	if _, err := pr.Postgres.Queries.CreateRequest(ctx, params); err != nil {
		return err
	}

	return nil

}

func (pr *PostgresRepository) GetRequestByID(ctx context.Context, id string) (*entity.Request, error) {
	result, err := pr.Postgres.Queries.GetRequestByID(ctx, uuid.Must(uuid.Parse(id)))
	if err != nil {
		logger.Info("DOMAIN::GetRequestByID, has failed with error , %v,id: %v", err, id)
		return nil, err
	}
	return &entity.Request{
		ID:          result.ID,
		RequesterID: result.RequesterID,
		FullName:    result.FullName,
		Age:         result.Age,
		LocationLat: result.LocationLat,
		LocationLon: result.LocationLon,
		Status:      ConvertPostgresStatusToEntityStatus(result.Status),
		CreatedAt:   result.CreatedAt,
	}, nil
}

// GetRequestByAdminID - responsible for calling sqlc
// TODO remove for each for send request to databases
func (pr *PostgresRepository) GetRequestByAdminID(ctx context.Context, adminID int64, limit, offset int32) ([]entity.Request, error) {

	requests, err := pr.Postgres.Queries.GetRequestCollaboratorsByAdminID(ctx, postgresClient.GetRequestCollaboratorsByAdminIDParams{
		AdminID: adminID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		logger.Info("DOMAIN::GetRequestByAdminID, has failed with error , %v,adminID: %v", err, adminID)
		return nil, err
	}

	var result []entity.Request
	for _, req := range requests {
		res, err := pr.Postgres.Queries.GetRequestByID(ctx, req.RequestID)
		if err != nil {
			logger.Info("DOMAIN::GetRequestByAdminID, has failed with error , %v, requestID: %v", err, req.RequestID)
			return nil, err
		}
		var msgs []entity.StoredMessage

		for _, msg := range res.Msgs {
			msgs = append(msgs, entity.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}
		result = append(result, entity.Request{
			ID:          res.ID,
			RequesterID: res.RequesterID,
			FullName:    res.FullName,
			Age:         res.Age,
			LocationLat: res.LocationLat,
			LocationLon: res.LocationLon,
			Photo: entity.StoredMessage{
				ChatID:    res.Photo.ChatID,
				MessageID: res.Photo.MessageID,
			},
			Msgs:      msgs,
			Status:    ConvertPostgresStatusToEntityStatus(res.Status),
			CreatedAt: res.CreatedAt,
		})
	}

	return result, nil
}

func (pr *PostgresRepository) CountRequestByAdminID(ctx context.Context, adminID int64) (int64, error) {

	requestCount, err := pr.Postgres.Queries.GetCountRequestCollaboratorsByAdminID(ctx, adminID)
	if err != nil {
		logger.Info("DOMAIN::GetRequestByAdminID, has failed with error , %v,adminID: %v", err, adminID)
		return 0, err
	}

	return requestCount, nil
}
func (pr *PostgresRepository) GetOrphanRequests(ctx context.Context, limit, offset int32) ([]entity.Request, error) {

	result, err := pr.Postgres.Queries.GetOrphanRequests(ctx, postgresClient.GetOrphanRequestsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	var output []entity.Request
	for _, pgRequest := range result {

		var msgs []entity.StoredMessage
		for _, msg := range pgRequest.Msgs {
			msgs = append(msgs, entity.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}
		output = append(output, entity.Request{
			ID:          pgRequest.ID,
			RequesterID: pgRequest.RequesterID,
			FullName:    pgRequest.FullName,
			Age:         pgRequest.Age,
			LocationLat: pgRequest.LocationLat,
			LocationLon: pgRequest.LocationLon,
			Photo: entity.StoredMessage{
				ChatID:    pgRequest.Photo.ChatID,
				MessageID: pgRequest.Photo.MessageID,
			},
			Msgs:      msgs,
			Status:    ConvertPostgresStatusToEntityStatus(pgRequest.Status),
			CreatedAt: pgRequest.CreatedAt,
		})
	}
	return output, nil
}

// GetRequestByRequesterID - is responsible for calling sqlc
func (pr *PostgresRepository) GetRequestByRequesterID(ctx context.Context, requesterID int64, limit, offset int32) ([]entity.Request, error) {
	result, err := pr.Postgres.Queries.GetRequestByRequesterID(ctx, postgresClient.GetRequestByRequesterIDParams{
		RequesterID: requesterID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		logger.Info("DOMAIN::GetRequestByRequesterID, has failed with error : %v", err)
		return nil, err
	}

	var requests []entity.Request

	for _, res := range result {
		var msgs []entity.StoredMessage
		for _, msg := range res.Msgs {
			msgs = append(msgs, entity.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}

		var qas []entity.QuestionAnswer
		for _, qa := range res.QuestionAnswers.QuestionAnswers {
			qas = append(qas, entity.QuestionAnswer{
				Question: qa.Question,
				Answer:   qa.Answer,
			})
		}

		requests = append(requests, entity.Request{
			ID:          res.ID,
			RequesterID: res.RequesterID,
			FullName:    res.FullName,
			Age:         res.Age,
			LocationLat: res.LocationLat,
			LocationLon: res.LocationLon,
			Photo: entity.StoredMessage{
				ChatID:    res.Photo.ChatID,
				MessageID: res.Photo.MessageID,
			},
			Msgs: msgs,
			QuestionAnswers: entity.QuestionAnswerSlice{
				QuestionAnswers: qas,
			},
			Status:    ConvertPostgresStatusToEntityStatus(res.Status),
			CreatedAt: res.CreatedAt,
		})

	}

	return requests, nil
}

// UpdateRequestStatus - is responsible for calling sqlc
func (pr *PostgresRepository) UpdateRequestStatus(ctx context.Context, requestID string, status entity.Status) error {
	if err := pr.Postgres.Queries.UpdateRequestStatus(ctx, postgresClient.UpdateRequestStatusParams{
		ID:     uuid.Must(uuid.Parse(requestID)),
		Status: postgresClient.RequestStatus(fmt.Sprint(status)),
	}); err != nil {
		logger.Info("DOMAIN::UpdateRequestStatus , has failed with error , %v", err)
		return err
	}
	return nil
}

// UpdateRequestCollaboratorsAdminID - is responsible for calling sqlc
func (pr *PostgresRepository) UpdateRequestCollaboratorsAdminID(ctx context.Context, requestID string, adminID int64) error {
	if err := pr.Postgres.Queries.UpdateRequestColloboratorsAdmin(ctx, postgresClient.UpdateRequestColloboratorsAdminParams{
		RequestID: uuid.Must(uuid.Parse(requestID)),
		AdminID:   adminID,
	}); err != nil {
		logger.Info("DOMAIN::UpdateRequestCollaboratorsAdminID, has failed with error %v", err)
		return err
	}
	return nil
}

// UpdateRequestCollaboratorsValidators - is responsible for calling sqlc
func (pr *PostgresRepository) UpdateRequestCollaboratorsValidators(ctx context.Context, requestID string, validators []int64) error {
	if err := pr.Postgres.Queries.UpdateRequestColloboratorsValidators(ctx, postgresClient.UpdateRequestColloboratorsValidatorsParams{
		RequestID:  uuid.Must(uuid.Parse(requestID)),
		Validators: validators,
	}); err != nil {
		logger.Info("DOMAIN::UpdateRequestCollaboratorsValidators, has failed with error %v", err)
		return err
	}
	return nil
}

// UpdateRequesterQuestionANDAnswers- is responsible for calling sqlc
func (pr *PostgresRepository) UpdateRequesterQuestionANDAnswers(ctx context.Context, requestID string, questionAnswers entity.QuestionAnswerSlice) error {
	if err := pr.Postgres.Queries.UpdateRequestQuestionAnswers(ctx, postgresClient.UpdateRequestQuestionAnswersParams{
		ID:              uuid.Must(uuid.Parse(requestID)),
		QuestionAnswers: questionAnswers,
	}); err != nil {
		logger.Info("DOMAIN::UpdateRequesterQuestionANDAnswers, has failed with error %v", err)
		return err
	}
	return nil
}

func (pr *PostgresRepository) GetRequestQuestionAnswer(ctx context.Context, requestID string) (*entity.QuestionAnswerSlice, error) {
	result, err := pr.Postgres.Queries.GetRequestQuestionAnswers(ctx, uuid.Must(uuid.Parse(requestID)))
	if err != nil {
		logger.Info("DOMAIN::GetRequestQuestionAnswer, has failed with error %v", err)
		return nil, err
	}

	return &result, nil
}

func ConvertPostgresStatusToEntityStatus(sIn postgresClient.RequestStatus) entity.Status {

	switch sIn {
	case postgresClient.RequestStatusFilling:
		return entity.StatusFilling
	case postgresClient.RequestStatusFilled:
		return entity.StatusFilled
	case postgresClient.RequestStatusUncompleted:
		return entity.StatusUncompleted
	case postgresClient.RequestStatusAdminPending:
		return entity.StatusAdminPending
	case postgresClient.RequestStatusAdminChecking:
		return entity.StatusAdminChecking
	case postgresClient.RequestStatusAdminAsk:
		return entity.StatusAdminAsk
	case postgresClient.RequestStatusRequesterAdminResponse:
		return entity.StatusRequesterAdminResponse
	case postgresClient.RequestStatusValidatorsChecking:
		return entity.StatusValidatorsChecking
	case postgresClient.RequestStatusRequesterValidatorResponse:
		return entity.StatusRequesterValidatorResponse
	case postgresClient.RequestStatusBlockchainPending:
		return entity.StatusBlockchainPending
	case postgresClient.RequestStatusAppoved:
		return entity.StatusApproved
	case postgresClient.RequestStatusRejected:
		return entity.StatusRejected
	case postgresClient.RequestStatusValidatorQuestionReady:
		return entity.StatusValidatorQuestionReady

	}
	return entity.StatusEmpty
}
