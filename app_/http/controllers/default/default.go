package DefaultController

import (
	"net/http"

	//Framework
	e "../../errors"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	e.ReturnSuccess(w, http.StatusOK, "hello world")
}

func SayError(w http.ResponseWriter, r *http.Request) {
	e.ReturnError(w, http.StatusBadRequest, e.GenericError, "something went wrong")
}
