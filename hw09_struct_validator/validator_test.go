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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
		in   interface{}
		errs []ValidationError
	}{
		{
			in: User{
				ID:     "54545dwewedfgfgfegasdvdsfas54545gfgf",
				Name:   "otus",
				Age:    22,
				Email:  "petya@yandex.ru",
				Role:   "stuff",
				Phones: []string{"79276176855", "79296179871"},
				meta:   []byte("[]"),
			},
			errs: []ValidationError{},
		},
		{
			in: App{
				Version: "53342",
			},
			errs: []ValidationError{},
		},
		{
			in: Token{
				Header: []byte("ghsrguida"),
			},
			errs: []ValidationError{},
		},
		{
			in: Response{
				Code: 2000,
			},
			errs: []ValidationError{
				{
					Field: "Response",
					Err:   ErrGroup,
				},
			},
		},
		{
			in: App{
				Version: "1313",
			},
			errs: []ValidationError{
				{
					Field: "App",
					Err:   ErrLen,
				},
			},
		},
		{
			in: User{
				ID:     "a2323sdg5454ga232323fasfgfhgjgj",
				Name:   "test",
				Age:    10,
				Email:  "@test@t@gmail.com",
				Role:   "admin_admin",
				Phones: []string{"79624678025", "111"},
				meta:   []byte("[]"),
			},
			errs: []ValidationError{
				{
					Field: "ID",
					Err:   ErrLen,
				},
				{
					Field: "Age",
					Err:   ErrMin,
				},
				{
					Field: "Email",
					Err:   ErrRegexp,
				},
				{
					Field: "Role",
					Err:   ErrInString,
				},
				{
					Field: "Phones",
					Err:   ErrLen,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			errs := Validate(tt.in)
			require.Equal(t, len(errs), len(tt.errs))
			for i := range errs {
				require.True(t, errors.Is(errs[i].Err, tt.errs[i].Err))
			}
		})
	}
}

// не структура

func TestCopyErrorOffset(t *testing.T) {
	type notStruct struct {
		in   string
		errs []ValidationError
	}
	errs := []ValidationError{
		{
			Field: "none",
			Err:   ErrNotStruct,
		},
	}
	tc := notStruct{
		in:   "string",
		errs: errs,
	}

	t.Run("interface not struct", func(t *testing.T) {
		errs := Validate(tc.in)
		for i := range errs {
			require.True(t, errors.Is(errs[i].Err, tc.errs[i].Err))
		}
	})
}
