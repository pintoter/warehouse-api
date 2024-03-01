package model

import "errors"

var (
	ErrInvalidInput               = errors.New("invalid input params")
	ErrInvalidCode                = errors.New("invalid code of product")
	ErrInvalidQuantity            = errors.New("required quantity of products is missing")
	ErrInvalidReservationQuantity = errors.New("too many products for release")
)
