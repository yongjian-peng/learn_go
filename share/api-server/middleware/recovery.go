package middleware

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"runtime"
	"share/common/pkg/appError"
	"share/common/pkg/config"
	"share/common/pkg/logger"

	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// ErrorHandler 业务异常返回
func ErrorHandler(ctx *fiber.Ctx, e error) error {
	// 业务异常
	if err, ok := e.(*appError.Error); ok {
		return ctx.JSON(map[string]interface{}{
			"code": err.Code,
			"msg":  err.Msg,
			"data": nil,
		})
	}

	if config.IsDevEnv() || config.IsTestEnv() {
		// 系统异常
		return ctx.JSON(map[string]interface{}{
			"code": appError.ServerError.Code,
			"msg":  e.Error(),
			"data": nil,
		})
	}

	// 系统异常
	return ctx.JSON(map[string]interface{}{
		"code": appError.ServerError.Code,
		"msg":  appError.ServerError.Msg,
		"data": nil,
	})
}

func Recover(isDevEnv bool) fiberRecover.Config {

	return fiberRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			//是否是开发环境
			if isDevEnv {
				fmt.Println(fmt.Sprintf("panic: %v\n%s\n", e, stack(4)))
			}
			// 如果不是业务异常，就抛出堆栈信息
			//if _, ok := e.(*appError.Error); !ok {
			logger.GetLogger("recover").Error("服务器异常：", zap.String("Stack", fmt.Sprintf("panic: %v\n%s\n", e, stack(4))))
			//}
		},
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; i < skip+12; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
