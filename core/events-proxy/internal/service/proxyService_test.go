package service

import (
	"context"
	"testing"
)

func Test_proxyService_Auth(t *testing.T) {
	type args struct {
		ctx            context.Context
		username       string
		hashedPassword string
	}
	tests := []struct {
		name    string
		svc     *proxyService
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.svc.Auth(tt.args.ctx, tt.args.username, tt.args.hashedPassword); (err != nil) != tt.wantErr {
				t.Errorf("proxyService.Auth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
