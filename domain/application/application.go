package application

import "github.com/mamadeusia/RequestSrv/entity"

// maybe in future we need some additional fields for the request .
// we called this application maybe in future we plan to add some additional data
// like whose the validators or sth like that
type Application struct {
	Request      *entity.Request // rootID
	AdminId      int64
	ValidatorIds []int64
}
