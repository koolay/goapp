package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineOSSURL(t *testing.T) {
	type args struct {
		endpoint string
		bucket   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				endpoint: "https://oss-cn.shenzhen.aliyun.com",
				bucket:   "goapp",
			},
			want:    "https://goapp.oss-cn.shenzhen.aliyun.com",
			wantErr: false,
		},
		{
			name: "InvalidEndpoint",
			args: args{
				endpoint: "xxaa",
				bucket:   "goapp",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ok2",
			args: args{
				endpoint: "https://oss-cn.shenzhen.aliyun.com/",
				bucket:   "goapp",
			},
			want:    "https://goapp.oss-cn.shenzhen.aliyun.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineOSSURL(tt.args.endpoint, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombineOSSURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CombineOSSURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandString(t *testing.T) {
	got1 := RandString(6)
	assert.Equal(t, len(got1), 6)

	got2 := RandString(8)
	assert.Equal(t, len(got2), 8)

	got22 := RandString(8)
	assert.NotEqual(t, got22, got2)
}
