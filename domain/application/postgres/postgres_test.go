package postgres

import (
	"time"

	"github.com/mamadeusia/RequestSrv/entity"
)

func (its *IntTestSuite) TestPostgresRepository_CreateRequest_Success() {
	// given
	testVals := struct {
		entity.Request
	}{
		entity.Request{
			// ID:              [16]byte{},
			RequesterID:     0,
			FullName:        "",
			Age:             0,
			LocationLat:     0,
			LocationLon:     0,
			Photo:           entity.StoredMessage{},
			Msgs:            []entity.StoredMessage{},
			QuestionAnswers: entity.QuestionAnswerSlice{},
			Status:          0,
			CreatedAt:       time.Time{},
		},
	}

	// when

	err := its.pgRepo.CreateRequest(its.ctx, testVals.Request)

	// then
	its.Nil(err)
}
