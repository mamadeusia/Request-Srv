package application

import (
	"context"

	AuthSrv "github.com/mamadeusia/AuthSrv/proto"
	NotificationSrv "github.com/mamadeusia/NotificationSrv/proto"
	"github.com/mamadeusia/RequestSrv/config"
	"github.com/mamadeusia/RequestSrv/domain/application"
	"github.com/mamadeusia/RequestSrv/entity"

	"go-micro.dev/v4/logger"
)

// CreateApplication - is responsible for business logic around creation of application aggregator and routing to repository
func (as *ApplicationService) CreateApplication(ctx context.Context, application application.Application) error {
	if err := as.persistRequests.CreateRequest(ctx, *application.Request); err != nil {
		logger.Info("SERVICE::CreateRequest, has failed with error , %v", err)
		return err
	}
	return nil
}

// GetRequestByAdminID - is responsible for business logic around fetching request by adminID and routing to repository
func (as *ApplicationService) GetRequestByAdminID(ctx context.Context, adminID int64, limit, offset int32) ([]entity.Request, error) {
	requests, err := as.persistRequests.GetRequestByAdminID(ctx, adminID, limit, offset)
	if err != nil {
		logger.Info("SERVICE::GetRequestByAdminID, has failed with error , %v", err)
		return nil, err
	}
	return requests, nil
}

// GetRequestByAdminID - is responsible for business logic around fetching request by adminID and routing to repository
func (as *ApplicationService) GetOrphanRequests(ctx context.Context, limit, offset int32) ([]entity.Request, error) {
	requests, err := as.persistRequests.GetOrphanRequests(ctx, limit, offset)
	if err != nil {
		logger.Info("SERVICE::GetRequestByAdminID, has failed with error , %v", err)
		return nil, err
	}
	return requests, nil
}

func (as *ApplicationService) CountRequestByAdminID(ctx context.Context, adminID int64) (int64, error) {
	count, err := as.persistRequests.CountRequestByAdminID(ctx, adminID)
	if err != nil {
		logger.Info("SERVICE::CountRequestByAdminID, has failed with error , %v", err)
		return 0, err
	}
	return count, nil

}

// GetRequestByRequesterID - is responsible for business logic around fetching request by requesterID and routing to repository
func (as *ApplicationService) GetRequestByRequesterID(ctx context.Context, requesterID int64, limit, offset int32) ([]entity.Request, error) {
	results, err := as.persistRequests.GetRequestByRequesterID(ctx, requesterID, limit, offset)
	if err != nil {
		logger.Info("SERVICE::GetRequestByRequesterID, has failed with error : %v", err)
		return nil, err
	}

	return results, nil
}

// UpdateRequest - is responsible for business logic around updating request's status and routing to repository
func (as *ApplicationService) UpdateRequest(ctx context.Context, requestID string, status entity.Status) error {
	if err := as.persistRequests.UpdateRequestStatus(ctx, requestID, status); err != nil {
		logger.Info("SERVICE::UpdateRequest, has failed with error %v", err)
		return err
	}
	return nil
}

// UpdateRequestCollaboratorsAdminID - is responsible for logic around updating requests colloborators admin and routing to repository
func (as *ApplicationService) UpdateRequestCollaboratorsAdminID(ctx context.Context, requestID string, adminID int64) error {

	if err := as.persistRequests.UpdateRequestCollaboratorsAdminID(ctx, requestID, adminID); err != nil {
		logger.Info("SERVICE::UpdateRequestCollaboratorsAdminID, has failed with error %v", err)
		return err
	}
	return nil
}

// UpdateRequestCollaboratorsValidators - is responsible for logic around updating requests collaborators validators and routing to respository
func (as *ApplicationService) UpdateRequestCollaboratorsValidators(ctx context.Context, requestID string) error {
	//TODO : should fetch data from notification service

	validators := make([]int64, 10)

	if err := as.persistRequests.UpdateRequestCollaboratorsValidators(ctx, requestID, validators); err != nil {
		logger.Info("SERVICE::UpdateRequestCollaboratorsValidators, has failed with error %v", err)
		return err
	}

	return nil
}

