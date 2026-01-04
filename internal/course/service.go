package course

import (
	"context"
	"log"
	"time"

	"github.com/juanjoaquin/back-g-domain/domain"
)

type Service interface {
	Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) /* Pasamos el Filtrado de params */ ([]domain.Course, error) // Get All
	Get(ctx context.Context, id string) (*domain.Course, error)
	Count(ctx context.Context, filters Filters) (int, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, name, startDate, endDate *string) error // Los dates fueron parseados a string
}

type Filters struct {
	Name string
}

type service struct {
	log *log.Logger
	// Ahora debemos pasar el Repository
	repo Repository
}

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {

	// Parseamos para los valores de Fecha
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, ErrInvalidStartDate // ErrInvalidStartDate
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, ErrInvalidStartDate // ErrInvalidStartDate
	}

	// Validacion que el Start Date no puede ser superior a End Date. Esto se usa con After de GO, relacionado a Times propios del lenguaje
	if startDateParsed.After(endDateParsed) {
		s.log.Println(ErrStartDateAfterEndDate)
		return nil, ErrStartDateAfterEndDate
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)

	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {

	course, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, name, startDate, endDate *string) error {
	var startDateParsed, endDateParsed *time.Time

	// Para la validacion del Start Date > End Date. Debemos traernos el curso
	course, err := s.repo.Get(ctx, id)

	if err != nil {
		return err
	}

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println(err)
			return ErrInvalidStartDate
		}

		// Una vez traido el Course. Hacemos la validacion
		if date.After(course.EndDate) {
			s.log.Println(ErrStartDateAfterEndDate)
			return ErrStartDateAfterEndDate
		}

		startDateParsed = &date
	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}

		// Aca tambien hacemos la validacion
		if course.StartDate.After(date) {
			s.log.Println(ErrStartDateAfterEndDate)
			return ErrStartDateAfterEndDate
		}

		endDateParsed = &date
	}

	return s.repo.Update(ctx, id, name, startDateParsed, endDateParsed)
}
