package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"thor/src/handler"
	"thor/src/properties"
	"thor/src/server/container"
	"thor/src/server/middleware"
	"thor/src/server/router"
	"time"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func StartHttpServer(container *container.DefaultContainer) {
	e := echo.New()

	tracer.Start(tracer.WithAgentAddr("datadog-agent:8126"))

	err := profiler.Start(
		profiler.WithService("thor-service"),
		profiler.WithEnv("Dev"),
		profiler.WithVersion("7.0"),
		profiler.WithTags("org:flux", "org:titipaja"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			// The profiles below are disabled by default to keep overhead
			// low, but can be enabled as needed.

			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	//e.Use(echomiddleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	//	c.Logger().Info(reqBody)
	//	c.Logger().Info(resBody)
	//}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &properties.CustomContext{c}
			return next(cc)
		}
	})
	e.Use(echomiddleware.Logger())
	e.Logger.SetLevel(log.INFO)
	validate := validator.New()
	validate.RegisterValidation("ISO8601date", IsISO8601Date)

	e.Validator = &CustomValidator{validator: validate}

	middleware.SetupMiddleware(e)
	router.InitializeRouter(e, handler.InitializeHandler(container))

	port := fmt.Sprintf("%s%s", ":", container.Config.Server.Port)
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	defer profiler.Stop()
	defer tracer.Stop()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func IsISO8601Date(fl validator.FieldLevel) bool {

	if fl.Field().String() == "" {
		return true
	}

	ISO8601DateRegexString := "^\\d{4}(-\\d\\d(-\\d\\d(T\\d\\d:\\d\\d(:\\d\\d)?(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?)?)?)?$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}
