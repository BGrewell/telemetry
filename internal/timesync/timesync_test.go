package timesync

import (
	"testing"
	"time"
)

func BenchmarkDelay(b *testing.B) {
	t1 := time.Now().UnixNano()
	t2 := time.Now().UnixNano()
	t3 := time.Now().UnixNano()
	t4 := time.Now().UnixNano()
	for i := 0; i < b.N; i++ {
        Delay(t1, t2, t3, t4)
    }
}

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
		{
			name: "1ns delay",
			args: args{
				t1: 1257894000000000000,
				t2: 1257894000000000400,
				t3: 1257894000000000500,
				t4: 1257894000000000101,
			},
			want: int64(1 * time.Nanosecond),
		},
		{
			name: "100ms delay",
			args: args{
				t1: 1257894000000000000,
                t2: 1257894000000000400,
                t3: 1257894000000000500,
                t4: 1257894000100000100,
			},
			want: int64(100 * time.Millisecond),
		},
		{
			name: "1s delay",
			args: args{
                t1: 1257894000000000000,
                t2: 1257894000000000400,
                t3: 1257894000000000500,
                t4: 1257894001000000100,
            },
			want: int64(1 * time.Second),
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delay(tt.args.t1, tt.args.t2, tt.args.t3, tt.args.t4); got != tt.want {
				t.Errorf("Delay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkOffset(b *testing.B) {
	t1 := time.Now().UnixNano()
	t2 := time.Now().UnixNano()
	t3 := time.Now().UnixNano()
	t4 := time.Now().UnixNano()
	for i := 0; i < b.N; i++ {
        Offset(t1, t2, t3, t4)
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
		{
			name: "1ns offset",
			args: args{
                t1: 1257894000000000000,
                t2: 1257894000100000001,
                t3: 1257894000110000001,
                t4: 1257894000210000000,
            },
			want: int64(1 * time.Nanosecond),
		},
		{
            name: "100ms offset",
            args: args{
                t1: 1257894000000000000,
                t2: 1257894000200000000,
                t3: 1257894000210000000,
                t4: 1257894000210000000,
            },
            want: int64(100 * time.Millisecond),
        },
        {
            name: "1s offset",
            args: args{
                t1: 1257894000000000000,
                t2: 1257894001100000000,
                t3: 1257894001110000000,
                t4: 1257894000210000000,
            },
            want: int64(1 * time.Second),
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Offset(tt.args.t1, tt.args.t2, tt.args.t3, tt.args.t4); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}
