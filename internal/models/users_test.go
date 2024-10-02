package models

import (
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestUserModel(t *testing.T) {
	db := newTestDB(t)
	m := UserModel{db}

	tests := []struct {
		name     string
		email    string
		password string
		errorMsg error
	}{
		{
			name:     "Valid User Insert",
			email:    "test2@example.com",
			password: "password",
			errorMsg: nil,
		},
		{
			name:     "Duplicate Email",
			email:    "test@example.com",
			password: "password",
			errorMsg: ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Insert(tt.email, tt.password)
			assert.Equal(t, err, tt.errorMsg)
		})
	}
}
