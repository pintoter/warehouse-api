package model

import "errors"

var (
	ErrInvalidInput               = errors.New("invalid input params")
	ErrInvalidCode                = errors.New("invalid code of product")
	ErrInvalidQuantity            = errors.New("required quantity of products is missing")
	ErrInvalidReservationQuantity = errors.New("too many products for release")
	ErrInternalServer             = errors.New("internal server error, try later")
	ErrFailedReservation          = errors.New("failed to reserve item from warehouse")
)
