package examplelogx

import (
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"testing"
	"time"
)

func TestSample(t *testing.T) {
	type log logx.Logger
	var sample log = logrusx.New()
	sample.Infof("hello")
}

func TestSampleJson(t *testing.T) {
	var Logger logx.Logger
	var formatter logrusx.JsonFormatter
	Logger = logrusx.New(logrusx.WithFormatter(formatter))
	Logger.Infof("format json")
}

func TestLogFiled(t *testing.T) {
	var Logger logx.Logger
	var formatter logrusx.JsonFormatter
	Logger = logrusx.New(logrusx.WithFormatter(formatter))
	today := time.Now().String()
	Logger.WithFieldsX(Logger.Info, "your_date_module", "today", today)
}
