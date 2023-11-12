package reback

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

var dbConn *sql.DB

var redisConn *redis.Client

// in minutes
var timeoutRedis = 5

// SetDBConn set db connection
// this function should be called in main
// this function will panic if db connection is not set
// this function is set all statements that needed
func SetDBConn(ctx context.Context, db *sql.DB) {
	if db == nil {
		panic("db connection is nil")
	}
	dbConn = db
	// set role statements
	setStatementCreateRole(ctx)
	setStatementUpdateRole(ctx)
	setStatementReadAllRole(ctx)
	setStatementReadRoleByID(ctx)
	setStatementHardDeleteRole(ctx)
	setStatementSoftDeleteRole(ctx)

	// set permission statements
	setStatementReadPermissionByRoleId(ctx)
}

func SetRedisConn(redis *redis.Client) {
	if redis == nil {
		panic("redis connection is nil")
	}
	redisConn = redis
}

func checkDBConn() {
	if dbConn == nil {
		panic("db connection not set")
	}
}

func isConnectToRedis() bool {
	return redisConn != nil
}
