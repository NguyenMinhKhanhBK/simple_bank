package api

import (
	"os"
	"testing"
	"time"

	db "github.com/NguyenMinhKhanhBK/simple_bank/db/sqlc"
	"github.com/NguyenMinhKhanhBK/simple_bank/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	assert.NoError(t, err)

	return server
}
