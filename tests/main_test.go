package tests

import "testing"

func TestStart(t *testing.T) {
	TestUserCreation(t)
	TestCreateRole(t)
	TestRoleWithPermissions(t)
	TestRoleCreationAttachDetach(t)
}
