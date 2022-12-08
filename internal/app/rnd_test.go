package app

import "testing"

func Test_randSeq(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "len of 0",
			args: args{
				n: 0,
			},
		},
		{
			name: "len of 3",
			args: args{
				n: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randSeq(tt.args.n)
			if len(got) != tt.args.n {
				t.Errorf("randSeq() = %v, want len %d", got, tt.args.n)
			}

			if tt.args.n == 0 {
				return
			}

			got2 := randSeq(tt.args.n)
			if got == got2 {
				t.Errorf("%v == %v, want !=", got, got2)
			}
		})
	}
}
