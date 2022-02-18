package middlewares

import (
	"errors"
	"github.com/huyct/CRUD-go/auth"
	"github.com/huyct/CRUD-go/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CheckJwt(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := jwt.Verify(r)

		if err != nil {
			res.ERROR(w, 401, errors.New("Unthorized"))
			return
		}

		next(w, r, ps)
	}
}
