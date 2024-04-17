package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status uint32

const (
	Unknown Status = iota
	//Filling the Request by Requester
	StatusFilling
	//Filled
	StatusFilled
	//Uncompleted
	StatusUncompleted
	//Wait to watch by admin
	StatusAdminPending
	//Admin what the Request data need to approve
	StatusAdminChecking
	//maybe we need admin Ask state and a state for requester response state
	StatusAdminAsk
	//it can be a loop
	StatusRequesterAdminResponse
	//in this state data has been sent to the validators and wait for receive questions
	StatusValidatorsChecking
	//in this state question send to the user and user must answer to them
	StatusRequesterValidatorResponse
	//in this state all validator questions answered and now user must wait for voting process on blockchain
	StatusBlockchainPending
	//this state request has been approved or rejected
	StatusRejected

	StatusApproved

	StatusValidatorQuestionReady

	//error status
	StatusEmpty
)

// Request - is responsible for demonstrating a request from requester
type Request struct {
	ID          uuid.UUID
	RequesterID int64
	FullName    string
	//FullNameAnswerTime time.Time statistical data
	Age int32
	// AgeAnswerTime time.Time
	LocationLat float64
	LocationLon float64
	// LocationAnswerTime time.Time statistical data

	//the messages forwarded to private chat , we need to access it .
	Photo           StoredMessage       `json:"photo"`
	Msgs            []StoredMessage     `json:"msgs"` //TODO : should be tested of that doesn't work return to the initial phase
	QuestionAnswers QuestionAnswerSlice `json:"question_answers"`
	Status          Status
	CreatedAt       time.Time
}

func (dm StoredMessage) Value() (driver.Value, error) {
	return json.Marshal(dm)
}
func (dm *StoredMessage) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &dm)

}

type QuestionAnswerSlice struct {
	QuestionAnswers []QuestionAnswer
}

func (dm QuestionAnswerSlice) Value() (driver.Value, error) {
	return json.Marshal(dm)
}

func (dm *QuestionAnswerSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &dm)

}

type StoredMessage struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int32 `json:"message_id"`
	// AnswerTime time.Time `json:"answer_time"`
}

// QuestionAnswer is responsible for storing admin questions and requester answer
type QuestionAnswer struct {
	Question string
	Answer   []StoredMessage
}
