package cascadeenv

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
)

func Test_checkOSEnv(t *testing.T) {
	type args struct {
		names *[]string
	}
	os.Setenv("a", "1")
	os.Setenv("b", "2")
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{names: &[]string{"a", "b"}},
			want: true,
		},
		{
			name: "2",
			args: args{names: &[]string{"c", "d"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkOSEnv(tt.args.names); got != tt.want {
				t.Errorf("checkOSEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadAndCheckEnv(t *testing.T) {
	type args struct {
		names    *[]string
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{names: &[]string{"a1", "b1"}, filename: "test.ENV"},
			want: true,
		},
		{
			name: "2",
			args: args{names: &[]string{"c", "d"}, filename: "test.ENV"},
			want: false,
		},
		{
			name: "2",
			args: args{names: &[]string{"c1", "c2"}, filename: "test.ENV"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadAndCheckEnv(tt.args.names, tt.args.filename); got != tt.want {
				t.Errorf("loadAndCheckEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkAWSParamStore(t *testing.T) {
	type args struct {
		names   *[]string
		session *session.Session
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{names: &[]string{"a1", "b1"}, session: nil},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkAWSParamStore(tt.args.names, tt.args.session); got != tt.want {
				t.Errorf("checkAWSParamStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitEnvVar(t *testing.T) {
	type args struct {
		names       []string
		envFilename string
		session     *session.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{names: []string{}, envFilename: "test.ENV", session: nil},
			wantErr: true,
		},
		{
			name:    "2",
			args:    args{names: []string{"a13", "b13"}, envFilename: "test2.ENV", session: nil},
			wantErr: true,
		},
		{
			name:    "3",
			args:    args{names: []string{"a1", "b1"}, envFilename: "test.ENV", session: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitEnvVar(tt.args.names, tt.args.envFilename, tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("InitEnvVar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExportEnvVar(t *testing.T) {
	type args struct {
		names       []string
		types       []reflect.Kind
		envFilename string
		session     *session.Session
	}
	os.Setenv("a", "1")
	os.Setenv("b", "string")
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				names:       []string{"a", "b"},
				types:       []reflect.Kind{reflect.Int64, reflect.String},
				envFilename: "test.ENV",
				session:     nil,
			},
			want: map[string]interface{}{
				"a": int64(1),
				"b": "string",
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				names:       []string{"c1", "c2"},
				types:       []reflect.Kind{reflect.String, reflect.String},
				envFilename: "test.ENV",
				session:     nil,
			},
			want: map[string]interface{}{
				"c1": "",
				"c2": "",
			},
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				names:       []string{"c3"},
				types:       []reflect.Kind{reflect.String},
				envFilename: "test.ENV",
				session:     nil,
			},
			want:    map[string]interface{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExportEnvVar(tt.args.names, tt.args.types, tt.args.envFilename, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportEnvVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}
