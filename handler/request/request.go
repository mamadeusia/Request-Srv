package request

import (
	"context"
	"errors"

	"github.com/mamadeusia/RequestSrv/domain/application"
	"github.com/mamadeusia/RequestSrv/entity"
	RequestSrv "github.com/mamadeusia/RequestSrv/proto"
	pb "github.com/mamadeusia/RequestSrv/proto"

	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
)

// CreateRequest implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) CreateRequest(ctx context.Context, req *pb.CreateRequestRequest, res *pb.CreateRequestResponse) error {
	if req.Age == 0 || req.FullName == "" || req.LocationLat == 0 || req.LocationLon == 0 ||
		req.RequesterID == 0 {
		logger.Info("HANDLER::CreateRequest, has failed with error: invalid req type %+v", req)
		res.Msg = "Failed"
		return errors.New("invalid input type")
	}

	var msgs []entity.StoredMessage

	for _, msg := range req.Msgs {
		msgs = append(msgs, entity.StoredMessage{
			ChatID:    msg.ChatID,
			MessageID: msg.MessageID,
		})
	}

	if req.Status != pb.Status_StatusFilled {
		logger.Info("HANDLER::CreateRequest, has failed with error : the status is not statusFilled!")
		res.Msg = "Failed"
		return errors.New("status is not StatusFilled")
	}

	r.service.CreateApplication(ctx, application.Application{
		Request: &entity.Request{
			ID:          uuid.New(),
			RequesterID: req.RequesterID,
			FullName:    req.FullName,
			Age:         req.Age,
			LocationLat: req.LocationLat,
			LocationLon: req.LocationLat,
			Photo: entity.StoredMessage{
				ChatID:    req.Photo.ChatID,
				MessageID: req.Photo.MessageID,
			},
			Msgs:   msgs,
			Status: entity.StatusFilled,
		},
		AdminId:      0,
		ValidatorIds: []int64{},
	})

	res.Msg = "Success"
	return nil
}

// GetRequestByAdminID implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) GetRequestByAdminID(ctx context.Context, req *pb.GetRequestByAdminIDRequest, res *pb.GetRequestByAdminIDResponse) error {
	if req.AdminID == 0 {
		logger.Info("HANDLER::GetRequestByAdminID, has failed with error : invalid input type : %v", req)
		return errors.New("invalid input type")
	}

	if req.Limit == 0 {
		//put a default value
		req.Limit = 10
	}

	if req.Offset == 0 {
		//put a default value
		req.Offset = 10
	}

	requests, err := r.service.GetRequestByAdminID(ctx, req.AdminID, req.Limit, req.Offset)
	if err != nil {
		logger.Info("HANDLER::GetRequestByAdminID , has failed with error: %v", err)
		return err
	}
	if len(requests) == 0 {
		logger.Info("HANDLER::GetRequestByAdminID , has failed with error: requests are empty")
		return errors.New("requests are empty")
	}

	for _, request := range requests {
		res.Count = int32(len(requests))

		var msgs []*pb.StoredMessage
		for _, msg := range request.Msgs {
			msgs = append(msgs, &pb.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}
		var questionAnswers []*pb.QuestionAnswer

		for _, qa := range request.QuestionAnswers.QuestionAnswers {
			var sm []*pb.StoredMessage
			for _, a := range qa.Answer {
				sm = append(sm, &pb.StoredMessage{
					ChatID:    a.ChatID,
					MessageID: a.MessageID,
				})
			}
			questionAnswers = append(questionAnswers, &pb.QuestionAnswer{
				Question: qa.Question,
				Answers:  sm,
			})
		}
		res.Requests = append(res.Requests, &pb.Request{
			ID:              request.ID.String(),
			RequesterID:     request.RequesterID,
			FullName:        request.FullName,
			Age:             request.Age,
			LocationLat:     request.LocationLat,
			LocationLon:     request.LocationLon,
			Photo:           &pb.StoredMessage{ChatID: request.Photo.ChatID, MessageID: request.Photo.MessageID},
			Msgs:            msgs,
			QuestionAnswers: questionAnswers,
			Status:          pb.Status(request.Status),
		})
	}
	return nil
}

