package course

import (
	"errors"
	"fmt"
)

// Este archivo son los errores customizados para los campos. Por ejemplo: First Name, Last Name del User, etc...

var ErrInvalidStartDate = errors.New("Invalid Start Date")
var ErrInvalidEndDate = errors.New("Invalid End Date")
var ErrNameRequired = errors.New("Name is Required")
var ErrStartRequired = errors.New("Start date is Required")
var ErrEndRequired = errors.New("End date is Required")

// Manejo de Errores con Parametros Dinamicos
type ErrNotFound struct {
	CourseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesnt exists", e.CourseID)
}
