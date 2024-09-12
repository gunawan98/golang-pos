package log

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // Capture the response body in the buffer
	return w.ResponseWriter.Write(b)
}

func WrapHandler(f http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Clone the request body to avoid losing it after reading
		var reqBodyClone []byte
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				Logger.Error("Failed to read request body: " + err.Error())
			}
			reqBodyClone = bodyBytes

			// Restore the request body so the next handler can read it
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Use CustomResponseWriter to capture the response body
		crw := &CustomResponseWriter{
			ResponseWriter: res,
			body:           bytes.NewBuffer(nil),
		}

		// Call the next handler
		f.ServeHTTP(crw, req)

		// Log the request and response
		responseTime := time.Now()
		requestTime := time.Now()

		// set request response
		fields := []zapcore.Field{
			zap.String("url", req.URL.Path),
			zap.String("unique_id", uuid.New().String()),
			zap.String("request", string(reqBodyClone)),
			zap.String("response", crw.body.String()),
		}
		if req != nil {
			fields = append(fields, zap.String("request_time", requestTime.String()))
		}

		if res != nil {
			fields = append(fields, zap.String("response_time", responseTime.String()))
			processingTime := time.Since(requestTime)
			fields = append(fields, zap.Int("processing_time_nano_second", int(processingTime.Nanoseconds())))
		}
		Logger.Info("log global request and response", fields...)
	}
}
