package hit

type IExpectBody interface {
	IStep
	JSON(data ...interface{}) IExpectBodyJSON
	Equal(data interface{}) IStep
	NotEqual(data interface{}) IStep
	Contains(data interface{}) IStep
	NotContains(data interface{}) IStep
}

type expectBody struct {
	expect IExpect
}

func newExpectBody(expect IExpect) IExpectBody {
	return &expectBody{expect}
}

func (body *expectBody) when() StepTime {
	return body.expect.when()
}

func (body *expectBody) exec(hit Hit) error {
	return body.expect.exec(hit)
}

// JSON expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect().Body().JSON([]int{1, 2, 3})
//           Expect().Body().JSON().Contains(1)
func (body *expectBody) JSON(data ...interface{}) IExpectBodyJSON {
	if arg, ok := getLastArgument(data); ok {
		return finalExpectBodyJSON{newExpectBodyJSON(body.expect).Equal("", arg)}
	}
	return newExpectBodyJSON(body.expect)
}

// Compare functions

// Equal expects the body to be equal to the specified value
// Example:
//           Expect().Body().Equal("Hello World")
func (body *expectBody) Equal(data interface{}) IStep {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.equalOnlyNativeTypes(data, true) {
			return
		}
		hit.Expect().Body().JSON().Equal("", data)
	})
}

// NotEqual expects the body to be not equal to the specified value
// Example:
//           Expect().Body().NotEqual("Hello World")
func (body *expectBody) NotEqual(data interface{}) IStep {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.equalOnlyNativeTypes(data, false) {
			return
		}
		hit.Expect().Body().JSON().NotEqual("", data)
	})
}

// Contains expects the body to contains the specified value
// Example:
//           Expect().Body().Contains("Hello World")
func (body *expectBody) Contains(data interface{}) IStep {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.containsOnlyNativeTypes(data, true) {
			return
		}
		hit.Expect().Body().JSON().Contains("", data)
	})
}

// NotContains expects the body to not contain the specified value
// Example:
//           Expect().Body().NotContains("Hello World")
func (body *expectBody) NotContains(data interface{}) IStep {
	return body.expect.Custom(func(hit Hit) {
		if hit.Response().body.containsOnlyNativeTypes(data, false) {
			return
		}
		hit.Expect().Body().JSON().NotContains("", data)
	})
}

type finalExpectBody struct {
	IStep
}

func (f finalExpectBody) JSON(data ...interface{}) IExpectBodyJSON {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Equal(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotEqual(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) Contains(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
func (f finalExpectBody) NotContains(data interface{}) IStep {
	panic("only usable with Expect().Body() not with Expect().Body(value)")
}
