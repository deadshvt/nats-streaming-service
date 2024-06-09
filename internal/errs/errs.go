package errs

import (
	"errors"
	"fmt"
)

var (
	ErrOrderExists = errors.New("order already exists")

	ErrInvalidOrderID = errors.New("invalid order id")

	ErrUnsupportedDBType    = errors.New("unsupported database type")
	ErrUnsupportedCacheType = errors.New("unsupported cache type")

	ErrJSONUnmarshal = errors.New("failed to unmarshal json")
	ErrValidateOrder = errors.New("failed to validate order")
	ErrCreateOrder   = errors.New("failed to create order")
	ErrGetOrderByID  = errors.New("failed to get order by id")
	ErrParseTemplate = errors.New("failed to parse template")

	ErrOrderNotFound = errors.New("order not found")
)

func WrapError(wrap, err error) error {
	if err == nil {
		return wrap
	}

	return fmt.Errorf("%s: %w", wrap.Error(), err)
}
