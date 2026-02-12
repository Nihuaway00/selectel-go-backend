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
