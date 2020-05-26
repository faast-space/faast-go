package log

import "github.com/sirupsen/logrus"

type wLogger struct {
	fields logrus.Fields
}

func For(fields logrus.Fields) logrus.FieldLogger {
	return &wLogger{fields}
}

func (l *wLogger) WithField(key string, value interface{}) *logrus.Entry {
	fs := logrus.Fields{
		key: value,
	}

	for k, v := range l.fields {
		fs[k] = v
	}

	return WithFields(fs)
}
func (l *wLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	for k, v := range l.fields {
		fields[k] = v
	}

	return WithFields(fields)
}
func (l *wLogger) WithError(err error) *logrus.Entry {
	return WithFields(l.fields).WithError(err)
}
func (l *wLogger) Debugf(format string, args ...interface{}) {
	WithFields(l.fields).Debugf(format, args...)
}
func (l *wLogger) Infof(format string, args ...interface{}) {
	WithFields(l.fields).Infof(format, args...)
}
func (l *wLogger) Printf(format string, args ...interface{}) {
	WithFields(l.fields).Printf(format, args...)
}
func (l *wLogger) Warnf(format string, args ...interface{}) {
	WithFields(l.fields).Warnf(format, args...)
}
func (l *wLogger) Warningf(format string, args ...interface{}) {
	WithFields(l.fields).Warningf(format, args...)
}
func (l *wLogger) Errorf(format string, args ...interface{}) {
	WithFields(l.fields).Errorf(format, args...)
}
func (l *wLogger) Fatalf(format string, args ...interface{}) {
	WithFields(l.fields).Fatalf(format, args...)
}
func (l *wLogger) Panicf(format string, args ...interface{}) {}
func (l *wLogger) Debug(args ...interface{})                 {}
func (l *wLogger) Info(args ...interface{})                  {}
func (l *wLogger) Print(args ...interface{})                 {}
func (l *wLogger) Warn(args ...interface{})                  {}
func (l *wLogger) Warning(args ...interface{})               {}
func (l *wLogger) Error(args ...interface{})                 {}
func (l *wLogger) Fatal(args ...interface{})                 {}
func (l *wLogger) Panic(args ...interface{})                 {}
func (l *wLogger) Debugln(args ...interface{})               {}
func (l *wLogger) Infoln(args ...interface{})                {}
func (l *wLogger) Println(args ...interface{})               {}
func (l *wLogger) Warnln(args ...interface{})                {}
func (l *wLogger) Warningln(args ...interface{})             {}
func (l *wLogger) Errorln(args ...interface{})               {}
func (l *wLogger) Fatalln(args ...interface{})               {}
func (l *wLogger) Panicln(args ...interface{})               {}
