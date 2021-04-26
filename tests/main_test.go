package tests

import (
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	handler "github.com/MAVIKE/yad-backend/internal/delivery/http"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
)

const (
	userDB     = "postgres"
	hostDB     = "localhost"
	passwordDB = "secret"
	nameDB     = "postgres"
	portDB     = "5433"
	sslmodeDB  = "disable"

	signingKey     = "test"
	accessTokenTTL = 720

	schemaDir = "../schema/"

	adminType   = "admin"
	userType       = "user"
	courierType    = "courier"
	restaurantType = "restaurant"
)

type APITestSuite struct {
	suite.Suite

	db       *sqlx.DB
	repos    *repository.Repository
	services *service.Service
	handlers *handler.Handler

	app *echo.Echo

	tokenManager *auth.Manager

	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) SetupSuite() {
	s.initPool()
	s.initApp()

	if err := s.downDB(); err != nil {
		s.FailNow("Failed to down db schema", err)
	}

	if err := s.initDB(); err != nil {
		s.FailNow("Failed to init db schema", err)
	}
}

func (s *APITestSuite) initPool() {
	var err error
	s.pool, err = dockertest.NewPool("")
	if err != nil {
		s.FailNow("Failed to connect to docker", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.2",
		Env: []string{
			"POSTGRES_USER=" + userDB,
			"POSTGRES_PASSWORD=" + passwordDB,
			"POSTGRES_DB=" + nameDB,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: portDB},
			},
		},
	}

	s.resource, err = s.pool.RunWithOptions(&opts)
	if err != nil {
		s.FailNow("Failed to start resource", err)
	}
}

func (s *APITestSuite) initApp() {
	var err error
	err = s.pool.Retry(func() error {
		s.db, err = repository.NewPostgresDB(repository.Config{
			Host:     hostDB,
			Port:     portDB,
			Username: userDB,
			DBName:   nameDB,
			SSLMode:  sslmodeDB,
			Password: passwordDB,
		})
		return err
	})
	if err != nil {
		s.FailNow("Failed to initialize db", err)
	}

	s.repos = repository.NewRepository(s.db)

	s.tokenManager, err = auth.NewManager(signingKey)
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	deps := service.Deps{
		Repos:          s.repos,
		TokenManager:   s.tokenManager,
		AccessTokenTTL: time.Duration(accessTokenTTL) * time.Hour,
	}

	s.services = service.NewService(deps)
	s.handlers = handler.NewHandler(s.services, s.tokenManager)

	s.app = echo.New()
	s.app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: ioutil.Discard,
	}))
	s.handlers.Init(s.app)
}

func (s *APITestSuite) initDB() error {
	schema, err := ioutil.ReadFile(schemaDir + "init.sql")
	if err != nil {
		return err
	}
	s.db.MustExec(string(schema))
	return err
}

func (s *APITestSuite) downDB() error {
	schema, err := ioutil.ReadFile(schemaDir + "down.sql")
	if err != nil {
		return err
	}
	s.db.MustExec(string(schema))
	return err
}

func (s *APITestSuite) TearDownSuite() {
	if err := s.pool.Purge(s.resource); err != nil {
		s.FailNow("Failed to purge resource", err)
	}
}

func (s *APITestSuite) TestPing() {
	req, err := http.NewRequest("GET", "/api/v1/ping", nil)
	if err != nil {
		s.FailNow("Failed to build request", err)
	}

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	s.Require().Equal(http.StatusOK, resp.Result().StatusCode)
}

// https://pkg.go.dev/github.com/stretchr/testify/suite
// SetupTest or BeforeTest ?
func (s *APITestSuite) SetupTest() {
	schema, err := ioutil.ReadFile(schemaDir + "truncate.sql")
	if err != nil {
		s.FailNow("Failed to truncate db", err)
	}
	s.db.MustExec(string(schema))

	filenames := []string{
		"admin.sql",
		"user.sql",
		"courier.sql",
		"restaurant.sql",
		"category.sql",
		"order.sql",
	}

	for _, filename := range filenames {
		schema, err := ioutil.ReadFile(schemaDir + "test/" + filename)
		if err != nil {
			s.FailNow("Failed to populate db", err)
		}
		s.db.MustExec(string(schema))
	}
}

func (s *APITestSuite) getJWT(clientId int, clientType string) (string, error) {
	return s.tokenManager.NewJWT(clientId, clientType, accessTokenTTL)
}
