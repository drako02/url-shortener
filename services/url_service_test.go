package services

import (
	"fmt"
	"testing"
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

func TestDeleteUrl(t *testing.T) {
	type args struct {
		shortCode string
	}

	tests := []struct{
		name string
		arg args
		wants error

		}{
		{"valid", args{"hhhhh"}, nil},
		{"empty string", args{""}, fmt.Errorf("")},
		{"non existing code", args{"jjj"}, fmt.Errorf("")},
	}

}
