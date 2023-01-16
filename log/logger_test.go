package log

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSimpleLogger(t *testing.T) {
	// ------------INIT--------------------
	logger := New("PROJECTKEY", "LOCAL")

	logger.SetPathURI("/hello")                                // Log for where api route is calling from
	logger.SetAPIKey("aosidufpaisudfojasdofijapsodfijapsdjfa") //Log for what API key is used [Masked befor log]
	logger.SetToken("TOKEN_STRING")                            // Log for what is token string
	logger.SetSpanID("SPANID")                                 // Log for Span ID for trace log
	logger.SetTraceID("TRACEID")                               // Log for Trace ID for Trace log

	// ------------USE--------------------
	logger.Warn("Hello world Info 2", "")  // Log msg,info with Warn Log Level
	logger.Error("Hello world Info 2", "") // Log msg,info with Error Log Level

	logger.Access("Request Begin", "")                    // Log msg,info with Infor Log Level will have separate index on server
	logger.SetResponseTime(20).Access("Response End", "") // Log msg,info with Infor Log Level will have separate index on server

}

// v0.0.0-20221025-db424b6f
func TestLogger(t *testing.T) {
	//init logger
	logger := New("PROJECTKEY", "LOCA", fmt.Sprintf("%vlogbeat%vlog", string(os.PathSeparator), string(os.PathSeparator))) // path params not work if env. are DEVELOP /UAT /PRODUCTION

	//set essential information dot function
	logger.SetPathURI("/hello").SetAPIKey("aosidufpaisudfojasdofijapsodfijapsdjfa").SetToken("TOKEN_STRING").SetSpanID("SPANID").SetTraceID("TRACEID")

	//set essential information individual function equal to above line
	logger.SetPathURI("/hello")                                // Log for where api route is calling from
	logger.SetAPIKey("aosidufpaisudfojasdofijapsodfijapsdjfa") //Log for what API key is used [Masked befor log]
	logger.SetToken("TOKEN_STRING")                            // Log for what is token string
	logger.SetSpanID("SPANID")                                 // Log for Span ID for trace log
	logger.SetTraceID("TRACEID")                               // Log for Trace ID for Trace log

	//save log to file
	logger.Trace("Hello world Info 2", "") // Log msg,info with Trace Log Level
	logger.Debug("Hello world Info", "")   // Log msg,info with Debug Log Level
	logger.Info("Hello world Info 2", "")  // Log msg,info with Info Log Level
	logger.Warn("Hello world Info 2", "")  // Log msg,info with Warn Log Level
	logger.Error("Hello world Info 2", "") // Log msg,info with Error Log Level

	// //save access log to file / response with setresponsetime
	logger.Access("Request Begin", "")                    // Log msg,info with Infor Log Level will have separate index on server
	logger.SetResponseTime(20).Access("Response End", "") // Log msg,info with Infor Log Level will have separate index on server

	// !!OPTIONAL  Log path prefix!!

	// Only works if env. not equal DEVELOP, SIT, UAT, PRODUCTION Default Value is always "/logbeat"
	customlogpathPrefix := fmt.Sprintf("%vcustom%vlog", string(os.PathSeparator), string(os.PathSeparator))
	logger2 := New("PROJECTKEY", "LOCAL", customlogpathPrefix)
	logger2.SetToken("xxxxxxxxxxxxx")
}

func TestRemovePath(t *testing.T) {
	path := removeTrailSlash("", "")
	assert.Equal(t, "", path)

	path = removeTrailSlash("/custom/log", "")
	assert.Equal(t, "/custom/log", path)

	path = removeTrailSlash(" /custom/log / ", "")
	assert.Equal(t, "/custom/log", path)

}

func TestRemovePathShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code should have panic")
		}
	}()

	removeTrailSlash(" / custom/log / ", "") // Log msg,info with
}
