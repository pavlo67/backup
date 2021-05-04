package server_http_v2

import (
	"fmt"
	"net/http"

	"github.com/pavlo67/common/common/errors"
)

func CommonFragments(title, htmlHeader, htmlMessage, htmlError, htmlIndex, htmlContent string) map[string]string {
	return map[string]string{
		"title":   title,
		"header":  htmlHeader,
		"message": htmlMessage,
		"error":   htmlError,
		"index":   htmlIndex,
		"content": htmlContent,
	}
}

func ErrorPage(httpStatus int, err error, publicDetails string, req *http.Request) (ResponsePage, error) {
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}

	fragments := CommonFragments(
		"помилка",
		"",
		"",
		"На жаль, виникла помилка:-(\n<p>"+publicDetails,
		"",
		"",
	)

	if req != nil {
		err = errors.CommonError(fmt.Errorf(" on [%s %s]", req.Method, req.URL), publicDetails, err)
	} else {
		err = errors.CommonError(publicDetails, err)
	}

	return ResponsePage{
		Status:    httpStatus,
		Fragments: fragments,
	}, err
}
