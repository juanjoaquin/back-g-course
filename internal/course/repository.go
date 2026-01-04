package course

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/juanjoaquin/back-g-domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, course *domain.Course) error
	Get(ctx context.Context, id string) (*domain.Course, error)
	GetAll(ctx context.Context, filters Filters, offset int, limit int) ([]domain.Course, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters Filters) (int, error)
	Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

// POST: Crear Course
func (repo *repo) Create(ctx context.Context, course *domain.Course) error {
	if err := repo.db.WithContext(ctx).Create(course).Error; err != nil {
		repo.log.Printf("error: %v", err)
		return err
	}
	repo.log.Println("Course created with id:", course.ID)
	return nil
}

// GET: Get All Courses
func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var c []domain.Course

	log.Println("REPO => ejecutando query...")

	tx := repo.db.WithContext(ctx).Model(&domain.Course{})
	tx = applyFilters(tx, filters)

	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&c)

	if result.Error != nil {
		return nil, result.Error
	}

	return c, nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.Course, error) {
	course := domain.Course{ID: id}

	if err := repo.db.WithContext(ctx).First(&course).Error; err != nil {
		repo.log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}
	/* 	result := repo.db.WithContext(ctx).First(&course)

	   	if result.Error != nil {
	   		return nil, result.Error
	   	} */

	return &course, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	course := domain.Course{ID: id}

	result := repo.db.WithContext(ctx).Delete(&course)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("course %s doesnt exists", id)
		return ErrNotFound{id}
	}

	return nil
}

func (repo *repo) Update(ctx context.Context, id string, name *string, startDate *time.Time, endDate *time.Time) error {
	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := repo.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Printf("course %s doesnt exists", id)
		return ErrNotFound{id}
	}

	return nil

}

// FUNCION PARA EL APLICADO DE FILTROS
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" { // Basicamente que si viene vacio, no pasa nada, y que lo devuelva en lower o uppercase
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name) // Query de GORM para la consulta
	}

	return tx
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
