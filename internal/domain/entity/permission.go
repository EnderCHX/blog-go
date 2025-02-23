package entity

type Permission struct {
	PermissionId int    `json:"permission_id" gorm:"primaryKey"`
	Permission   string `json:"permission"`
}

var Permissions = []Permission{
	{PermissionId: 0, Permission: "access"},  //访问
	{PermissionId: 1, Permission: "read"},    //读文章
	{PermissionId: 2, Permission: "edit"},    //写文章
	{PermissionId: 3, Permission: "comment"}, //评论
}
