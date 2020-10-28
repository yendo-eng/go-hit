// +build !generate

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //
package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// IClearExpectInt provides methods to clear steps.
type IClearExpectInt interface {
	IStep
	// Between clears all matching Between steps
	Between(value ...int) IStep
	// Equal clears all matching Equal steps
	Equal(value ...int) IStep
	// GreaterOrEqualThan clears all matching GreaterOrEqualThan steps
	GreaterOrEqualThan(value ...int) IStep
	// GreaterThan clears all matching GreaterThan steps
	GreaterThan(value ...int) IStep
	// LessOrEqualThan clears all matching LessOrEqualThan steps
	LessOrEqualThan(value ...int) IStep
	// LessThan clears all matching LessThan steps
	LessThan(value ...int) IStep
	// NotBetween clears all matching NotBetween steps
	NotBetween(value ...int) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...int) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...int) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...int) IStep
}
type clearExpectInt struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectInt(cp callPath) IClearExpectInt {
	return &clearExpectInt{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectInt) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectInt) when() StepTime {
	return CleanStep
}
func (v *clearExpectInt) callPath() callPath {
	return v.cp
}
func (v *clearExpectInt) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectInt) Between(value ...int) IStep {
	return removeStep(v.callPath().Push("Between", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) Equal(value ...int) IStep {
	return removeStep(v.callPath().Push("Equal", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) GreaterOrEqualThan(value ...int) IStep {
	return removeStep(v.callPath().Push("GreaterOrEqualThan", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) GreaterThan(value ...int) IStep {
	return removeStep(v.callPath().Push("GreaterThan", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) LessOrEqualThan(value ...int) IStep {
	return removeStep(v.callPath().Push("LessOrEqualThan", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) LessThan(value ...int) IStep {
	return removeStep(v.callPath().Push("LessThan", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) NotBetween(value ...int) IStep {
	return removeStep(v.callPath().Push("NotBetween", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) NotEqual(value ...int) IStep {
	return removeStep(v.callPath().Push("NotEqual", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) NotOneOf(value ...int) IStep {
	return removeStep(v.callPath().Push("NotOneOf", intSliceToInterfaceSlice(value)))
}
func (v *clearExpectInt) OneOf(value ...int) IStep {
	return removeStep(v.callPath().Push("OneOf", intSliceToInterfaceSlice(value)))
}
