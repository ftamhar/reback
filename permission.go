package reback

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var stmtReadPermissionsByRoleNameAndResourceName *sql.Stmt

func setStatementReadPermissionsByRoleName(ctx context.Context) {
	stmtReadPermissionsByRoleNameAndResourceName, _ = dbConn.PrepareContext(ctx, `SELECT is_create, is_read, is_update, is_delete
	FROM permissions p
	LEFT JOIN roles r on p.role_id = r.id
	WHERE r.name = $1 AND p.resource = $2`)
}

func ReadPermissionsByRoleNameAndResourceName(ctx context.Context, roleName, resourceName string) (Permission, error) {
	if isConnectToRedis() {
		s, err := redisConn.Get(ctx, fmt.Sprintf("ReadPermissionsByRoleNameAndResourceName.%s.%s", roleName, resourceName)).Result()
		if err == nil {
			var permission Permission
			err := json.Unmarshal([]byte(s), &permission)
			if err != nil {
				return permission, errors.WithStack(err)
			}
			return permission, nil
		}
	}

	var permission Permission
	err := stmtReadPermissionsByRoleNameAndResourceName.
		QueryRowContext(ctx, roleName, resourceName).
		Scan(&permission.IsCreate, &permission.IsRead, &permission.IsUpdate, &permission.IsDelete)
	if err != nil {
		return permission, errors.WithStack(err)
	}

	if isConnectToRedis() {
		b, err := json.Marshal(permission)
		if err != nil {
			return permission, errors.WithStack(err)
		}
		err = redisConn.Set(ctx, fmt.Sprintf("ReadPermissionsByRoleNameAndResourceName.%s.%s", roleName, resourceName), string(b), time.Duration(timeoutRedis)*time.Minute).Err()
		if err != nil {
			return permission, errors.WithStack(err)
		}
	}

	return permission, nil
}

var stmtReadPermissionsByRoleId *sql.Stmt

func setStatementReadPermissionByRoleId(ctx context.Context) {
	stmtReadPermissionsByRoleId, _ = dbConn.PrepareContext(ctx, "SELECT resource, is_create, is_read, is_update, is_delete FROM permissions WHERE role_id = $1")
}

func ReadPermissionByRoleId(ctx context.Context, roleId string) ([]Permission, error) {
	rows, err := stmtReadPermissionsByRoleId.QueryContext(ctx, roleId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var permission Permission
		err = rows.Scan(&permission.Resource, &permission.IsCreate, &permission.IsRead, &permission.IsUpdate, &permission.IsDelete)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return permissions, nil
}

func CreatePermissions(ctx context.Context, permissions []Permission) error {
	if len(permissions) == 0 {
		return errors.New("permissions is empty")
	}

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Rollback()
	}()

	result, err := tx.ExecContext(ctx, "DELETE FROM permissions WHERE role_id = ", permissions[0].RoleId)
	if err != nil {
		return errors.WithStack(err)
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Insert("permissions").Columns("resource", "role_id", "is_create", "is_read", "is_update", "is_delete", "created_by")

	for _, permission := range permissions {
		query = query.Values(permission.Resource, permission.RoleId, permission.IsCreate, permission.IsRead, permission.IsUpdate, permission.IsDelete, permission.CreatedBy)
	}

	_, err = query.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
