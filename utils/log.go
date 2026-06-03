package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"
	"transaction_api/model/log_monitoring"
	"transaction_api/service"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type RequestExtractor struct {
	Message string `json:"message"`
	Request string `json:"username"`
}

func LogMonitoringMiddleware(service *service.LogMonitoringService) gin.HandlerFunc {

	return func(c *gin.Context) {
		startTime := time.Now()
		str := startTime.Format("2006-01-02 15:04:05")
		responseMsg := ""
		userName := ""

		var requestBodyBytes []byte
		if c.Request.Body != nil {
			requestBodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		var body map[string]interface{}
		if err := json.Unmarshal(requestBodyBytes, &body); err == nil {
			if username, ok := body["username"].(string); ok {
				userName = username
			}
		}

		if len(string(blw.body.String())) > 0 {
			var extractor RequestExtractor
			if err := json.Unmarshal(blw.body.Bytes(), &extractor); err == nil {
				if extractor.Message != "" {
					responseMsg = extractor.Message
				}
				if extractor.Request != "" {
					userName = extractor.Request
				}
			}
		}

		if userName == "" {
			userName = "guest"
		}

		oldValue, err := c.Get("oldValue")
		if !err {
			oldValue = ""
		}

		var parts []string
		if c.Request.Body != nil {
			parts = strings.Split(c.Request.URL.Path, "/")
		}
		var action = ""
		var menu = ""

		if len(parts) > 1 {
			menu = parts[1]
			action = parts[2]
		}

		logData := log_monitoring.Log_monitoring{
			ACTION:       strings.ToLower(action),
			MENU:         strings.ToLower(menu),
			OLD_VALUE:    oldValue.(string),
			NEW_VALUE:    string(requestBodyBytes),
			RESPONSE_MSG: responseMsg,
			CHANGED_BY:   userName,
			CHANGED_DATE: str,
			IP_CLIENT:    c.ClientIP(),
			USER_AGENT:   c.Request.Header.Get("User-Agent"),
		}

		go func() {

			test, insertErr := service.InsertLogMonitoringService(context.Background(), logData) // jangan panic-in request
			log.Printf("Log monitoring insert ID: %v", test)
			if insertErr != nil {
				log.Printf("Error inserting log monitoring: %v", insertErr)
			} else {

				LogToES(logData.MENU+"_"+logData.ACTION, logData.ACTION, logData.RESPONSE_MSG, map[string]interface{}{
					"full_request":  string(requestBodyBytes),
					"full_response": blw.body.String(),
					"message":       logData.RESPONSE_MSG,
					"user_id":       logData.CHANGED_BY,
					"ip_client":     logData.IP_CLIENT,
					"user_agent":    logData.USER_AGENT,
					"url":           c.Request.URL.Path,
				})

			}
		}()
	}
}
