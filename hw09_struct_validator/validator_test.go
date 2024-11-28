package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"regexp/syntax"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
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

	TestStruct struct {
		Ints  []int `validate:"in:10,20,30"`
		Ints2 []int `validate:"max:100"`
		int3  int   `validate:"min:1000"`
	}

	WrongTestStructTag struct {
		TestInt int `validate:"max"`
	}

	WrongTestStructTagName struct {
		TestInt int `validate:"test:11"`
	}

	WrongTestStructTagValue struct {
		TestInt int `validate:"max:a"`
	}

	WrongTestStructTag2 struct {
		TestString string `validate:"len"`
	}

	WrongTestStructTagName2 struct {
		TestString string `validate:"test:test"`
	}

	WrongTestStructTagValue2 struct {
		TestString string `validate:"len:r"`
	}

	WrongTestStructTagRegex struct {
		TestString string `validate:"regexp:^\\w(+@\\w+\\.\\w+$"`
	}

	AppResponseUser struct {
		A App      `validate:"nested"`
		R Response `validate:"nested"`
		U User     `validate:"nested"`
	}

	AppResponse struct {
		A App
		R Response `validate:"nested"`
	}

	WrongComposeStruct struct {
		W  WrongTestStructTagName   `validate:"nested"`
		W2 WrongTestStructTagValue2 `validate:"nested"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{
				Version: "12345",
			},
			nil,
		},
		{
			Response{
				Code: 404,
			},
			nil,
		},
		{
			Response{
				Code: 200,
			},
			nil,
		},
		{
			User{
				ID:     "123456789",
				Name:   "test1",
				Age:    18,
				Email:  "example@test.com",
				Role:   "admin",
				Phones: []string{"88009997755", "8800999775", "880099977"},
			},
			nil,
		},
		{
			User{
				ID:     "12345678",
				Name:   "test2",
				Age:    50,
				Email:  "example2@test.com",
				Role:   "admin",
				Phones: []string{"88009997755", "8800999775"},
			},
			nil,
		},
		{
			User{
				ID:     "1234567890",
				Age:    36,
				Email:  "example3@test.com",
				Role:   "stuff",
				Phones: []string{"88009997755"},
			},
			nil,
		},
		{
			TestStruct{
				Ints:  []int{10, 20, 30},
				Ints2: []int{1, 100, 0, 50},
				int3:  10,
			},
			nil,
		},
		{
			AppResponseUser{
				A: App{
					Version: "12345",
				},
				R: Response{
					Code: 200,
				},
				U: User{
					ID:     "1234567890",
					Age:    36,
					Email:  "example3@test.com",
					Role:   "stuff",
					Phones: []string{"88009997755"},
				},
			},
			nil,
		},
		{
			AppResponse{
				A: App{
					Version: "123456",
				},
				R: Response{
					Code: 200,
				},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			require.ErrorIs(t, Validate(tt.in), tt.expectedErr)
			_ = tt
		})
	}
}

func TestIncorrectValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{
				Version: "123456",
			},
			ValidationErrors{
				{Field: "Version", Err: ErrTooBigString},
			},
		},
		{
			Response{
				Code: 403,
			},
			ValidationErrors{
				{Field: "Code", Err: ErrNotInSeq},
			},
		},
		{
			User{
				ID:     "111111111111111111111111111111111111111111111111111111111111",
				Name:   "test1",
				Age:    15,
				Email:  "exampletest.com",
				Role:   "admi",
				Phones: []string{"880099977552", "880099977512", "8800999771234"},
			},
			ValidationErrors{
				{Field: "ID", Err: ErrTooBigString},
				{Field: "Age", Err: ErrTooSmallInt},
				{Field: "Email", Err: ErrInvalidByRegexString},
				{Field: "Role", Err: ErrNotInSeq},
				{Field: "Phones", Err: ErrTooBigString},
			},
		},
		{
			User{
				ID:     "111111111",
				Name:   "test1",
				Age:    195,
				Email:  "exampletest.com",
				Role:   "admi",
				Phones: []string{"880099977552", "8800999775", "880099977"},
			},
			ValidationErrors{
				{Field: "Age", Err: ErrTooBigInt},
				{Field: "Email", Err: ErrInvalidByRegexString},
				{Field: "Role", Err: ErrNotInSeq},
				{Field: "Phones", Err: ErrTooBigString},
			},
		},
		{
			TestStruct{
				Ints:  []int{0, 40, 50},
				Ints2: []int{1, 100, 0, 500},
			},
			ValidationErrors{
				{Field: "Ints", Err: ErrNotInSeq},
				{Field: "Ints2", Err: ErrTooBigInt},
			},
		},
		{
			AppResponseUser{
				A: App{
					Version: "123456",
				},
				R: Response{
					Code: 403,
				},
				U: User{
					ID:     "111111111111111111111111111111111111111111111111111111111111",
					Name:   "test1",
					Age:    15,
					Email:  "exampletest.com",
					Role:   "admi",
					Phones: []string{"880099977552", "880099977512", "8800999771234"},
				},
			},
			ValidationErrors{
				{Field: "A.Version", Err: ErrTooBigString},
				{Field: "R.Code", Err: ErrNotInSeq},
				{Field: "U.ID", Err: ErrTooBigString},
				{Field: "U.Age", Err: ErrTooSmallInt},
				{Field: "U.Email", Err: ErrInvalidByRegexString},
				{Field: "U.Role", Err: ErrNotInSeq},
				{Field: "U.Phones", Err: ErrTooBigString},
			},
		},
		{
			AppResponse{
				A: App{
					Version: "12345",
				},
				R: Response{
					Code: 201,
				},
			},
			ValidationErrors{
				{Field: "R.Code", Err: ErrNotInSeq},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			require.ErrorAs(t, err, &tt.expectedErr)

			for _, e := range strings.Split(tt.expectedErr.Error(), "\n") {
				require.ErrorContains(t, err, e)
			}
		})
	}
}

func TestIncorrectValidateExecute(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			1234,
			ErrNotStruct,
		},
		{
			WrongTestStructTag{
				TestInt: 12,
			},
			ErrInvalidTag,
		},
		{
			WrongTestStructTagName{
				TestInt: 12,
			},
			ErrInvalidTag,
		},
		{
			WrongTestStructTagValue{
				TestInt: 12,
			},
			strconv.ErrSyntax,
		},
		{
			WrongTestStructTag2{
				TestString: "test",
			},
			ErrInvalidTag,
		},
		{
			WrongTestStructTagName2{
				TestString: "test",
			},
			ErrInvalidTag,
		},
		{
			WrongTestStructTagValue2{
				TestString: "test",
			},
			strconv.ErrSyntax,
		},
		{
			WrongComposeStruct{
				W:  WrongTestStructTagName{TestInt: 1},
				W2: WrongTestStructTagValue2{TestString: "rewr"},
			},
			ErrInvalidTag,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestIncorrectValidateRegex(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			WrongTestStructTagRegex{
				TestString: "12",
			},
			&syntax.Error{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorAs(t, err, &tt.expectedErr)
		})
	}
}
