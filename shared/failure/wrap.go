package failure

import (
	"fmt"
	"strings"
)

func Wrap(originalErr error, addedStr string) error {
	return fmt.Errorf("%s:: %w", addedStr, originalErr)
}

func WrapE(originalErr error, addedErr error) error {
	return fmt.Errorf("%w:: %w", addedErr, originalErr)
}

func Split(err error) []string {
	if err == nil {
		return nil
	}
	return strings.Split(err.Error(), ":: ")
}
