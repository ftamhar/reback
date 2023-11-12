package reback

import "context"

func IsUserCanRead(ctx context.Context, permissions []Permission) bool {
	for _, permission := range permissions {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		if permission.IsRead != nil && *permission.IsRead {
			return true
		}
	}
	return false
}

func IsUserCanCreate(ctx context.Context, permissions []Permission) bool {
	for _, permission := range permissions {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		if permission.IsCreate != nil && *permission.IsCreate {
			return true
		}
	}
	return false
}

func IsUserCanUpdate(ctx context.Context, permissions []Permission) bool {
	for _, permission := range permissions {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		if permission.IsUpdate != nil && *permission.IsUpdate {
			return true
		}
	}
	return false
}

func IsUserCanDelete(ctx context.Context, permissions []Permission) bool {
	for _, permission := range permissions {
		select {
		case <-ctx.Done():
			return false
		default:
		}
		if permission.IsDelete != nil && *permission.IsDelete {
			return true
		}
	}
	return false
}
