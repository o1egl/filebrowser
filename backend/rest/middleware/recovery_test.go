package middleware

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/filebrowser/filebrowser/v3/log"
)

func TestMain(m *testing.M) {
	oldLogger := log.DefaultLogger
	defer func() {
		log.DefaultLogger = oldLogger
	}()
	os.Exit(m.Run())
}

type logWriter struct {
	bytes.Buffer
}

func (w *logWriter) Sync() error {
	return nil
}

type RecoveryTestSuite struct {
	suite.Suite
	logOut *logWriter
}

func TestRecoverySuite(t *testing.T) {
	suite.Run(t, new(RecoveryTestSuite))
}

func (suite *RecoveryTestSuite) SetupTest() {
	suite.logOut = &logWriter{}
	var err error
	log.DefaultLogger, err = log.NewLogger(log.Configuration{
		LogLevel: log.LevelDebug,
		Format:   log.FormatPlain,
		Output:   suite.logOut,
	})
	if err != nil {
		panic(err)
	}
}

// TestPanicInHandler assert that panic has been recovered.
func (suite *RecoveryTestSuite) TestPanicInHandler() {
	router := gin.New()
	router.Use(Recovery)
	router.GET("/recovery", func(_ *gin.Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(suite.logOut.String(), "panic recovered")
	suite.Contains(suite.logOut.String(), "Oupps, Houston, we have a problem")
	suite.Contains(suite.logOut.String(), "TestPanicInHandler")
	suite.NotContains(suite.logOut.String(), "GET /recovery")

	// Debug mode prints the request
	gin.SetMode(gin.DebugMode)
	// RUN
	w = performRequest(router, "GET", "/recovery")
	// TEST
	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(suite.logOut.String(), "/recovery")

	gin.SetMode(gin.TestMode)
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func (suite *RecoveryTestSuite) TestPanicWithAbort() {
	router := gin.New()
	router.Use(Recovery)
	router.GET("/recovery", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	// TEST
	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *RecoveryTestSuite) TestSource() {
	bs := source(nil, 0)
	suite.Equal([]byte("???"), bs)

	in := [][]byte{
		[]byte("Hello world."),
		[]byte("Hi, gin.."),
	}
	bs = source(in, 10)
	suite.Equal([]byte("???"), bs)

	bs = source(in, 1)
	suite.Equal([]byte("Hello world."), bs)
}

func (suite *RecoveryTestSuite) TestFunction() {
	bs := function(1)
	suite.Equal([]byte("???"), bs)
}

// TestPanicWithBrokenPipe asserts that recovery specifically handles
// writing responses to broken pipes
func (suite *RecoveryTestSuite) TestPanicWithBrokenPipe() {
	const expectCode = 204

	expectMsgs := map[syscall.Errno]string{
		syscall.EPIPE:      "broken pipe",
		syscall.ECONNRESET: "connection reset by peer",
	}

	for errno, expectMsg := range expectMsgs {
		suite.Run(expectMsg, func() {
			router := gin.New()
			router.Use(Recovery)
			router.GET("/recovery", func(c *gin.Context) {
				// Start writing response
				c.Header("X-Test", "Value")
				c.Status(expectCode)

				// Oops. Client connection closed
				e := &net.OpError{Err: &os.SyscallError{Err: errno}}
				panic(e)
			})
			// RUN
			w := performRequest(router, "GET", "/recovery")
			// TEST
			suite.Equal(expectCode, w.Code)
			suite.Contains(strings.ToLower(suite.logOut.String()), expectMsg)
		})
	}
}

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
