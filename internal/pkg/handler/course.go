package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/juanjoaquin/back-g-course/internal/course"
	"github.com/juanjoaquin/back-g-response/response"
)

// Definimos la funcion. Recibira el Context y los Endpoints definidos.
func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {

	router := mux.NewRouter()

	// Manejo de Errores con Go Kit
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	//No usamos router.HandleFunc() como estabamos usando. Usaremos Handle de Gorilla Mux
	// Tampoco nos traeremos el userEndpoint. Usaremos el httptransport.NewServer() de Go Kit
	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create), // Debemos hacer una conversion. Llamamos al Endpoint de GO KIT, y lo encapsulamos dentro del endpoints.Create del Controller
		decodeCreateCourse,
		encodeResponse,
		opts..., // Tambien le pasamos el OPTS del Middleware para descrifar los errores
	)).Methods("POST")
	/* EXPLICACION DE LOS PARAMETROS Y LA FUNCIONES:
	El handle primero va a enviar al Decode el POST para crear el usuario. Este Decode ejecuta la funcion, y hace la conversion.
	En caso de que no puede generara un error. Si esta OK, enviara la Request 200.

	Despues pasa el encodeResponse donde recibe la respuesta, y accede a la creacion y al status 200.
	*/

	router.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourses,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	router.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

// Esta funcion se encarga de hacer un Decode dentro del request cuando nosotros hagamos el store de un User
func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	// Definimos la Request del CreateReq
	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error())) // Le pasamos el package del Response
	}

	return req, nil
}

// Hacemos un Enconde del Response.
// Esto lo que va a devolver despues el Endpoint una vez que retorne
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	// Hacemos un reconverse de nuestro Package de Response
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())       // Esto tambien
	return json.NewEncoder(w).Encode(r) // Retornamos el response
}

// Aqui pasara por otra instancia donde decodifica el Error. En caso de haber un error por ejemplo un 400. Lo descifra.
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response) // Hacemos una conversion del Error al Response
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp) // No debemos hacer un return. Solo mapearle al response que recibimos por parametro, lo que queremos retornar al cliente

}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format '%v'", err.Error()))
	}
	path := mux.Vars(r)
	req.ID = path["id"]
	return req, nil

}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {
	path := mux.Vars(r)
	req := course.DeleteReq{
		ID: path["id"],
	}

	return req, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := course.GetReq{
		ID: p["id"],
	}
	return req, nil
}

func decodeGetAllCourses(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil

}