// GetRequestByRequesterID implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) GetRequestByRequesterID(ctx context.Context, req *pb.GetRequestByRequesterIDRequest, res *pb.GetRequestByRequesterIDResponse) error {
	if req.RequesterID == 0 {
		logger.Info("HANDLER::GetRequestByRequesterID, has failed with error : invalid input %v", req)
		return errors.New("invalid input")
	}

	if req.Limit == 0 {
		//put a default value
		req.Limit = 10
	}

	if req.Offset == 0 {
		//put a default value
		req.Offset = 10
	}

	results, err := r.service.GetRequestByRequesterID(ctx, req.RequesterID, req.Limit, req.Offset)
	if err != nil {
		logger.Info("HANDLER::GetRequestByRequesterID, has failed with error : %v", err)
		return err
	}

	if len(results) == 0 {
		logger.Info("HANDLER::GetRequestByAdminID , has failed with error: requests are empty")
		return errors.New("requests are empty")
	}

	for _, result := range results {
		var msgs []*pb.StoredMessage
		for _, msg := range result.Msgs {
			msgs = append(msgs, &pb.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}

		var questionAnswers []*pb.QuestionAnswer

		for _, qa := range result.QuestionAnswers.QuestionAnswers {
			var sm []*pb.StoredMessage
			for _, a := range qa.Answer {
				sm = append(sm, &pb.StoredMessage{
					ChatID:    a.ChatID,
					MessageID: a.MessageID,
				})
			}
			questionAnswers = append(questionAnswers, &pb.QuestionAnswer{
				Question: qa.Question,
				Answers:  sm,
			})
		}
		res.Requests = append(res.Requests, &pb.Request{
			ID:          result.ID.String(),
			RequesterID: result.RequesterID,
			FullName:    result.FullName,
			Age:         result.Age,
			LocationLat: result.LocationLat,
			LocationLon: result.LocationLon,
			Photo: &pb.StoredMessage{
				ChatID:    result.Photo.ChatID,
				MessageID: result.Photo.MessageID,
			},
			Msgs:            msgs,
			QuestionAnswers: questionAnswers,
			Status:          pb.Status(result.Status),
		})
	}

	return nil
}

// UpdateRequest implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) UpdateRequest(ctx context.Context, req *pb.UpdateRequestRequest, res *pb.UpdateRequestResponse) error {
	if req.RequestID == "" {
		logger.Info("HANDLER::UpdateRequest, has failed with error : invalid input : %v", req)
		res.Msg = "Failed"
		return errors.New("invalid input")
	}
	var status entity.Status
	switch req.Status {
	case pb.Status_StatusAdminAsk:
		status = entity.StatusAdminAsk
	case pb.Status_StatusAdminPending:
		status = entity.StatusAdminPending
	case pb.Status_StatusAdminChecking:
		status = entity.StatusAdminChecking
	case pb.Status_StatusBlockchainPending:
		status = entity.StatusBlockchainPending
	case pb.Status_StatusApproved:
		status = entity.StatusApproved
	case pb.Status_StatusRejected:
		status = entity.StatusRejected
	case pb.Status_StatusFilled:
		status = entity.StatusFilled
	case pb.Status_StatusFilling:
		status = entity.StatusFilling
	case pb.Status_StatusRequesterAdminResponse:
		status = entity.StatusRequesterAdminResponse
	case pb.Status_StatusRequesterValidatorResponse:
		status = entity.StatusRequesterValidatorResponse
	case pb.Status_StatusUncompleted:
		status = entity.StatusUncompleted
	case pb.Status_StatusValidatorsChecking:
		status = entity.StatusValidatorsChecking
	case pb.Status_StatusValidatorQuestionReady:
		status = entity.StatusValidatorQuestionReady
	default:
		status = entity.Unknown

	}

	if status == entity.Unknown {
		logger.Info("HANDLER::UpdateRequest, has failed with error: invalid status %v", status)
		res.Msg = "Failed"
		return errors.New("invalid status")
	}

	if err := r.service.UpdateRequest(ctx, req.RequestID, status); err != nil {
		logger.Info("HANDLER::UpdateRequest, has failed with error %v", err)
		res.Msg = "Failed"
		return err
	}

	res.Msg = "Success"

	return nil
}

// SetAdminForRequest implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) SetAdminForRequest(ctx context.Context, req *pb.SetAdminForRequestRequest, res *pb.SetAdminForRequestResponse) error {
	if req.AdminID == 0 || req.RequestID == "" {
		logger.Info("HANDLER::SetAdminForRequest, has failed with error: input type is not valid")
		res.Msg = "Failed"
		return errors.New("input type is not valid")
	}

	if err := r.service.UpdateRequestCollaboratorsAdminID(ctx, req.RequestID, req.AdminID); err != nil {
		logger.Info("HANDLER::SetAdminForRequest, has failed with error %v ", err)
		res.Msg = "Failed"
		return err
	}
	res.Msg = "Success"
	return nil
}

// SetValidatorsForRequest implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) SetValidatorsForRequest(ctx context.Context, req *pb.SetValidatorsForRequestRequest, res *pb.SetValidatorsForRequestResponse) error {
	if req.RequestID == "" {
		logger.Info("HANDLER::SetValidatorsForRequest, has failed with error : input type is not valid")
		res.Msg = "Failed"
		return errors.New("input type is not valid")
	}

	if err := r.service.UpdateRequestCollaboratorsValidators(ctx, req.RequestID); err != nil {
		logger.Info("HANDLER::UpdateRequestCollaboratorsValidators, has failed with error %v", err)
		res.Msg = "Failed"
		return err
	}
	res.Msg = "Success"
	return nil
}

