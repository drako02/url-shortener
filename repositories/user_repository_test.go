package repositories

import "testing"

func Test_isValidUserField(t *testing.T) {
	type args struct {
		fields map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid map",
			args{
				fields: map[string]string{
					"first_name": "andrew",
					"last_name":  "appah",
					"email":     "jjj",
				},
			},
			true},
		{"invalid map",
			args{
				fields: map[string]string{
					"first_name": "andrew",
					"last_name":  "appah",
					"created_at":     "jjj",
				},
			},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidUserField(tt.args.fields); got != tt.want {
				t.Errorf("isValidUserField() = %v, want %v", got, tt.want)
			}
		})
	}
}
