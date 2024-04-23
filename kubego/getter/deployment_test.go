package getter

import (
	appsv1 "k8s.io/api/apps/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestGetDeployment(t *testing.T) {
	type args struct {
		ns   string
		name string
		cl   client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *appsv1.Deployment
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
		//	name 是‘app-name’, args 是{ns: 'default', name: 'app-name', cl: fakeClient}
		//	期望返回的是一个Deployment对象，和一个bool值为true，没有错误

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetDeployment(tt.args.ns, tt.args.name, tt.args.cl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDeployment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetDeployment() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
