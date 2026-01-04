package course

import (
	"context"
	"errors"

	"github.com/juanjoaquin/back-g-meta/pkg/meta"
	"github.com/juanjoaquin/back-g-response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Delete Controller
		Update Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateReq struct {
		ID        string
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	GetAllReq struct {
		Name  string
		Limit int
		Page  int
	}

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEdnpoint(s, config),
		Get:    makeGetEndpoint(s),
		Delete: makeDeleteEndpoint(s),
		Update: makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStartRequired.Error())

		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndRequired.Error())
		}

		course, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {

			if err == ErrInvalidStartDate || err == ErrInvalidEndDate {
				return nil, response.BadRequest(err.Error())
			}
		}

		return response.Created("success", course, nil), nil
	}
}

func makeGetAllEdnpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			Name: req.Name,
		}
		count, err := s.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.InternalServerError(err.Error())

		}
		return response.OK("success", courses, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		course, err := s.Get(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())

		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStartRequired.Error())

		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndRequired.Error())

		}
		err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil

	}
}
