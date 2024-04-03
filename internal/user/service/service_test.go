package usersvc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/arfan21/mertani/internal/entity"
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/internal/user"
	userrepo "github.com/arfan21/mertani/internal/user/repository"
	"github.com/arfan21/mertani/pkg/constant"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	pgxMock     pgxmock.PgxPoolIface
	redisClient *redis.Client
	redisMock   redismock.ClientMock
	repoPG      user.Repository
	servcie     user.Service

	defaultPassword       = "test123qwe"
	defaultHashedPassword = "$2a$10$RKU1hsAXRPXvf2tPXdyGnuzfM.gikV04zp.D7cwWG0dwEGD519B9m"
)

func initDep(t *testing.T) {
	if pgxMock == nil {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)

		pgxMock = mockPool
	}

	if redisClient == nil || redisMock == nil {
		client, clientMock := redismock.NewClientMock()
		redisClient = client
		redisMock = clientMock
	}

	if repoPG == nil {
		repoPG = userrepo.New(pgxMock)
	}

	if servcie == nil {
		servcie = New(repoPG)
	}
}

func TestRegister(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Fullname: "test",
			Email:    "test@gmail.com",
			Password: "test123qwe",
		}

		pgxMock.ExpectExec("INSERT INTO users (.+)").
			WithArgs(req.Fullname, req.Email, pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := servcie.Register(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("failed, email already registered", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Fullname: "test",
			Email:    "test@gmail.com",
			Password: "test123qwe",
		}

		pgxMock.ExpectExec("INSERT INTO users (.+)").
			WithArgs(req.Fullname, req.Email, pgxmock.AnyArg()).
			WillReturnError(&pgconn.PgError{Code: "23505"}) // unique violation

		err := servcie.Register(context.Background(), req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrEmailAlreadyRegistered)
	})

	t.Run("failed, invalid request", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Fullname: "test",
			Email:    "test",
			Password: "test",
		}

		err := servcie.Register(context.Background(), req)

		assert.Error(t, err)

		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)
	})
}

func TestLogin(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.UserLoginRequest{
			Email:    "test@gmail.com",
			Password: "test123qwe",
		}

		id := uuid.New()

		pgxMock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(req.Email).
			WillReturnRows(pgxMock.NewRows([]string{"id", "fullname", "email", "password"}).
				AddRow(id, "test", req.Email, defaultHashedPassword))

		redisPayload := entity.UserRefreshToken{
			ID:    id,
			Email: req.Email,
		}

		redisPayloadJson, err := json.Marshal(redisPayload)
		assert.NoError(t, err)

		redisMock.Regexp().ExpectSet(constant.RedisRefreshTokenKeyPrefix+"asd", string(redisPayloadJson), 0).SetVal("OK")

		res, err := servcie.Login(context.Background(), req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("failed, email not found", func(t *testing.T) {
		req := model.UserLoginRequest{
			Email:    "test1@gmail.com",
			Password: "test123qwe",
		}

		pgxMock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(req.Email).
			WillReturnError(pgx.ErrNoRows)

		_, err := servcie.Login(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrEmailOrPasswordInvalid)
	})
}
