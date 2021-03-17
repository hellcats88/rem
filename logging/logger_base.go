package logging

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hellcats88/abstracte/logging"
)

func defaultExtraParameterFormat(extras []logging.K) string {
	return fmt.Sprintf("%+v", extras)
}

type WriteMessage func(level logging.Level, msg string)

type base struct {
	config               logging.Config
	writeMessageCallback WriteMessage
}

func New(config logging.Config, writeMessageCallback WriteMessage) logging.Logger {
	log := base{
		config:               config,
		writeMessageCallback: writeMessageCallback,
	}

	if len(log.config.Order) == 0 {
		log.config.Order = append(log.config.Order, logging.PartOrderLevel, logging.PartOrderTimestamp,
			logging.PartOrderCorrelationId, logging.PartOrderMessage, logging.PartOrderExtra)
	}

	if log.config.ExtraParametersFormat == nil {
		log.config.ExtraParametersFormat = defaultExtraParameterFormat
	}

	if !log.config.SkipExtraParameterPrefix {
		if log.config.ExtraParametersPrefix == "" {
			log.config.ExtraParametersPrefix = "extra-parameters="
		}
	}

	if log.config.CustomTime == nil {
		log.config.CustomTime = func() (string, time.Time) {
			now := time.Now()
			if log.config.TimeFormat != "" {
				return now.Format(log.config.TimeFormat), now
			} else {
				return now.String(), now
			}
		}
	}

	return log
}

func (cns base) composeMessage(level string, ctx logging.Context, message string) string {
	msg := ""

	for _, part := range cns.config.Order {
		switch part {
		case logging.PartOrderLevel:
			if !cns.config.SkipPrintLevel {
				msg += level + ":"
			}

		case logging.PartOrderTimestamp:
			if !cns.config.SkipPrintTimestamp {
				str, _ := cns.config.CustomTime()
				msg += "[" + str + "]:"
			}

		case logging.PartOrderCorrelationId:
			if !cns.config.SkipPrintCorrelationID {
				msg += "[" + ctx.CorrID() + "]:"
			}

		case logging.PartOrderMessage:
			msg += message

		case logging.PartOrderExtra:
			if len(ctx.GetExtras()) > 0 {
				if msg != "" {
					msg += ":" + cns.config.ExtraParametersPrefix + cns.config.ExtraParametersFormat(ctx.GetExtras())
				} else {
					msg += cns.config.ExtraParametersPrefix + cns.config.ExtraParametersFormat(ctx.GetExtras())
				}
			}

		default:
			break
		}
	}

	return msg[:len(msg)-1]
}

func (cns base) _print(ctx logging.Context, referenceLevel logging.Level, levelText string, msg string, params ...interface{}) {
	if referenceLevel <= cns.config.Level {
		if cns.config.CustomLogFormat != nil {
			_, now := cns.config.CustomTime()
			cns.writeMessageCallback(referenceLevel, cns.config.CustomLogFormat(logging.CustomLogFormatData{
				Level:         referenceLevel,
				CorrelationID: ctx.CorrID(),
				CurrentTime:   now,
				Message:       fmt.Sprintf(msg, params...),
				ExtraParams:   ctx.GetExtras(),
			}))
		} else {
			cns.writeMessageCallback(referenceLevel, cns.composeMessage(levelText, ctx, fmt.Sprintf(msg, params...)))
		}
	}
}

func (cns base) Debug(ctx logging.Context, msg string, params ...interface{}) {
	cns._print(ctx, logging.Debug, "DEBUG", msg, params...)
}

func (cns base) Trace(ctx logging.Context, msg string, params ...interface{}) {
	cns._print(ctx, logging.Trace, "TRACE", msg, params...)
}

func (cns base) Error(ctx logging.Context, msg string, params ...interface{}) {
	cns._print(ctx, logging.Error, "ERROR", msg, params...)
}

func (cns base) Info(ctx logging.Context, msg string, params ...interface{}) {
	cns._print(ctx, logging.Info, "INFO", msg, params...)
}

func (cns base) Warn(ctx logging.Context, msg string, params ...interface{}) {
	cns._print(ctx, logging.Warn, "WARN", msg, params...)
}

func (cns base) BeginMethod(ctx logging.Context) {
	if logging.Debug <= cns.config.Level {
		fpcs := make([]uintptr, 1)
		runtime.Callers(2, fpcs)
		caller := runtime.FuncForPC(fpcs[0] - 1)
		cns.Debug(ctx, "Begin "+caller.Name())
	}
}

func (cns base) BeginMethodParams(ctx logging.Context, format string, params ...interface{}) {
	if logging.Debug <= cns.config.Level {
		fpcs := make([]uintptr, 1)
		runtime.Callers(2, fpcs)
		caller := runtime.FuncForPC(fpcs[0] - 1)
		cns.Debug(ctx, "Begin "+caller.Name()+" "+format, params...)
	}
}

func (cns base) EndMethod(ctx logging.Context) {
	if logging.Debug <= cns.config.Level {
		fpcs := make([]uintptr, 1)
		runtime.Callers(2, fpcs)
		caller := runtime.FuncForPC(fpcs[0] - 1)
		cns.Debug(ctx, "End "+caller.Name())
	}
}

func (cns base) EndMethodParams(ctx logging.Context, format string, params ...interface{}) {
	if logging.Debug <= cns.config.Level {
		fpcs := make([]uintptr, 1)
		runtime.Callers(2, fpcs)
		caller := runtime.FuncForPC(fpcs[0] - 1)
		cns.Debug(ctx, "End "+caller.Name()+" "+format, params...)
	}
}
