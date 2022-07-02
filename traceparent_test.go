package trex

import (
	"testing"
)

func TestGenerateRadomHexString(t *testing.T) {

	cases := []struct {
		name   string // Name of the test case
		bytes  int    // Number of bytes to process
		length int    // Character length of the result
		hasErr bool   // Should the function throw error
	}{
		{"normal case", 6, 12, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str, err := GenerateRadomHexString(c.bytes)

			if c.hasErr == true && err == nil {
				t.Error("expected an error but it's nil", err)
			} else if c.hasErr == false && err != nil {
				t.Error("expected no error but it's not nil", err)
			}

			if len(str) != c.length {
				t.Errorf("string should be %d characters, but got %d", c.length, len(str))
			}
		})
	}
}

func TestDecode(t *testing.T) {

	cases := []struct {
		name    string
		trcprnt string
		ver     string
		tid     string
		pid     string
		flg     string
		hasErr  bool
	}{
		{
			"normal case",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
			"00",
			"0123456789abcdef0123456789abcdef",
			"0123456789abcdef",
			"01",
			false,
		},
		{
			"empty traceparent",
			"",
			"", "", "", "", true,
		},
		{
			"incorrect version format",
			"0x-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
			"", "", "", "", true,
		},
		{
			"incorrect version length",
			"0-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
			"", "", "", "", true,
		},
		{
			"incorrect trace id format",
			"00-0123456789abcdef0123456789abcdex-0123456789abcdef-01",
			"", "", "", "", true,
		},
		{
			"incorrect trace id length",
			"00-0-0123456789abcdef-01",
			"", "", "", "", true,
		},
		{
			"incorrect parent id format",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdex-01",
			"", "", "", "", true,
		},
		{
			"incorrect parent id length",
			"00-0123456789abcdef0123456789abcdef-0-01",
			"", "", "", "", true,
		},
		{
			"incorrect flag format",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdef-0x",
			"", "", "", "", true,
		},
		{
			"incorrect flag length",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdef-0",
			"", "", "", "", true,
		},
		{
			"trace id is zero",
			"00-00000000000000000000000000000000-0123456789abcdef-0",
			"", "", "", "", true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ver, tid, pid, flg, err := DecodeTraceparent(c.trcprnt)

			if c.hasErr == true && err == nil {
				t.Error("expected an error but it's nil", err)
			} else if c.hasErr == false && err != nil {
				t.Error("expected no error but it's not nil", err)
			}

			if ver != c.ver {
				t.Errorf("version should be \"%s\" , but got \"%s\"", c.ver, ver)
			}
			if tid != c.tid {
				t.Errorf("trace id should be \"%s\" , but got \"%s\"", c.tid, tid)
			}
			if pid != c.pid {
				t.Errorf("parent id should be \"%s\" , but got \"%s\"", c.pid, pid)
			}
			if flg != c.flg {
				t.Errorf("flag should be \"%s\" , but got \"%s\"", c.flg, flg)
			}
		})
	}
}
