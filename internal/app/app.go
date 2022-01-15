package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	q "github.com/core-go/sql"
	"github.com/core-go/sql/query"
	_ "github.com/go-sql-driver/mysql"
	"reflect"

	_ "github.com/lib/pq"
	. "go-service/internal/usecase/user"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(root.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(root.Status)
	action := sv.InitializeAction(root.Action)
	validator := v.NewValidator()

	userType := reflect.TypeOf(User{})
	userQueryBuilder := query.NewBuilder(db, "usertest", userType)
	userSearchBuilder, err := q.NewSearchBuilder(db, userType, userQueryBuilder.BuildQuery)
	if err != nil {
		return nil, err
	}

	userRepository := NewUserRepository(db)
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
