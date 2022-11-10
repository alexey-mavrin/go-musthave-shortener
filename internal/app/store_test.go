package app

import (
	"testing"
)

func Test_storage_store(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple url",
			args: args{
				url: "http://www.kiae.su/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStorage()
			k, err := s.store(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.store() error = %v, wantErr %v", err, tt.wantErr)
			}
			url, err := s.get(k)
			if err != nil {
				t.Errorf("get returns error %s while should not", err)
			}
			if url != tt.args.url {
				t.Errorf("stored url does not match original %s", tt.args.url)
			}
		})
	}
}
