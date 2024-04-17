package application

import (
	authpb "github.com/mamadeusia/AuthSrv/proto"
	notifpb "github.com/mamadeusia/NotificationSrv/proto"
	"go-micro.dev/v4/events"

	"github.com/mamadeusia/RequestSrv/domain/application"
)

type ApplicationService struct {
	persistRequests     application.Repository
	authService         authpb.AuthSrvService
	notificationService notifpb.NotificationSrvService
	natsjs              events.Stream
}

func NewApplicationService(persistRequests application.Repository, authService authpb.AuthSrvService,
	notificationService notifpb.NotificationSrvService, natsjs events.Stream) *ApplicationService {
	return &ApplicationService{
		persistRequests:     persistRequests,
		authService:         authService,
		notificationService: notificationService,
		natsjs:              natsjs,
	}
}