// SendPotentialValidatorsNotifications implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) SendPotentialValidatorsNotifications(ctx context.Context, req *RequestSrv.SendPotentialValidatorsNotificationsRequest, res *RequestSrv.SendPotentialValidatorsNotificationsResponse) error {
	if req.RequestID == "" {
		logger.Info("HANDLER::SendPotentialValidatorsNotifications, has failed with error: invalid input type")
		return errors.New("invalid input type")
	}
	if err := r.service.SendPotentialValidatorsNotifications(ctx, req.RequestID); err != nil {
		logger.Info("HANDLER::SendPotentialValidatorsNotifications, has failed with error : %v", err)
		return err
	}
	return nil
}

// GetRequestByAdminID implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) GetCountRequestByAdminID(ctx context.Context, req *pb.GetCountRequestByAdminIDRequest, res *pb.GetCountRequestByAdminIDResponse) error {
	countRequestByAdminIDRequest := entity.GetCountRequestByAdminIDRequest{
		AdminID: req.AdminID,
	}

	err := r.Validate.Struct(countRequestByAdminIDRequest)
	if err != nil {
		logger.Info("HANDLER::GetRequestByAdminID, has failed with error : invalid input type : %v", req)
		return errors.New("invalid input type")
	}

	requestsCount, err := r.service.CountRequestByAdminID(ctx, req.AdminID)
	if err != nil {
		logger.Info("HANDLER::GetRequestByAdminID , has failed with error: %v", err)
		return err
	}

	//TODO : add pending counter
	res.RequestCounter = requestsCount
	res.PendingCounter = requestsCount

	return nil
}

// rpc GetOrphanRequests(GetOrphanRequestsRequest) returns (GetOrphanRequestsResponse) {}
func (r *RequestHandler) GetOrphanRequests(ctx context.Context, req *pb.GetOrphanRequestsRequest, res *pb.GetOrphanRequestsResponse) error {
	if req.Limit < 0 || req.Offset < 0 {
		logger.Info("HANDLER::GetOrphanRequests, has failed with error: invalid input type")
		return errors.New("invalid input type")
	}
	results, err := r.service.GetOrphanRequests(ctx, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	for _, result := range results {
		var msgs []*pb.StoredMessage
		for _, msg := range result.Msgs {
			msgs = append(msgs, &pb.StoredMessage{
				ChatID:    msg.ChatID,
				MessageID: msg.MessageID,
			})
		}
		res.Requests = append(res.Requests, &pb.Request{
			ID:          result.ID.String(),
			RequesterID: result.RequesterID,
			FullName:    result.FullName,
			Age:         result.Age,
			LocationLat: result.LocationLat,
			LocationLon: result.LocationLon,
			Photo: &pb.StoredMessage{
				ChatID:    result.Photo.ChatID,
				MessageID: result.Photo.MessageID,
			},
			Msgs:   msgs,
			Status: pb.Status(result.Status),
		})
	}
	return nil
}

// AddQuestionAnswer implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) AddQuestionAnswer(ctx context.Context, req *RequestSrv.AddQuestionAnswerRequest, res *RequestSrv.AddQuestionAnswerResponse) error {
	if req.Question == "" {
		logger.Info("HANDLER::AddQuestionAnswer, has failed with error: invalid input type")
		return errors.New("invalid input type")
	}
	var answers []entity.StoredMessage
	for _, a := range req.Answers {
		answers = append(answers, entity.StoredMessage{
			ChatID:    a.ChatID,
			MessageID: a.MessageID,
		})
	}

	if err := r.service.UpdateRequestQuestionAnswer(ctx, req.RequestID, entity.QuestionAnswer{
		Question: req.Question,
		Answer:   answers,
	}); err != nil {
		return err
	}
	return nil
}

// GetQuestionAnswer implements RequestSrv.RequestSrvHandler
func (r *RequestHandler) GetQuestionAnswer(ctx context.Context, req *RequestSrv.GetQuestionAnswerRequest, res *RequestSrv.GetQuestionAnswerResponse) error {
	if req.RequestID == "" {
		logger.Info("HANDLER::GetQuestionAnswer, has failed with error: invalid input type")
		return errors.New("invalid input type")
	}

	result, err := r.service.GetRequestQuestionAnswer(ctx, req.RequestID)
	if err != nil {
		return err
	}

	for _, qa := range result.QuestionAnswers {
		var answers []*pb.StoredMessage
		for _, sm := range qa.Answer {
			answers = append(answers, &pb.StoredMessage{
				ChatID:    sm.ChatID,
				MessageID: sm.MessageID,
			})
		}
		res.Result = append(res.Result, &pb.QuestionAnswer{
			Question: qa.Question,
			Answers:  answers,
		})
	}

	return nil
}
