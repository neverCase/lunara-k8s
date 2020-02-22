package service

import (
	"context"
	"os"
	"testing"

	"github.com/nevercase/lunara-k8s/configs"
)

func TestService_checkAuth(t *testing.T) {
	type fields struct {
		c           *configs.Config
		output      *os.File
		httpService *httpService
		ctx         context.Context
		cancel      context.CancelFunc
	}
	type args struct {
		login *loginRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test_1",
			fields: fields{
				httpService: &httpService{},
			},
			args: args{
				login: &loginRequest{
					Account: "abc",
					Password: "111",
				},
			},
			want: false,

		},
		{
			name: "test_2",
			fields: fields{},
			args: args{
				login: &loginRequest{
					Account: "admin",
					Password: "123456",
				},
			},
			want: true,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				c:           tt.fields.c,
				output:      tt.fields.output,
				httpService: tt.fields.httpService,
				ctx:         tt.fields.ctx,
				cancel:      tt.fields.cancel,
			}
			if got := s.httpService.checkAuth(tt.args.login); got != tt.want {
				t.Errorf("Service.checkAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
