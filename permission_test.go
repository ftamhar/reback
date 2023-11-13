package reback

import (
	"context"
	"testing"
)

func TestPermission(t *testing.T) {
	roles := []Role{
		{
			Name:        "admin",
			Description: "for admin",
		},
		{
			Name:        "user",
			Description: "for user",
		},
	}

	tr := true
	fl := false

	// permission admin
	permissions := []Permission{
		{
			Resource: "get-all",
			IsCreate: &tr,
			IsRead:   &tr,
			IsUpdate: &tr,
			IsDelete: &tr,
		},
	}

	// permissions user
	userPermissions := []Permission{
		{
			Resource: "get-all",
			IsCreate: &fl,
			IsRead:   &fl,
			IsUpdate: &fl,
			IsDelete: &fl,
		},
	}

	// create roles

	for i, v := range roles {
		id, err := CreateRole(context.Background(), v.Name, v.Description)
		if err != nil {
			t.Fatal(err)
		}

		roles[i].ID = id
	}

	// create permissions
	permissions[0].RoleId = roles[0].ID
	err := CreatePermissions(context.Background(), permissions)
	if err != nil {
		t.Fatal(err)
	}

	userPermissions[0].RoleId = roles[1].ID
	err = CreatePermissions(context.Background(), userPermissions)
	if err != nil {
		t.Fatal(err)
	}
}
