package timesync

import "testing"

func TestDelay(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
		t3 int64
		t4 int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delay(tt.args.t1, tt.args.t2, tt.args.t3, tt.args.t4); got != tt.want {
				t.Errorf("Delay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	type args struct {
		t1 int64
		t2 int64
		t3 int64
		t4 int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Offset(tt.args.t1, tt.args.t2, tt.args.t3, tt.args.t4); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}
