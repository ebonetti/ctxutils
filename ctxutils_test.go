package ctxutils

import (
	"context"
	"errors"
	"testing"
)

var errCorr = errors.New("Dummy Error")
var errWrong = errors.New("Wrong Error")

func TestWithFail(t *testing.T) {
	t.Run("NilFail", func(t *testing.T) {
		_, fail := WithFail(context.Background())
		err := fail(nil)
		if err != nil {
			t.Error("fail(nil) should return nil, instead it returns ", err)
		}
	})
	t.Run("CancelledFail", func(t *testing.T) {
		parent, cancel := context.WithCancel(context.Background())
		cancel()
		_, fail := WithFail(parent)
		err := fail(errWrong)
		if parent.Err() != err {
			t.Error("cancel followed by fail should return ", parent.Err(), ", instead it returns ", err)
		}
	})
	t.Run("CancelFail", func(t *testing.T) {
		parent, cancel := context.WithCancel(context.Background())
		_, fail := WithFail(parent)
		cancel()
		err := fail(errWrong)
		if parent.Err() != err {
			t.Error("cancel followed by fail should return ", parent.Err(), ", instead it returns ", err)
		}
	})
	t.Run("FailFail", func(t *testing.T) {
		_, fail := WithFail(context.Background())
		err := fail(errCorr)
		err1 := fail(errWrong)
		switch {
		case err != errCorr:
			t.Error("fail(errCorr) should return errCorr, instead it returns ", err)
		case err1 != errCorr:
			t.Error("fail(errCorr) followed by fail should return errCorr, instead it returns ", err1)
		}
	})
	t.Run("FailCancel", func(t *testing.T) {
		parent, cancel := context.WithCancel(context.Background())
		_, fail := WithFail(parent)
		fail(errCorr)
		cancel()
		err := fail(errWrong)
		if err != errCorr {
			t.Error("fail followed by cancel followed by fail should return errCorr, instead it returns ", err)
		}
	})
}
