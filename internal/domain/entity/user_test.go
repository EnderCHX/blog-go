package entity

import (
	"testing"
)

func TestGetEmail(t *testing.T) {
	user := NewUser("test", "test", "test@test.com", "test", "test")
	email, err := user.GetEmail("https://api.passport.hrbeu.top/user/info", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiaXNzIjoiY2h4Yy5jYyIsImV4cCI6MTc0MDk5Njc2NywiaWF0IjoxNzQwOTkzMTY3fQ.f0T8OoqhHo5gB4nVPP4hxK1C7uCefi5CqHJwULABfIo")
	if err != nil {
		t.Errorf("GetEmail() error = %v", err)
	}
	t.Logf("GetEmail() = %v", email)
}
