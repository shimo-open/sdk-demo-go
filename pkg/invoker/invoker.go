package invoker

import (
	"log"
	"net/url"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/elog"

	"sdk-demo-go/pkg/models/db"
	"sdk-demo-go/pkg/services"
	"sdk-demo-go/ui"

	"github.com/ego-component/egorm"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/server/egin"
	sdk "github.com/shimo-open/sdk-kit-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// Gin is the HTTP server component
	Gin *egin.Component
	// Services holds all business service instances
	Services *services.Services
	// DB is the database connection instance
	DB *gorm.DB
	// SdkMgr is the Shimo SDK manager instance
	SdkMgr *sdk.Manager
)

func Init() error {
	var err error
	Gin = egin.Load("server.http").Build(egin.WithEmbedFs(ui.WebUI))
	Services = services.NewServices()
	if econf.GetBool("mysql.use") {
		DB = egorm.Load("mysql").Build()
	} else {
		// Connect to an in-memory database so each run shares the same connection
		dsn := "file::memory:?cache=shared" // Use an in-memory database with a shared cache
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		// Configure the connection pool
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}

		// Set the maximum number of open connections
		sqlDB.SetMaxOpenConns(100)
		// Set the maximum number of idle connections
		sqlDB.SetMaxIdleConns(50)
		// Set the maximum connection lifetime
		sqlDB.SetConnMaxLifetime(0)
		// Create the tables
		if err = initTables(); err != nil {
			return err
		}
	}
	InitShimo()
	return nil
}

func initTables() error {
	models := []interface{}{
		&db.Team{}, // Base table
		&db.User{}, // Base table

		&db.AppClient{},  // Depends on users
		&db.Department{}, // Depends on teams
		&db.TeamRole{},   // Depends on teams and users

		&db.DeptMember{}, // Depends on departments and users
		&db.File{},       // File table

		&db.Event{},           // Depends on files
		&db.FilePermissions{}, // Depends on files and users

		&db.KnowledgeBase{}, // Knowledge base table
		&db.TestApi{},       // Standalone table
	}
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Panicf("failed to migrate tables: %v", err)
		return err
	}

	return nil
}

func InitShimo() {
	// Initialize the sdk-sdk service
	shimoHost := econf.GetString("shimoSDK.host")
	_, err := url.Parse(shimoHost)
	if err != nil {
		elog.Panic("invalid endpoint", l.E(err))
	}
	SdkMgr = sdk.NewManager(
		sdk.WithAppID(econf.GetString("shimoSDK.appId")),
		sdk.WithAppSecret(econf.GetString("shimoSDK.appSecret")),
		sdk.WithHTTPClient(ehttp.Load("").Build(ehttp.WithAddr(shimoHost), ehttp.WithRawDebug(true))),
	)
	Services = services.NewServices()
}
