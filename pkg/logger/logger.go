package logger

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/rs/zerolog"
)

func Init() (zerolog.Logger, error) {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return zerolog.Logger{}, err
	}

	file, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return zerolog.Logger{}, err
	}

	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	logger := zerolog.New(multi).With().Caller().Timestamp().Logger()

	return logger, nil
}

func NewLogger(baseLogger zerolog.Logger, component string) zerolog.Logger {
	return baseLogger.With().Str("COMPONENT", component).Logger()
}

func LogWithParams(logger zerolog.Logger, msg string, params interface{}) {
	event := logger.Info()

	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		logStructFields(event, val)
	case reflect.Slice:
		logSliceValues(event, val)
	default:
		logSingleValue(event, "value", val.Interface())
	}

	event.Msg(msg)
}

func logStructFields(event *zerolog.Event, val reflect.Value) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i).Interface()

		event = logSingleValue(event, field.Name, fieldValue)
	}
}

func logSliceValues(event *zerolog.Event, val reflect.Value) {
	for i := 0; i < val.Len(); i++ {
		event = logSingleValue(event, fmt.Sprintf("[%d]", i), val.Index(i).Interface())
	}
}

func logSingleValue(event *zerolog.Event, key string, value interface{}) *zerolog.Event {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		event = event.Str(key, val.String())
	case reflect.Bool:
		event = event.Bool(key, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		event = event.Int64(key, val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		event = event.Uint64(key, val.Uint())
	case reflect.Float32, reflect.Float64:
		event = event.Float64(key, val.Float())
	default:
		event = event.Interface(key, value)
	}

	return event
}
