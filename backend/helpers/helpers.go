package helpers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func EncodeBase62(n uint64) string {
	var result string = ""
	for n > 0 {
		remainder := n % 62
		result = string(alphabet[remainder]) + result
		n = n / 62
	}
	return result
}

func DecodeBase62(encodedStr string) (uint64, error) {
	var n uint64 = 0
	for i := 0; i < len(encodedStr); i++ {
		char := encodedStr[i]
		index := strings.IndexByte(alphabet, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid character '%c' in base62 string", char)
		}
		n = n*uint64(62) + uint64(index)
	}
	return n, nil
}

func SetupPrometheus(metricsEndpoint string, server *gin.Engine) {
	// get global Monitor object
	monitor := ginmetrics.GetMonitor()

	monitor.SetMetricPath(metricsEndpoint)

	// set middleware for gin
	monitor.Use(server)
}
