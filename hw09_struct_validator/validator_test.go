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
		Code int    `validate:"in:200,404,500" json:"code"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:   "id3214324564352345671234590876274639",
				Name: "TEST Name", Age: 23, Email: "Correct@email.com",
				Role: "admin", Phones: []string{"12345678912", "22222222222"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID: "id123123212", Name: "TEST Name", Age: 23,
				Email: "WrongEmail.com", Role: "admin", Phones: []string{"12345678912", "22222222222"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrValidateStrLen},
				ValidationError{Field: "Email", Err: ErrValidateStrRegexp},
			},
		},
		{
			in:          Response{Code: 200, Body: "body"},
			expectedErr: ValidationErrors{},
		},
		{
			in:          Response{Code: 232, Body: "body"},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: ErrValidateIn}},
		},
		{
			in:          App{Version: "ver 2"},
			expectedErr: ValidationErrors{},
		},
		{
			in:          App{Version: "version 2"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrValidateStrLen}},
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr.Error())
			_ = tt
		})
	}
}

func TestGetValFuncInt(t *testing.T) {
	tests := []struct {
		name           string
		validateParams string
		in             int
		expexted       error
	}{
		{
			name:           "Validate int: MIN True",
			validateParams: "min:7",
			in:             9,
			expexted:       nil,
		},
		{
			name:           "Validate int: MAX True",
			validateParams: "max:7",
			in:             4,
			expexted:       nil,
		},
		{
			name:           "Validate int: MAX True",
			validateParams: "max:7",
			in:             4,
			expexted:       nil,
		},
		{
			name:           "Validate int: IN True",
			validateParams: "in:7,8,10,15",
			in:             10,
			expexted:       nil,
		},
		{
			name:           "Validate int: COMBINATE True",
			validateParams: "min:7|max:15|in:8,9,10,12",
			in:             12,
			expexted:       nil,
		},
		{
			name:           "Validate int: MIN False",
			validateParams: "min:7",
			in:             5,
			expexted:       ErrValidateIntMin,
		},
		{
			name:           "Validate int: IN False",
			validateParams: "in:7,8,10,15",
			in:             11,
			expexted:       ErrValidateIn,
		},

		{
			name:           "Validate int: MAX False",
			validateParams: "max:7",
			in:             9,
			expexted:       ErrValidateIntMax,
		},
		{
			name:           "Validate int: COMBINATE False",
			validateParams: "min:7|max:15|in:8,9,10,12",
			in:             11,
			expexted:       errors.Join(ErrValidateIn),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := getValFuncInt(test.validateParams)
			err := f(test.in)
			if test.expexted != nil {
				require.EqualError(t, err, test.expexted.Error())
			} else {
				require.Equal(t, test.expexted, err)
			}
		})
	}
}

func TestGetValFuncString(t *testing.T) {
	tests := []struct {
		name           string
		validateParams string
		in             string
		expexted       error
	}{
		{
			name:           "Validate string: LEN True",
			validateParams: "len:5",
			in:             "hello",
			expexted:       nil,
		},
		{
			name:           "Validate string: REGEXP True",
			validateParams: "regexp:\\d+",
			in:             "67234",
			expexted:       nil,
		},
		{
			name:           "Validate string: IN True",
			validateParams: "in:hello,bye,world",
			in:             "bye",
			expexted:       nil,
		},
		{
			name:           "Validate string: COMBINATE True",
			validateParams: "len:5|regexp:\\d+|in:hello,67234,res",
			in:             "67234",
			expexted:       nil,
		},
		{
			name:           "Validate string: len False",
			validateParams: "len:7",
			in:             "hello",
			expexted:       ErrValidateStrLen,
		},
		{
			name:           "Validate string: REGEXP False",
			validateParams: "regexp:\\d+",
			in:             "bye",
			expexted:       ErrValidateStrRegexp,
		},
		{
			name:           "Validate string: IN False",
			validateParams: "in:hello,bye,world",
			in:             "hell",
			expexted:       ErrValidateIn,
		},
		{
			name:           "Validate string: COMBINATE False",
			validateParams: "len:4|regexp:\\d+|in:hello,67234",
			in:             "hello",
			expexted:       errors.Join(ErrValidateStrLen, ErrValidateStrRegexp),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, _ := getValFuncString(test.validateParams)
			err := f(test.in)
			if test.expexted != nil {
				require.EqualError(t, err, test.expexted.Error())
			} else {
				require.Equal(t, test.expexted, err)
			}
		})
	}
}
