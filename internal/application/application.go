package application

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"simplesurance-test-task/internal/counter"
	"simplesurance-test-task/internal/limiter"
	"simplesurance-test-task/internal/storage"
	"simplesurance-test-task/internal/webserver"
)

type Application struct {
	ctx             context.Context
	cancel          context.CancelFunc
	state           uint64
	wg              sync.WaitGroup
	shutdownTimeout time.Duration
	webserver       *webserver.Server
	logger          *slog.Logger
}

func New() *Application {
	return &Application{
		shutdownTimeout: DefaultShutdownTimeout,
	}
}

func (app *Application) Init() (err error) {
	app.ctx = context.Background()
	app.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	config := NewConfig()
	config.LoadFromEnv()

	requestCounter, err := counter.New(storage.NewFileStorage(config.Filename))
	if err != nil {
		return err
	}

	parallelLimiter := limiter.NewParallelRateLimiter(config.MaxParallelRequests)

	if app.webserver, err = webserver.New(
		``,
		config.Port,
		app.logger,
		requestCounter,
		parallelLimiter,
	); err != nil {
		return err
	}

	app.logger.Info(`application is initialized`)

	return nil
}

func (app *Application) Run() error {
	for {
		state := atomic.LoadUint64(&app.state)
		if state == running {
			return ErrAppRunning
		}
		if atomic.CompareAndSwapUint64(&app.state, state, running) {
			break
		}
		runtime.Gosched()
	}

	idleConnectionsClosed := make(chan struct{})
	ctx := context.Background()
	go app.signalHandler(ctx, idleConnectionsClosed)

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		if err := app.webserver.Serve(); err != nil {
			app.logger.Error(`webserver error`, `error`, err.Error())
		}
		<-idleConnectionsClosed
	}()

	app.logger.Info(`application is running`)
	app.wg.Wait()
	app.logger.Info(`application stopped`)

	return nil
}

func (app *Application) Shutdown() error {
	for {
		state := atomic.LoadUint64(&app.state)
		if state != running {
			return ErrAppNotRunning
		}
		if atomic.CompareAndSwapUint64(&app.state, state, none) {
			break
		}
		runtime.Gosched()
	}

	app.logger.Info(`application stops`)

	app.ctx, app.cancel = context.WithTimeout(app.ctx, app.shutdownTimeout)
	defer func() {
		app.cancel()
	}()

	return app.webserver.Shutdown(app.ctx)
}

func (app *Application) signalHandler(ctx context.Context, idleConnectionsClosed chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, os.Kill)
	sig := <-ch
	app.logger.Info(`OS signal received`, `signal`, sig.String())

	if err := app.Shutdown(); err != nil {
		app.logger.Error(`failed to shutdown app`, `error`, err.Error())
	}
	close(idleConnectionsClosed)
}

var (
	ErrAppRunning    = errors.New(`application is already running`)
	ErrAppNotRunning = errors.New(`application is not running `)
)

const (
	DefaultShutdownTimeout = 60 * time.Second

	none = uint64(1) << iota
	running
)
