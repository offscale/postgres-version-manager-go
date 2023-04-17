package pvm

import (
	"errors"
	"os"
)

func ensureDirsExist(dirs ...string) error {
	for _, d := range dirs {
		if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(d, os.ModeSticky|os.ModePerm|os.ModeDir); err != nil {
				return err
			}
		}
	}
	return nil
}
