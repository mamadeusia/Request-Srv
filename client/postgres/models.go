// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package postgres

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	entity "github.com/mamadeusia/RequestSrv/entity"
)

type RequestStatus string

const (
	RequestStatusFilling                    RequestStatus = "filling"
	RequestStatusFilled                     RequestStatus = "filled"
	RequestStatusUncompleted                RequestStatus = "uncompleted"
	RequestStatusAdminPending               RequestStatus = "admin_pending"
	RequestStatusAdminChecking              RequestStatus = "admin_checking"
	RequestStatusAdminAsk                   RequestStatus = "admin_ask"
	RequestStatusRequesterAdminResponse     RequestStatus = "requester_admin_response"
	RequestStatusValidatorsChecking         RequestStatus = "validators_checking"
	RequestStatusRequesterValidatorResponse RequestStatus = "requester_validator_response"
	RequestStatusBlockchainPending          RequestStatus = "blockchain_pending"
	RequestStatusAppoved                    RequestStatus = "appoved"
	RequestStatusRejected                   RequestStatus = "rejected"
	RequestStatusValidatorQuestionReady     RequestStatus = "ValidatorQuestionReady"
)

func (e *RequestStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RequestStatus(s)
	case string:
		*e = RequestStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for RequestStatus: %T", src)
	}
	return nil
}

type NullRequestStatus struct {
	RequestStatus RequestStatus
	Valid         bool // Valid is true if RequestStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRequestStatus) Scan(value interface{}) error {
	if value == nil {
		ns.RequestStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RequestStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRequestStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.RequestStatus, nil
}

type RequestCollaborators struct {
	RequestID   uuid.UUID `db:"request_id" json:"request_id"`
	RequesterID int64     `db:"requester_id" json:"requester_id"`
	AdminID     int64     `db:"admin_id" json:"admin_id"`
	Validators  []int64   `db:"validators" json:"validators"`
}

type Requests struct {
	ID              uuid.UUID                  `db:"id" json:"id"`
	FullName        string                     `db:"full_name" json:"full_name"`
	Age             int32                      `db:"age" json:"age"`
	LocationLat     float64                    `db:"location_lat" json:"location_lat"`
	LocationLon     float64                    `db:"location_lon" json:"location_lon"`
	Status          RequestStatus              `db:"status" json:"status"`
	RequesterID     int64                      `db:"requester_id" json:"requester_id"`
	Photo           entity.StoredMessage       `db:"photo" json:"photo"`
	Msgs            []entity.StoredMessage     `db:"msgs" json:"msgs"`
	QuestionAnswers entity.QuestionAnswerSlice `db:"question_answers" json:"question_answers"`
	CreatedAt       time.Time                  `db:"created_at" json:"created_at"`
}
