package main_test

import (
	"errors"
	"fmt"
	"testing"

	main "github.com/matejb/gofer-testing/weather/with-good-tests"
)

func TestIsForShorts(t *testing.T) {
	type testCase struct {
		name                string
		sourceTemperature   int
		sourceError         error
		expectedIsForShorts bool
		expectedError       error
	}

	cases := []testCase{
		{
			name:                "temperature bellow 21 Celsius is not for shorts",
			sourceTemperature:   20,
			expectedIsForShorts: false,
		},

		{
			name:                "temperature above 20 Celsius is for shorts",
			sourceTemperature:   21,
			expectedIsForShorts: true,
		},

		{
			name:          "error in weather source returns error",
			sourceError:   errors.New("some error"),
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range cases {
		source := mockWeatherSource{
			temperature: tc.sourceTemperature,
			err:         tc.sourceError,
		}

		isForShorts, err := main.IsForShorts(source)

		if isForShorts != tc.expectedIsForShorts {
			t.Errorf("case %q: expected %v got %v", tc.name, tc.expectedIsForShorts, isForShorts)
		}

		err = errorsEquality(err, tc.expectedError)
		if err != nil {
			t.Errorf("case %q: %s", tc.name, err)
		}
	}

}

type mockWeatherSource struct {
	temperature int
	err         error
}

func (m mockWeatherSource) CurrentTemperature() (int, error) {
	return m.temperature, m.err
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
