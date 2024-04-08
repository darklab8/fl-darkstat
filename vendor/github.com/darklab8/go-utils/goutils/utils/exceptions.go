package utils

import (
	"fmt"
)

type ErrorNotFound struct {
}

func (err ErrorNotFound) Error() string {
	return fmt.Sprintf("unable to find object")
}
