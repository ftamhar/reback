package reback

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var (
	ErrNoResultFound = errors.New("no result found")
	ErrNoRoleSet     = errors.New("no role set")
)

func GetPermissionsByResourceNameAndRoleNames(ctx context.Context, resourceName string, roleNames []string) ([]Permission, error) {
	if len(roleNames) == 0 {
		return nil, ErrNoRoleSet
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("is_create, is_read, is_update, is_delete").
		From("permissions").
		LeftJoin("roles ON permissions.role_id = roles.id").
		Where(sq.Eq{
			"permissions.resource": resourceName,
			"roles.name":           roleNames,
		})

	rows, err := query.RunWith(dbConn).QueryContext(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		err = rows.Scan(&permission.IsCreate, &permission.IsRead, &permission.IsUpdate, &permission.IsDelete)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	if len(permissions) == 0 {
		return nil, ErrNoResultFound
	}

	return permissions, nil
}

func GetPermissionsByResourceNameAndRoleIds(ctx context.Context, resourceName string, roleIds []string) ([]Permission, error) {
	if len(roleIds) == 0 {
		return nil, ErrNoRoleSet
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("is_create, is_read, is_update, is_delete").
		From("permissions").
		Where(sq.Eq{
			"resource": resourceName,
			"role_id":  roleIds,
		})

	rows, err := query.RunWith(dbConn).QueryContext(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		err = rows.Scan(&permission.IsCreate, &permission.IsRead, &permission.IsUpdate, &permission.IsDelete)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	if len(permissions) == 0 {
		return nil, ErrNoResultFound
	}

	return permissions, nil
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
