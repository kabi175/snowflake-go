package generator_test

import (
	"testing"
	"time"

	generator "github.com/kabi175/snowflake-go"
)

type genCase struct {
	arg0 int64
	arg1 int64
}

var newGenSuccess = []genCase{
	{arg0: 0, arg1: 0},
	{arg0: 1023, arg1: 0},
	{arg0: 0, arg1: time.Now().UnixMilli()},
	{arg0: 478, arg1: time.Now().UnixMilli()},
	{arg0: 1023, arg1: time.Now().UnixMilli()},
}

var newGenFails = []genCase{
	{arg0: -1, arg1: 0},
	{arg0: 0, arg1: -1},
	{arg0: 0, arg1: time.Now().UnixMilli() + (1000 * 2)},
	{arg0: 1024, arg1: -1},
	{arg0: 1009, arg1: -1},
}

func TestNewGenerator(t *testing.T) {
	t.Parallel()
	var gen *generator.Generator
	var err error

	t.Run("should pass", func(t *testing.T) {
		for _, test := range newGenSuccess {
			t.Logf("Case arg0:%v arg1:%v", test.arg0, test.arg1)
			gen, err = generator.NewGenerator(test.arg0, test.arg1)
			if gen == nil {
				t.Errorf("gen: can't be nil")
			}
			if err != nil {
				t.Errorf("err: expected nil got %e", err)
			}
		}
	})

	t.Run("should fail", func(t *testing.T) {
		for _, test := range newGenFails {
			t.Logf("Case arg0:%v arg1:%v", test.arg0, test.arg1)
			gen, err = generator.NewGenerator(test.arg0, test.arg1)
			if gen != nil {
				t.Errorf("gen: should be nil")
			}
			if err == nil {
				t.Errorf("err: expected error got  nil")
			}
		}
	})

}

func TestNext(t *testing.T) {
	t.Parallel()
	t.Run("same id check", func(t *testing.T) {
		t.Parallel()
		gen1, err := generator.NewGenerator(0, 0)
		if err != nil {
			t.Errorf("err: expected nil got %e", err)
		}

		gen2, err := generator.NewGenerator(1, 0)
		if err != nil {
			t.Errorf("err: expected nil got %e", err)
		}

		for iter := 0; iter < 1000; iter++ {
			next1, err1 := gen1.Next()
			next2, err2 := gen2.Next()
			if err1 != nil || err2 != nil {
				continue
			}
			if next1 == next2 {
				t.Error("ids should't be same")
			}
		}
	})

	t.Run("backward id check", func(t *testing.T) {
		t.Parallel()
		gen, err := generator.NewGenerator(0, 0)
		if err != nil {
			t.Errorf("err: expected nil got %e", err)
		}

		var prev int64
		for iter := 0; iter < 10000; iter++ {
			curr, err := gen.Next()
			if curr == 0 && err == nil {
				t.Error("err: expected error got nil")
			}
			if curr < prev {
				t.Errorf("curr id value less than prev id curr:%v prev:%v", curr, prev)
			}
		}
	})

	t.Run("negative id check", func(t *testing.T) {
		t.Parallel()
		gen, err := generator.NewGenerator(0, 0)
		if err != nil {
			t.Errorf("err: expected nil got %e", err)
		}

		for iter := 0; iter < 10000; iter++ {
			_, err := gen.Next()
			if err != nil && err.Error() == "negative id genrated" {
				t.Errorf("err: expected nil got %v", err)
			}
		}
	})
}
