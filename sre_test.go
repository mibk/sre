package sre

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		0:  {`ab+`, `abbbb`, true},
		1:  {`ab+`, `a`, false},
		2:  {`x+y`, `xxxxxy`, true},
		3:  {`x+y`, `xyy`, false},
		4:  {`č+řř`, `ččřř`, true},
		5:  {`a?bb?ccc?d*e*`, `bbcce`, true},
		6:  {`\..\.`, `.+.`, true},
		7:  {`[ab]+`, `abbbaaab`, true},
		8:  {`[ab]+`, `abbbaaabx`, false},
		9:  {`ab|cd|ef`, `ab`, true},
		10: {`ab|cd|ef`, `cd`, true},
		11: {`ab|cd|ef`, `ef`, true},
		12: {`(xy)+`, `xyxyxy`, true},
		13: {`[^^]`, `^`, false},
		14: {`[^\^]`, `^`, false},
		15: {`\^[^\^]`, `^r`, true},

		16: {`x{}]`, `x{}]`, true},
		17: {`x{2}`, `xx`, true},
		18: {`x{2}`, `x`, false},
		19: {`x{2}`, `xxx`, false},
		20: {`x{2,5}`, `xx`, true},
		21: {`x{2,5}`, `xxxxx`, true},
		22: {`x{2,5}`, `x`, false},
		23: {`x{2,5}`, `xxxxxx`, false},
		24: {`x{,5}`, ``, true},
		25: {`x{2,}`, `xxxxx`, true},

		26: {`[0-9]+`, `2092389034`, true},
		27: {`[0-9-]+`, `209-238-9034`, true},
		28: {`[-0-9]+`, `209-238-9034`, true},
		29: {`[^0-9-]+`, `209-238-9034`, false},
	}

	for i, tt := range tests {
		rx, err := Compile(tt.pattern)
		if err != nil {
			t.Errorf("test[%d]: unexpected err: %v", i, err)
		}
		if got := rx.Match([]byte(tt.input)); got != tt.want {
			t.Errorf("test[%d]: got %v, want %v", i, got, tt.want)
		}
	}
}
