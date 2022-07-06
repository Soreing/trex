package trex

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Injects a transaction context dependency into the gin context
// The transaction context should contain traced dependencies
// such as dbs and https clients
func TxContextMiddleware(factory TxFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get and decode the traceparent header into smaller values
		trcprnt := c.GetHeader("traceparent")
		ver, tid, pid, flg, err := DecodeTraceparent(trcprnt)

		// If the header could not be decoded, generate a new header
		if err != nil {
			ver, flg = "00", "01"
			if tid, err = GenerateRadomHexString(16); err != nil {
				c.AbortWithError(
					http.StatusInternalServerError,
					fmt.Errorf("failed to generate trace id %w", err),
				)
				return
			}
		}

		// Generate a new resource id
		rid, err := GenerateRadomHexString(8)
		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				fmt.Errorf("failed to generate resource id %w", err),
			)
			return
		}

		// Generate a transaction context usin the factory
		txc, err := factory.Generate(ver, tid, pid, rid, flg)
		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				fmt.Errorf("failed to generate transaction context %w", err),
			)
			return
		}

		c.Set("tx-context", txc)
		c.Next()
	}
}
