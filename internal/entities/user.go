package entities

type User struct {
	InternalId     string      `json:"internal_id"`
	Email          string      `json:"email"`
	Phone          string      `json:"phone"`
	HashedPassword string      `json:"___password"`
	Roles          []Role      `json:"___roles"`
	Data           interface{} `json:"data"`
}

func (u *User) AddRole(role Role) {
	updated := u.UpdateRole(role)
	if !updated {
		u.Roles = append(u.Roles, role)
	}
}

func (u *User) UpdateRole(updatedRole Role) bool {
	for i, role := range u.Roles {
		if role.InternalID == updatedRole.InternalID {
			u.Roles[i] = updatedRole
			return true
		}
	}
	return false
}

func (u *User) DeleteRole(roleID string) {
	for i, role := range u.Roles {
		if role.InternalID == roleID {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			break
		}
	}
}
