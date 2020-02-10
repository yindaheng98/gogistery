package Registry

import (
	"fmt"
	"gogistery/Registry"
	"testing"
	"time"
)
import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func TestRegistry(t *testing.T) {
	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	beatProto := NewResponseBeatProtocol()
	registry := Registry.New(Info{"0", "0.registry"}, 5, NewTimeoutProtocol(), beatProto)
	go registry.Run()
	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		RegistryHandler(ctx, beatProto)
	})
	go func() {
		for {
			fmt.Print(registry.GetConnections())
			fmt.Print("\n")
			time.Sleep(1e9)
		}
	}()

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
