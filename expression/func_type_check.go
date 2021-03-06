//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"github.com/couchbase/query/value"
)

///////////////////////////////////////////////////
//
// IsArray
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISARRAY(expr).
It returns true if expr is an array; else false.
*/
type IsArray struct {
	UnaryFunctionBase
}

func NewIsArray(operand Expression) Function {
	rv := &IsArray{
		*NewUnaryFunctionBase("is_array", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsArray) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsArray) Type() value.Type { return value.BOOLEAN }

func (this *IsArray) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsArray) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsArray) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.ARRAY), nil
}

/*
Factory method pattern.
*/
func (this *IsArray) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsArray(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsAtom
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISATOM(expr).
Returns true if expr is a boolean, number, or string;
else false.
*/
type IsAtom struct {
	UnaryFunctionBase
}

func NewIsAtom(operand Expression) Function {
	rv := &IsAtom{
		*NewUnaryFunctionBase("is_atom", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsAtom) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsAtom) Type() value.Type { return value.BOOLEAN }

func (this *IsAtom) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsAtom) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

/*
Checks the type of input argument and returns true for boolean,
number and string and false for all other values.
*/
func (this *IsAtom) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	switch arg.Type() {
	case value.BOOLEAN, value.NUMBER, value.STRING:
		return value.TRUE_VALUE, nil
	default:
		return value.FALSE_VALUE, nil
	}
}

/*
Factory method pattern.
*/
func (this *IsAtom) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsAtom(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsBinary
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISBINARY(expr).
Returns true if expr is a boolean; else false.
*/
type IsBinary struct {
	UnaryFunctionBase
}

func NewIsBinary(operand Expression) Function {
	rv := &IsBinary{
		*NewUnaryFunctionBase("is_binary", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsBinary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsBinary) Type() value.Type { return value.BOOLEAN }

func (this *IsBinary) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsBinary) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsBinary) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.BINARY), nil
}

/*
Factory method pattern.
*/
func (this *IsBinary) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsBinary(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsBoolean
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISBOOLEAN(expr).
Returns true if expr is a boolean; else false.
*/
type IsBoolean struct {
	UnaryFunctionBase
}

func NewIsBoolean(operand Expression) Function {
	rv := &IsBoolean{
		*NewUnaryFunctionBase("is_boolean", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsBoolean) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsBoolean) Type() value.Type { return value.BOOLEAN }

func (this *IsBoolean) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsBoolean) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsBoolean) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.BOOLEAN), nil
}

/*
Factory method pattern.
*/
func (this *IsBoolean) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsBoolean(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsNumber
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISNUMBER(expr).
Returns true if expr is a number; else false.
*/
type IsNumber struct {
	UnaryFunctionBase
}

func NewIsNumber(operand Expression) Function {
	rv := &IsNumber{
		*NewUnaryFunctionBase("is_number", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsNumber) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsNumber) Type() value.Type { return value.BOOLEAN }

func (this *IsNumber) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsNumber) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsNumber) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.NUMBER), nil
}

/*
Factory method pattern.
*/
func (this *IsNumber) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsNumber(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsObject
//
///////////////////////////////////////////////////

/*
This represents the type checking function ISOBJECT(expr).
Returns true if expr is an object; else false.
*/
type IsObject struct {
	UnaryFunctionBase
}

func NewIsObject(operand Expression) Function {
	rv := &IsObject{
		*NewUnaryFunctionBase("is_object", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsObject) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsObject) Type() value.Type { return value.BOOLEAN }

func (this *IsObject) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsObject) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsObject) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.OBJECT), nil
}

/*
Factory method pattern.
*/
func (this *IsObject) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsObject(operands[0])
	}
}

///////////////////////////////////////////////////
//
// IsString
//
///////////////////////////////////////////////////

/*
This represents the Type checking function ISSTRING(expr).
Returns true if expr is a string; else false.
*/
type IsString struct {
	UnaryFunctionBase
}

func NewIsString(operand Expression) Function {
	rv := &IsString{
		*NewUnaryFunctionBase("is_string", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsString) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *IsString) Type() value.Type { return value.BOOLEAN }

func (this *IsString) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For boolean functions, simply list this expression.
*/
func (this *IsString) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

func (this *IsString) Apply(context Context, arg value.Value) (value.Value, error) {
	if arg.Type() == value.MISSING || arg.Type() == value.NULL {
		return arg, nil
	}

	return value.NewValue(arg.Type() == value.STRING), nil
}

/*
Factory method pattern.
*/
func (this *IsString) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsString(operands[0])
	}
}

///////////////////////////////////////////////////
//
// Type
//
///////////////////////////////////////////////////

/*
This represents the type checking function TYPE(expr).
Returns the type based on the value of the expr as a string.
*/
type Type struct {
	UnaryFunctionBase
}

func NewType(operand Expression) Function {
	rv := &Type{
		*NewUnaryFunctionBase("type", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *Type) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunction(this)
}

func (this *Type) Type() value.Type { return value.STRING }

func (this *Type) Evaluate(item value.Value, context Context) (value.Value, error) {
	return this.UnaryEval(this, item, context)
}

func (this *Type) PropagatesMissing() bool {
	return false
}

func (this *Type) PropagatesNull() bool {
	return false
}

func (this *Type) Apply(context Context, arg value.Value) (value.Value, error) {
	return value.NewValue(arg.Type().String()), nil
}

/*
Factory method pattern.
*/
func (this *Type) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewType(operands[0])
	}
}
