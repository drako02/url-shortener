package services

import (
	"context"
	"errors"
	"testing"

	"github.com/drako02/url-shortener/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestDeleteUrl(t *testing.T) {
// 	type args struct {
// 		shortCode string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := DeleteUrl(tt.args.shortCode); (err != nil) != tt.wantErr {
// 				t.Errorf("DeleteUrl() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_isValidField(t *testing.T) {
	type args struct {
		field any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidField(tt.args.field); got != tt.want {
				t.Errorf("isValidField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidField(t *testing.T) {
	cases := []struct {
		name  string
		field any
		wants bool
	}{
		{"valid string", "id", true},
		{"invalid string", "password", false},
		{"valid slice", []string{"short_code", "clicks"}, true},
		{"slice with one invalid", []string{"id", "bad_field"}, false},
		{"empty slice", []string{}, true},
		{"nil value", nil, false},
		{"unsupported type", 123, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if isValidField(c.field) != c.wants {
				t.Errorf("isValidField(%v) = %v; want %v", c.field, isValidField(c.field), c.wants)
			}
		})
	}
}

// func TestDeleteUrl(t *testing.T) {
// 	type args struct {
// 		shortCode string
// 	}

// 	tests := []struct{
// 		name string
// 		arg args
// 		wants error

// 		}{
// 		{"valid", args{"hhhhh"}, nil},
// 		{"empty string", args{""}, fmt.Errorf("")},
// 		{"non existing code", args{"jjj"}, fmt.Errorf("")},
// 	}

// }

type mockURLUpdater struct {
	mock.Mock
}

func (r *mockURLUpdater) UpdateById(ctx context.Context, id uint, data repositories.Data) error {
	args := r.Called(ctx, id, data)
	return args.Error(0)
}

func TestURLService_SetUrlActiveStatus(t *testing.T) {
	type fields struct {
		updater repositories.Updater
	}

	type args struct {
		ctx   context.Context
		id    uint
		value bool
	}

	tests := []struct {
		name            string
		arg             args
		field           fields
		expectedDataArg repositories.Data
		mockBehaviour   func(repo *mockURLUpdater, data repositories.Data)
		wantErr         bool
	}{
		{
			"should set url active successfully",
			args{ctx: context.Background(), id: 123, value: true},
			fields{updater: new(mockURLUpdater)},
			repositories.Data{Field: "active", Value: true},
			func(repo *mockURLUpdater, data repositories.Data) {
				repo.On("UpdateById", mock.Anything, uint(123), data).Return(nil)

			},
			false,
		},
		{
			"should return an error on update failure",
			args{context.Background(), 123, true},
			fields{updater: new(mockURLUpdater)},
			repositories.Data{Field: "active", Value: true},
			func(repo *mockURLUpdater, data repositories.Data) {
				repo.On("UpdateById", mock.Anything, uint(123), data).Return(errors.New("some error"))

			},
			true,
		},
		{
			"should deactivate url successfully",
			args{ctx: context.Background(), id: 123, value: false},
			fields{updater: new(mockURLUpdater)},
			repositories.Data{Field: "active", Value: false},
			func(repo *mockURLUpdater, data repositories.Data) {
				repo.On("UpdateById", mock.Anything, uint(123), data).Return(nil)

			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpdater := tt.field.updater.(*mockURLUpdater)
			// svc := URLService{updater: mockRepo, deleter: nil}
			svc := &URLService{updater: mockUpdater}

			if tt.mockBehaviour != nil {
				tt.mockBehaviour(mockUpdater, tt.expectedDataArg)
			}

			err := svc.SetUrlActiveStatus(tt.arg.ctx, tt.arg.id, tt.arg.value)

			mockUpdater.AssertExpectations(t)

			if !tt.wantErr {
				assert.NoError(t, err)
			}

			if tt.wantErr != (err != nil) {
				t.Errorf("URLService.SetUrlActive() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
