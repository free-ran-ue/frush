package subscriber

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/free-ran-ue/frush/constant"
	"github.com/gin-gonic/gin"
)

var testGetConsoleLoginTokenCases = []struct {
	name         string
	templatePath string
	token        string
}{
	{
		name:         "test_token",
		templatePath: "../" + constant.TEMPLATE_CONSOLE_ACCOUNT_JSON,
		token:        "test_token",
	},
}

func TestGetConsoleLoginToken(t *testing.T) {
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	router.POST(constant.CONSOLE_LOGIN_PATH, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			constant.CONSOLE_ACCESS_TOKEN: "test_token",
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", constant.CONSOLE_IP, constant.CONSOLE_PORT),
		Handler: router,
	}

	errCh := make(chan error, 1)
	ready := make(chan struct{})
	go func() {
		ready <- struct{}{}
		errCh <- srv.ListenAndServe()
	}()

	<-ready
	time.Sleep(100 * time.Millisecond)

	for _, testCase := range testGetConsoleLoginTokenCases {
		t.Run(testCase.name, func(t *testing.T) {
			token, err := getConsoleLoginToken(testCase.templatePath)
			if err != nil {
				t.Fatalf("failed to get console login token: %v", err)
			}
			if token != testCase.token {
				t.Errorf("getConsoleLoginToken() = %v, want %v", token, testCase.token)
			}
		})
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}

	if err := <-errCh; err != nil && err != http.ErrServerClosed {
		t.Fatalf("Server error: %v", err)
	}
}
