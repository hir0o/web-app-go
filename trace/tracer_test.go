package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Errorf("New() returned nil")
	} else {
		tracer.Trace("今日わ")
		if buf.String() != "今日わ\n" {
			t.Errorf("'%s'という文字列が出力されました。", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = off()
	silentTracer.Trace("データ")
}