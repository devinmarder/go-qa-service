package repository

import (
	"reflect"
	"testing"
)

func TestLocalRepository_UpdateServiceCoverage(t *testing.T) {
	type fields struct {
		Services []ServiceCoverage
	}
	type args struct {
		ServiceCoverage
	}
	tests := []struct {
		name     string
		initial  fields
		args     args
		expected fields
	}{
		{
			"add new service",
			fields{},
			args{ServiceCoverage{ServiceName: "service1", Coverage: 75}},
			fields{Services: []ServiceCoverage{{ServiceName: "service1", Coverage: 75}}},
		},
		{
			"update existing service",
			fields{Services: []ServiceCoverage{{ServiceName: "service1", Coverage: 25}}},
			args{ServiceCoverage{ServiceName: "service1", Coverage: 75}},
			fields{Services: []ServiceCoverage{{ServiceName: "service1", Coverage: 75}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lr := &LocalRepository{
				Services: tt.initial.Services,
			}
			lr.UpdateServiceCoverage(tt.args.ServiceCoverage)
			if !reflect.DeepEqual(lr.Services, tt.expected.Services) {
				t.Fatalf("expected state: %v \n actual state: %v ", tt.expected, lr.Services)
			}
		})
	}
}
