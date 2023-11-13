package reback

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

var ErrNoRowsAffected = errors.New("no rows affected")

var stmtCreateRole *sql.Stmt

func setStatementCreateRole(ctx context.Context) {
	var err error
	stmtCreateRole, err = dbConn.PrepareContext(ctx, "INSERT INTO roles(name, description) VALUES($1, $2) RETURNING id")
	if err != nil {
		panic(err)
	}
}

func CreateRole(ctx context.Context, name, description string) (string, error) {
	var id string
	err := stmtCreateRole.QueryRowContext(ctx, name, description).Scan(&id)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return id, nil
}

var stmtUpdateRole *sql.Stmt

func setStatementUpdateRole(ctx context.Context) {
	stmtUpdateRole, _ = dbConn.PrepareContext(ctx, "UPDATE roles SET name = $1, description = $2, updated_at = now(), updated_by = $3 WHERE id = $4")
}

func UpdateRole(ctx context.Context, id, name, description, updatedBy string) error {
	status, err := stmtUpdateRole.ExecContext(ctx, name, description, updatedBy, id)
	if err != nil {
		return errors.WithStack(err)
	}

	RowsAffected, err := status.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

var stmtHardDeleteRole *sql.Stmt

func setStatementHardDeleteRole(ctx context.Context) {
	stmtHardDeleteRole, _ = dbConn.PrepareContext(ctx, "DELETE FROM roles WHERE id = $1")
}

func HardDeleteRole(ctx context.Context, id string) error {
	status, err := stmtHardDeleteRole.ExecContext(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	RowsAffected, err := status.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

var stmtSoftDeleteRole *sql.Stmt

func setStatementSoftDeleteRole(ctx context.Context) {
	stmtSoftDeleteRole, _ = dbConn.PrepareContext(ctx, "UPDATE roles SET deleted_at = now(), deleted_by = $1 WHERE id = $2")
}

func SoftDeleteRole(ctx context.Context, id, deletedBy string) error {
	status, err := stmtSoftDeleteRole.ExecContext(ctx, deletedBy, id)
	if err != nil {
		return errors.WithStack(err)
	}

	RowsAffected, err := status.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

var stmtReadAllRole *sql.Stmt

func setStatementReadAllRole(ctx context.Context) {
	stmtReadAllRole, _ = dbConn.PrepareContext(ctx, `
	SELECT r.id, r.name, r.description FROM roles r`)
}

func ReadAllRole(ctx context.Context) ([]Role, error) {
	rows, err := stmtReadAllRole.QueryContext(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var roles []Role

	for rows.Next() {
		role := Role{}
		err = rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		roles = append(roles, role)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return roles, nil
}

var stmtReadRoleByID *sql.Stmt

func setStatementReadRoleByID(ctx context.Context) {
	stmtReadRoleByID, _ = dbConn.PrepareContext(ctx, `
	SELECT r.id, r.name, r.description
	FROM roles r WHERE r.id = $1`)
}

func ReadRoleByID(ctx context.Context, id string) (*Role, error) {
	role := Role{}
	err := stmtReadRoleByID.QueryRowContext(ctx, id).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &role, nil
}
