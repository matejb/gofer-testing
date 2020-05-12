package concurrent_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/matejb/gofer-testing/concurrent"
)

func TestOperation(t *testing.T) {
	type testCase struct {
		value          int
		expectedResult int
		expectedErr    error
	}

	cases := []testCase{
		{value: 1, expectedResult: 10},
		{value: 5, expectedResult: 50},
		{value: 10, expectedResult: 100},
		{value: 11, expectedResult: 110, expectedErr: errors.New("can't be larger then 100")},
		{value: 20, expectedResult: 200, expectedErr: errors.New("can't be larger then 100")},
	}

	op := concurrent.New()

	tester := func(t *testing.T, tc testCase) {
		result, err := op.Do(tc.value)

		if result != tc.expectedResult {
			t.Errorf("value %d: expected %d got %d", tc.value, tc.expectedResult, result)
		}

		err = errorsEquality(err, tc.expectedErr)
		if err != nil {
			t.Errorf("value %d: %s", tc.value, err)
		}
	}

	var wg sync.WaitGroup
	for _, tc := range cases {
		wg.Add(1)
		go func(tc testCase) {
			defer wg.Done()
			tester(t, tc)
		}(tc)
	}
	wg.Wait()
}

func errorsEquality(got, expected error) error {
	switch {
	case got == nil && expected == nil:
		return nil
	case got != nil && expected == nil:
		return fmt.Errorf("unexpected error: %s", got)
	case got == nil && expected != nil:
		return fmt.Errorf("expected error")
	case got.Error() != expected.Error():
		return fmt.Errorf("expected error %q but got %q", expected, got)
	}
	return nil
}
