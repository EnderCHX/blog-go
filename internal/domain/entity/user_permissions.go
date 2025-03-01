package entity

type UserPermissions struct {
	Username    string `json:"username" gorm:"primaryKey"`
	Permissions string `json:"permissions"`
}

type UserPermissionsOptions func(*UserPermissions)

func NewUserPermissions(option ...UserPermissionsOptions) *UserPermissions {
	userPermissions := &UserPermissions{}
	for _, opt := range option {
		opt(userPermissions)
	}
	return userPermissions
}

func WithPermissions(permissions string) UserPermissionsOptions {
	return func(userPermissions *UserPermissions) {
		userPermissions.Permissions = permissions
	}
}
