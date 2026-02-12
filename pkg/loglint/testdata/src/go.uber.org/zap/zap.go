package zap

type Logger struct{}
type SugaredLogger struct{}

func (l *Logger) Debug(msg string, fields ...any)  {}
func (l *Logger) Info(msg string, fields ...any)   {}
func (l *Logger) Warn(msg string, fields ...any)   {}
func (l *Logger) Error(msg string, fields ...any)  {}
func (l *Logger) DPanic(msg string, fields ...any) {}
func (l *Logger) Panic(msg string, fields ...any)  {}
func (l *Logger) Fatal(msg string, fields ...any)  {}

func (l *SugaredLogger) Debug(args ...any)  {}
func (l *SugaredLogger) Info(args ...any)   {}
func (l *SugaredLogger) Warn(args ...any)   {}
func (l *SugaredLogger) Error(args ...any)  {}
func (l *SugaredLogger) DPanic(args ...any) {}
func (l *SugaredLogger) Panic(args ...any)  {}
func (l *SugaredLogger) Fatal(args ...any)  {}

func (l *SugaredLogger) Debugf(template string, args ...any)  {}
func (l *SugaredLogger) Infof(template string, args ...any)   {}
func (l *SugaredLogger) Warnf(template string, args ...any)   {}
func (l *SugaredLogger) Errorf(template string, args ...any)  {}
func (l *SugaredLogger) DPanicf(template string, args ...any) {}
func (l *SugaredLogger) Panicf(template string, args ...any)  {}
func (l *SugaredLogger) Fatalf(template string, args ...any)  {}

func (l *SugaredLogger) Debugw(msg string, keysAndValues ...any)  {}
func (l *SugaredLogger) Infow(msg string, keysAndValues ...any)   {}
func (l *SugaredLogger) Warnw(msg string, keysAndValues ...any)   {}
func (l *SugaredLogger) Errorw(msg string, keysAndValues ...any)  {}
func (l *SugaredLogger) DPanicw(msg string, keysAndValues ...any) {}
func (l *SugaredLogger) Panicw(msg string, keysAndValues ...any)  {}
func (l *SugaredLogger) Fatalw(msg string, keysAndValues ...any)  {}