// SendPotentialValidatorsNotifications - is responsible for logic around sending notifications to potential validators
func (as *ApplicationService) SendPotentialValidatorsNotifications(ctx context.Context, requestId string) error {

	request, err := as.persistRequests.GetRequestByID(ctx, requestId)
	if err != nil {
		logger.Info("SERVICE::SendPotentialValidatorsNotifications, has failed with error : %v", err)
		return err
	}
	//1. call auth and get the potential validators
	result, err := as.authService.GetNearValidators(ctx, &AuthSrv.GetNearValidatorsRequest{
		LocationLat: request.LocationLat,
		LocationLon: request.LocationLon,
		Distance:    10, //should be in the config or env ??or what
		Limit:       1000,
		Offset:      0,
	})
	if err != nil {
		logger.Info("SERVICE::SendPotentialValidatorsNotifications, has failed with error : %v", err)
		return err
	}

	var bulkValidatorRequest NotificationSrv.CreateBulkValidatorNotificationRequest

	for _, validator := range result.Validators {
		bulkValidatorRequest.CreateValidatorNotificationRequest = append(bulkValidatorRequest.CreateValidatorNotificationRequest, &NotificationSrv.CreateValidatorNotificationRequest{
			From: 0,
			To:   validator,
			MessageOneof: &NotificationSrv.CreateValidatorNotificationRequest_NearRequestFoundDetails{
				NearRequestFoundDetails: &NotificationSrv.NearRequestFoundDetails{
					RequestID: requestId,
					FullName:  request.FullName,
				},
			},
		})

	}
	//2. call the notification and send the notification to all of them
	if _, err := as.notificationService.CreateBulkValidatorNotification(ctx, &bulkValidatorRequest); err != nil {
		logger.Info("SERVICE::SendPotentialValidatorsNotifications, has failed with error : %v", err)
		return err
	}

	//3. publish a message to broker to check if after 5 hours we have enough questions to choose from and if that was
	//ok we will send notification to requester
	//todo:: add minimum required question to something like config or what???
	if err := as.natsjs.Publish(config.NatsCallbackFinalTopic(), &entity.QuestionConfirmedReq{RequestID: requestId, MinRequiredQuestions: 1000, RequesterID: request.RequesterID}); err != nil {
		logger.Info("SERVICE::SendPotentialValidatorsNotifications, has failed with error: %v", err)
		return err
	}

	return nil
}

// UpdateRequestQuestionAnswer - is responsible for logic around question and answers between admin and requester
func (as *ApplicationService) UpdateRequestQuestionAnswer(ctx context.Context, requestID string, questionAnswers entity.QuestionAnswer) error {
	result, err := as.persistRequests.GetRequestQuestionAnswer(ctx, requestID)
	if err != nil {
		logger.Info("SERVICE::UpdateRequestQuestionAnswer, has failed with error: %v", err)
		return err
	}
	result.QuestionAnswers = append(result.QuestionAnswers, questionAnswers)
	if err := as.persistRequests.UpdateRequesterQuestionANDAnswers(ctx, requestID, *result); err != nil {
		logger.Info("SERVICE::UpdateRequestQuestionAnswer, has failed with error: %v", err)
		return err
	}
	return nil
}

// GetRequestQuestionAnswer - is responsible for logic around question and answer fetching
func (as *ApplicationService) GetRequestQuestionAnswer(ctx context.Context, requestID string) (*entity.QuestionAnswerSlice, error) {
	res, err := as.persistRequests.GetRequestQuestionAnswer(ctx, requestID)
	if err != nil {
		logger.Info("SERVICE::GetRequestQuestionAnswer, has failed with error: %v", err)
		return nil, err
	}
	return res, nil
}
