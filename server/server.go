package server

import (
	"bira.io/template/controller"
	"bira.io/template/httprouter"
	"bira.io/template/infra"
	"bira.io/template/repository"
	"bira.io/template/service"
	"github.com/gin-gonic/gin"
)

func SetupServer(migrationsPath string, sqlDatabaseName string) *gin.Engine {
	router := gin.Default()

	infraCollection := infra.NewInfraCollection(sqlDatabaseName)
	repositoryCollection := repository.NewRepositoryCollection(infraCollection)
	serviceCollection := service.NewServiceCollection(repositoryCollection, infraCollection)
	controllerCollection := controller.NewControllerCollection(serviceCollection)

	httprouter.AddRouters(router, controllerCollection)

	infra.MigrateSqlDatabase(
		migrationsPath,
		infra.BuildSqlConnectionString("cockroach", infra.Config.SqlDbConnectionString, sqlDatabaseName, "sslmode=disable"),
		true,
	)

	return router
}
