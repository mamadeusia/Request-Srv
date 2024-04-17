package postgres

import (
	"context"
	"testing"

	"github.com/mamadeusia/RequestSrv/client/postgres"
	pgclient "github.com/mamadeusia/RequestSrv/client/postgres"
	"github.com/mamadeusia/RequestSrv/config"
	"github.com/stretchr/testify/suite"
	"go-micro.dev/v4/logger"
)

type IntTestSuite struct {
	suite.Suite

	pgClient *pgclient.PostgresClient
	pgRepo   *PostgresRepository
	ctx      context.Context
}

// connection to db
// this is fully necessary for the suite to run.
func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	//should connect to redis
	if err := config.Load(); err != nil {
		its.FailNow("redis config load failed", err)
	}

	its.ctx = context.Background()

	pgClient, err := postgres.NewPostgres(its.ctx, config.PostgresURL())
	if err != nil {
		logger.Fatal(err)
	}
	its.pgClient = pgClient

	its.pgRepo = NewPostgresRepository(its.pgClient)

}

func (its *IntTestSuite) TearDownSuite() {
	//should clean up connection to postgres
	its.pgClient.DB.Close()

}

func (its *IntTestSuite) BeforeTest(suiteName, testName string) {
	// switch testName {
	// case "TestRedisRepository_ReserveItem_Success":
	// 	if err := its.redisRepo.CreateCacheBundle(
	// 		its.ctx,
	// 		1231231,
	// 		default_source_key,
	// 		2,
	// 		0,
	// 		0,
	// 		time.Now().Add(1*time.Hour),
	// 		[]string{"SH_1"},
	// 	); err != nil {
	// 		its.FailNow("before: TestRedisRepository_ReserveItem_Success failed")
	// 	}
	// }
	// should put data in redis based on test name .

}
func (its *IntTestSuite) TearDownTest() {
	// should clean up the data that we wrote in redis\
	// empty for now.
	// its.pgClient.

}
