package diffsummary

import "errors"

var ErrLocaleNotFound = errors.New("locale not found")
var ErrLocaleKeyNotFound = errors.New("locale key not found")
var ErrWrongSourceScope = errors.New("wrong replacement source scope")
var ErrWrongSourceValue = errors.New("source must be a struct")
