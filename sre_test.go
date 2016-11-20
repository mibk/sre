package sre

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		0: {`ab+`, `abbbb`, true},
		1: {`ab+`, `a`, false},
		2: {`x+y`, `xxxxxy`, true},
		3: {`x+y`, `xyy`, false},
		4: {`č+řř`, `ččřř`, true},
	}

	for i, tt := range tests {
		rx := MustCompile(tt.pattern)
		if got := rx.Match([]byte(tt.input)); got != tt.want {
			t.Errorf("test[%d]: got %v, want %v", i, got, tt.want)
		}
	}
}
