package models

import (
	"fmt"
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestUserModel(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := UserModel{db}

	testsInsert := []struct {
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

	for _, tt := range testsInsert {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Insert(tt.email, tt.password)
			assert.Equal(t, err, tt.errorMsg)
		})
	}

	testsAuthenticate := []struct {
		name     string
		email    string
		password string
		id       int
		errorMsg error
	}{
		{
			name:     "Valid User Authentication",
			email:    "test@example.com",
			password: "password",
			id:       1,
			errorMsg: nil,
		},
		{
			name:     "Invalid User Auth, Incorrect Email",
			email:    "incorret@example.com",
			password: "password",
			id:       0,
			errorMsg: ErrInvalidCredentials,
		},
		{
			name:     "Invalid User Auth, Incorrect Password",
			email:    "test@example.com",
			password: "pAssword",
			id:       0,
			errorMsg: ErrInvalidCredentials,
		},
	}

	for _, tt := range testsAuthenticate {
		t.Run(tt.name, func(t *testing.T) {
			id, err := m.Authenticate(tt.email, tt.password)
			assert.Equal(t, id, tt.id)
			assert.Equal(t, err, tt.errorMsg)
		})
	}

	testsExists := []struct {
		name     string
		id       int
		expected bool
		errorMsg error
	}{
		{
			name:     "User Exists",
			id:       1,
			expected: true,
			errorMsg: nil,
		},
		{
			name:     "User Doesn't Exist",
			id:       10,
			expected: false,
			errorMsg: nil,
		},
	}

	for _, tt := range testsExists {
		t.Run(tt.name, func(t *testing.T) {
			exist, err := m.Exists(tt.id)
			assert.Equal(t, exist, tt.expected)
			assert.Equal(t, err, tt.errorMsg)
		})
	}
}
