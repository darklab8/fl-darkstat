package main

import (
	"errors"
	"fmt"
	"testing"
)

func DownloadFileTest(give_error bool, status_code int) error {
	if give_error {
		return ErrDownloadFile{
			status_code:  status_code,
			original_err: errors.New("smth"),
		}
	}

	return nil
}

func TestError(t *testing.T) {

	err := DownloadFileTest(true, 403)
	var mErr ErrDownloadFile
	if errors.As(err, &mErr) {
		fmt.Println("1 error is present", mErr)
	} else {
		fmt.Println("1 error is NOT present", mErr)
	}

	err = DownloadFileTest(true, 403)
	var mErr2 ErrDownloadFile
	if errors.As(err, &mErr2) {
		fmt.Println("2 error is present", mErr)
	} else {
		fmt.Println("2 error is NOT present", mErr)
	}
}
