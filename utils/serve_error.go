package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-openapi/errors"
)

func errorAsJSON(err errors.Error) []byte {
	b, _ := json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{strconv.Itoa(int(err.Code())), err.Error()})
	return b
}

func asHTTPCode(input int) int {
	if input >= 600 {
		return 422
	}
	return input
}

func flattenComposite(errs *errors.CompositeError) *errors.CompositeError {
	var res []error
	for _, er := range errs.Errors {
		switch e := er.(type) {
		case *errors.CompositeError:
			if len(e.Errors) > 0 {
				flat := flattenComposite(e)
				if len(flat.Errors) > 0 {
					res = append(res, flat.Errors...)
				}
			}
		default:
			if e != nil {
				res = append(res, e)
			}
		}
	}
	return errors.CompositeValidationError(res...)
}

// ServeError the error handler interface implemenation
// it's the same as in github.com/go-openapi/errors
// except it uses code as string, so it's the same as our models.Error
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")
	switch e := err.(type) {
	case *errors.CompositeError:
		er := flattenComposite(e)
		ServeError(rw, r, er.Errors[0])
	case *errors.MethodNotAllowedError:
		rw.Header().Add("Allow", strings.Join(err.(*errors.MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	case errors.Error:
		if e == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(errorAsJSON(errors.New(http.StatusInternalServerError, "Unknown error")))
			return
		}
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(errors.New(http.StatusInternalServerError, err.Error())))
		}
	}
}
