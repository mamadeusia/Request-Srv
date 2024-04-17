package main

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/mamadeusia/RequestSrv/client/postgres"
	"github.com/mamadeusia/RequestSrv/config"
	handler "github.com/mamadeusia/RequestSrv/handler/request"
	pb "github.com/mamadeusia/RequestSrv/proto"
	service "github.com/mamadeusia/RequestSrv/service/application"
	"github.com/mamadeusia/go-micro-plugins/events/natsjs"

	authpb "github.com/mamadeusia/AuthSrv/proto"
	notifpb "github.com/mamadeusia/NotificationSrv/proto"

	pRepo "github.com/mamadeusia/RequestSrv/domain/application/postgres"

	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var (
	serviceName = "requestsrv"
	version     = "latest"
)

func main() {

	// Load configurations
	if err := config.Load(); err != nil {
		logger.Fatal(err)
	}

	// Create service
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	srv.Init(
		micro.Name(serviceName),
		micro.Version(version),
	)

	applicationRepo, err := postgres.NewPostgres(context.Background(), config.PostgresURL())
	if err != nil {
		logger.Fatal(err)
	}

	AuthService := authpb.NewAuthSrvService(config.GetAuthServiceName(), srv.Client())
	NotifService := notifpb.NewNotificationSrvService(config.GetNotificationServiceName(), srv.Client())

	rep := pRepo.NewPostgresRepository(applicationRepo)

	publisherClient, err := natsjs.NewStream(
		natsjs.Address(config.NatsURL()),
		natsjs.NkeyConfig(config.NatsNkey()),
	)
	if err != nil {
		logger.Fatal(err)
	}

	applicationService := service.NewApplicationService(rep, AuthService, NotifService, publisherClient)
	// Register handler
	if err := pb.RegisterRequestSrvHandler(srv.Server(), handler.NewRequestHandler(applicationService, validator.New())); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
