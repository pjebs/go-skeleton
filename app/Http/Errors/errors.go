package Error

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jeffail/gabs"
	"github.com/pjebs/jsonerror"
	"gopkg.in/unrolled/render.v1"
)

type ErrorCode uint8

// Add more specific error codes here
const (
	RequestEmptyOrInvalidCode ErrorCode = 1
	GenericErrorCode                    = 2
	SpecificErrorCode                   = 3
)

// Add more specific error descriptions here
var (
	RequestEmptyOrInvalid = errors.New("request is empty or invalid")
	GenericError          = errors.New("generic error")
	SpecificError         = errors.New("specific error")
)

var errorMap = map[error]ErrorCode{}

func init() {
	errorMap[RequestEmptyOrInvalid] = RequestEmptyOrInvalidCode
	errorMap[GenericError] = GenericErrorCode
	errorMap[SpecificError] = SpecificErrorCode
}

func New(err error, message ...string) error {
	code := errorMap[err]
	var je jsonerror.JE

	if len(message) == 0 {
		je = jsonerror.New(int(code), err.Error(), "")
	} else {
		je = jsonerror.New(int(code), err.Error(), message[0])
	}

	return je
}

// Returns API Error response in a standard format
func ReturnError(w http.ResponseWriter, httpStatusCode int, err error, message ...string) {

	code := errorMap[err]
	var je jsonerror.JE

	if len(message) == 0 {
		je = jsonerror.New(int(code), err.Error(), "")
	} else {
		je = jsonerror.New(int(code), err.Error(), message[0])
	}

	r := render.New(render.Options{})
	r.JSON(w, httpStatusCode, je.Render())
	return
}

// Returns API response in json format
func ReturnSuccess(w http.ResponseWriter, httpStatusCode int, data interface{}) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(httpStatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}

// ParseBody validates and converts the JSON Post body into gabs.Container structs
func ParseBody(r *http.Request) (*gabs.Container, error) {

	if r.Body == nil {
		return nil, RequestEmptyOrInvalid
	}

	reqData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	reqJSON, err := gabs.ParseJSON(reqData)
	if err != nil {
		return nil, err
	}
	return reqJSON, nil
}
