package global

import (
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"sync"
	"sync/atomic"
)

// loggerProvider is a placeholder for a configured SDK LoggerProvider.
//
// All LoggerProvider functionality is forwarded to a delegate once
// configured.
type loggerProvider struct {
	mtx      sync.Mutex
	loggers  map[il]*logger
	delegate logs.LoggerProvider
}

// Compile-time guarantee that loggerProvider implements the LoggerProvider
// interface.
var _ logs.LoggerProvider = &loggerProvider{}

// setDelegate configures p to delegate all TracerProvider functionality to
// provider.
//
// All Tracers provided prior to this function call are switched out to be
// Tracers provided by provider.
//
// It is guaranteed by the caller that this happens only once.
func (p *loggerProvider) setDelegate(provider logs.LoggerProvider) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.delegate = provider

	if len(p.loggers) == 0 {
		return
	}

	for _, t := range p.loggers {
		t.setDelegate(provider)
	}

	p.loggers = nil
}

// Logger implements LoggerProvider.
func (p *loggerProvider) Logger(name string, opts ...logs.LoggerOption) logs.Logger {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.delegate != nil {
		return p.delegate.Logger(name, opts...)
	}

	// At this moment it is guaranteed that no sdk is installed, save the tracer in the tracers map.

	c := logs.NewLoggerConfig(opts...)
	key := il{
		name:    name,
		version: c.InstrumentationVersion(),
	}

	if p.loggers == nil {
		p.loggers = make(map[il]*logger)
	}

	if val, ok := p.loggers[key]; ok {
		return val
	}

	t := &logger{name: name, opts: opts, provider: p}
	p.loggers[key] = t
	return t
}

type il struct {
	name    string
	version string
}

// logger is a placeholder for a logs.Logger.
//
// All Logger functionality is forwarded to a delegate once configured.
// Otherwise, all functionality is forwarded to a NoopLogger.
type logger struct {
	name     string
	opts     []logs.LoggerOption
	provider *loggerProvider

	delegate atomic.Value
}

// Compile-time guarantee that logger implements the logs.Logger interface.
var _ logs.Logger = &logger{}

func (t *logger) Emit(logRecord logs.LogRecord) {
	delegate := t.delegate.Load()
	if delegate != nil {
		delegate.(logs.Logger).Emit(logRecord)
	}
}

// setDelegate configures t to delegate all Tracer functionality to Loggers
// created by provider.
//
// All subsequent calls to the Logger methods will be passed to the delegate.
//
// It is guaranteed by the caller that this happens only once.
func (t *logger) setDelegate(provider logs.LoggerProvider) {
	t.delegate.Store(provider.Logger(t.name, t.opts...))
}
