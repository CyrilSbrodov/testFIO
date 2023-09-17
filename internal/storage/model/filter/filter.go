package filter

import (
	"fmt"

	"testFIO/cmd/loggers"
)

const (
	DataTypeStr = "str"
	DataTypeInt = "int"

	OperatorEq            = "eq"
	OperatorNotEq         = "neq"
	OperatorLowerThan     = "lt"
	OperatorLowerThanEq   = "lte"
	OperatorGreaterThan   = "gt"
	OperatorGreaterThanEq = "gte"
	OperatorBetween       = "between"
	OperatorSubString     = "like"
	Eq                    = "="
	NotEq                 = "!="
	LowerThan             = "<"
	LowerThanEq           = "<="
	GreaterThan           = ">"
	GreaterThanEq         = ">="
)

type options struct {
	operators map[string]string
	limit     int
	fields    []Field
	logger    *loggers.Logger
}

func NewOptions(limit int) Options {
	logger := loggers.NewLogger()
	operators := addOperatorsToMap()
	return &options{
		operators: operators,
		limit:     limit,
		fields:    nil,
		logger:    logger,
	}
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

type Options interface {
	Limit() int
	AddField(name, operator, value, dType string) error
	Fields() []Field
}

func (o *options) Limit() int {
	return o.limit
}

func (o *options) AddField(name, operator, value, dType string) error {
	err := o.validateOperator(operator)
	if err != nil {
		o.logger.LogErr(err, "bad operator")
		return err
	}
	oper := o.selectOperator(operator)
	o.fields = append(o.fields, Field{
		Name:     name,
		Value:    value,
		Operator: oper,
		Type:     dType,
	})
	return nil
}

func (o *options) Fields() []Field {
	return o.fields
}

func (o *options) validateOperator(operator string) error {
	if _, ok := o.operators[operator]; !ok {
		return fmt.Errorf("bad operator")
	}
	return nil
}

func addOperatorsToMap() map[string]string {
	var operators = make(map[string]string)
	operators[OperatorEq] = OperatorEq
	operators[OperatorNotEq] = OperatorNotEq
	operators[OperatorLowerThan] = OperatorLowerThan
	operators[OperatorLowerThanEq] = OperatorLowerThanEq
	operators[OperatorGreaterThan] = OperatorGreaterThan
	operators[OperatorGreaterThanEq] = OperatorGreaterThanEq
	operators[OperatorBetween] = OperatorBetween
	operators[OperatorSubString] = OperatorSubString
	return operators
}

func (o *options) selectOperator(operator string) string {
	newOperator := ""
	switch operator {
	case OperatorEq:
		newOperator = Eq
	case OperatorNotEq:
		newOperator = NotEq
	case OperatorLowerThan:
		newOperator = LowerThan
	case OperatorLowerThanEq:
		newOperator = LowerThanEq
	case OperatorGreaterThan:
		newOperator = GreaterThan
	case OperatorGreaterThanEq:
		newOperator = GreaterThanEq
	default:
		return operator
	}
	return newOperator
}
