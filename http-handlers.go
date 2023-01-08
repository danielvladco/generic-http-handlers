package generic_http_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TODO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

type Handler[Req any, Req2 interface{ *Req }, Res any] func(ctx context.Context, req Req2) (Res, error)

type ErrorHandler[E Encoder] func(w http.ResponseWriter, err error, status int, encoderFactory func(w io.Writer) E)

func GET[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res]) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, http.MethodGet)
}

func POST[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res]) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, http.MethodPost)
}

func PUT[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res]) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, http.MethodPut)
}

func PATCH[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res]) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, http.MethodPatch)
}

func DELETE[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res]) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, http.MethodDelete)
}

func HandleJson[
	Req any,
	Req2 interface{ *Req },
	Res any,
](handler Handler[Req, Req2, Res],
	methods ...string) http.Handler {
	return Handle(json.NewEncoder, json.NewDecoder, BasicJsonErrorHandler, handler, methods...)
}

func Handle[
	E Encoder,
	D Decoder,
	Req any,
	Req2 interface{ *Req },
	Res any,
](
	newEncoder func(w io.Writer) E,
	newDecoder func(r io.Reader) D,
	handleError ErrorHandler[E],
	handler Handler[Req, Req2, Res],
	methods ...string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				break
			}

			handleError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed, newEncoder)
			return
		}
		var req = Req2(new(Req))

		if err := newDecoder(r.Body).Decode(req); err != nil {
			handleError(w, err, http.StatusInternalServerError, newEncoder)
			return
		}
		res, err := handler(r.Context(), req)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, newEncoder)
			return
		}
		if err := newEncoder(w).Encode(res); err != nil {
			handleError(w, err, http.StatusInternalServerError, newEncoder)
			return
		}
	})
}

func BasicJsonErrorHandler(w http.ResponseWriter, err error, status int, newEncoder func(w io.Writer) *json.Encoder) {
	_ = newEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
	w.WriteHeader(status)
}
