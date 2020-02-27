package app

import "testing"

func Test_parseToken(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "bad argument",
			args: args{
				value: "foo bar",
			},
			want: "",
			wantErr: true,
		},
		{
			name: "good argument",
			args: args{
				value: "Bearer token",
			},
			want: "token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseToken(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
