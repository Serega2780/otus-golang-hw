package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// Place your code here.
			in: User{
				ID:     "111111111111111111111111111111111111",
				Name:   "user",
				Age:    49,
				Email:  "user@mail.ru",
				Role:   "admin",
				Phones: []string{"22222222222", "33333333333", "44444444444"},
			},
			expectedErr: nil,
		},
		{
			// Place your code here.
			in: User{
				ID:     "11111111111111111111111111111111111",
				Name:   "user",
				Age:    49,
				Email:  "user@mail.ru",
				Role:   "admin",
				Phones: []string{"22222222222", "33333333333", "44444444444"},
			},
			expectedErr: make(ValidationErrors, 1),
		},
		{
			// Place your code here.
			in: User{
				ID:     "11111111111111111111111111111111111",
				Name:   "user",
				Age:    80,
				Email:  "user@mail.ru",
				Role:   "nobody",
				Phones: []string{"2222222222", "33333333333", "44444444444"},
			},
			expectedErr: make(ValidationErrors, 4),
		},
		{
			// Place your code here.
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			// Place your code here.
			in:          App{Version: "1234"},
			expectedErr: make(ValidationErrors, 1),
		},
		{
			// Place your code here.
			in:          Token{Header: make([]byte, 5), Payload: make([]byte, 5), Signature: make([]byte, 5)},
			expectedErr: nil,
		},
		{
			// Place your code here.
			in:          Response{Code: 200, Body: "just_a_body"},
			expectedErr: nil,
		},
		{
			// Place your code here.
			in:          Response{Code: 300, Body: "just_a_body"},
			expectedErr: make(ValidationErrors, 1),
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			er := &ValidationErrors{}
			expEr := &ValidationErrors{}
			if errors.As(err, er) {
				if errors.As(tt.expectedErr, expEr) {
					require.Equal(t, er.Len(), expEr.Len())
				}
			} else {
				require.Equal(t, tt.expectedErr, err)
			}
			_ = tt
		})
	}
}
