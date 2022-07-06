package trex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type TCtx struct {
	ver, tid, pid, rid, flg string
}

type MockTxFactory struct {
	fail bool
}

func (f *MockTxFactory) Generate(
	ver string,
	tid string,
	pid string,
	rid string,
	flg string,
) (interface{}, error) {
	if f.fail {
		return nil, fmt.Errorf("Some Error")
	} else {
		return &TCtx{ver, tid, pid, rid, flg}, nil
	}
}

func TestTxContextMiddleware(t *testing.T) {
	cases := []struct {
		name    string // Test case name
		trcprnt string // Header value of traceparent
		decFail bool   // Decoding of traceparent will fail
		txfFail bool   // Generate Transaction Context will fail
		status  int    // HTTP response status
	}{
		{
			"normal case",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
			false,
			false,
			200,
		},
		{
			"factory fails with error",
			"00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
			false,
			true,
			500,
		},
		{"missing traceparent", "", true, false, 200},
		{"incorrect traceparent", "Hello World", true, false, 200},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Arrange
			gin.SetMode("test")
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/", nil)

			if c.trcprnt != "" {
				req.Header.Add("traceparent", c.trcprnt)
			}

			ctx.Request = req

			factory := &MockTxFactory{c.txfFail}
			sut := TxContextMiddleware(factory)

			// Act
			sut(ctx)

			// Assert
			val, set := ctx.Get("tx-context")

			if w.Code != c.status {
				t.Errorf("expected a %d status code, but got %d", c.status, w.Code)
			}
			if c.txfFail && set {
				t.Error("expected trnsaction context to be unset but it's set")
			}
			if !c.txfFail && set {
				if !set {
					t.Error("expected trnsaction context to be set but it's unset")
				} else {
					tctx := val.(*TCtx)
					if tctx.ver == "" {
						t.Error("expected version to not be empty")
					}
					if tctx.tid == "" {
						t.Error("expected trace id to not be empty")
					}
					if c.decFail && tctx.pid != "" {
						t.Error("expected parent id to be empty")
					}
					if !c.decFail && tctx.pid == "" {
						t.Error("expected parent id to not be empty")
					}
					if tctx.flg == "" {
						t.Error("expected flag to not be empty")
					}
				}
			}
		})
	}
}
