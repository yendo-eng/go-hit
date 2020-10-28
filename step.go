package hit

import (
	"bytes"
	"fmt"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/errortrace"
)

type StepTime uint8

const (
	CombineStep StepTime = iota + 1
	CleanStep
	BeforeSendStep
	SendStep
	AfterSendStep
	BeforeExpectStep
	ExpectStep
	AfterExpectStep
)

func (s StepTime) String() string {
	switch s {
	case CleanStep:
		return "CleanStep"
	case CombineStep:
		return "CombineStep"
	case BeforeSendStep:
		return "BeforeSendStep"
	case SendStep:
		return "SendStep"
	case AfterSendStep:
		return "AfterSendStep"
	case BeforeExpectStep:
		return "BeforeExpectStep"
	case ExpectStep:
		return "ExpectStep"
	case AfterExpectStep:
		return "AfterExpectStep"
	}
	return ""
}

type IStep interface {
	trace() *errortrace.ErrorTrace
	when() StepTime
	callPath() callPath
	exec(instance *hitImpl) error
}

type hitStep struct {
	Trace    *errortrace.ErrorTrace
	When     StepTime
	CallPath callPath
	Exec     func(hit *hitImpl) error
}

func (step *hitStep) trace() *errortrace.ErrorTrace {
	return step.Trace
}

func (step *hitStep) when() StepTime {
	return step.When
}

func (step *hitStep) callPath() callPath {
	return step.CallPath
}

func (step *hitStep) exec(hit *hitImpl) (err error) {
	if step.Exec == nil {
		return nil
	}
	return step.Exec(hit)
}

func execStep(hit *hitImpl, step IStep) (err error) {
	setError := func(r interface{}) {
		if r == nil {
			return
		}

		setMeta := func() {
			step.trace().SetDescription(hit.Description())
			var b bytes.Buffer
			if newDebug(step.callPath().Push("Debug", nil), &b).exec(hit) == nil {
				step.trace().SetContext(b.String())
			}
		}

		switch v := r.(type) {
		case *errortrace.ErrorTrace:
			// this is already a errortrace
			err = v
			return
		case error:
			step.trace().SetError(v)
			setMeta()
			err = step.trace()
		default:
			step.trace().SetError(xerrors.New(fmt.Sprint(r)))
			setMeta()
			err = step.trace()
		}
	}

	defer func() {
		setError(recover())
	}()
	setError(step.exec(hit))
	return err
}

func StepCallPath(step IStep, withArguments bool) string {
	return step.callPath().CallString(withArguments)
}
