package repository

import (
	"context"
	"fmt"
	"housekeeper/internal/booking-service/model"

	"github.com/kamva/mgm/v3"
)

type repository struct {
}

func NewRepository() *repository {
	// TODO: create index if need

	return &repository{}
}

func (*repository) CreateNewJob(ctx context.Context, job *model.Job) error {
	coll := mgm.Coll(job)

	err := coll.CreateWithCtx(ctx, job)
	if err != nil {
		return fmt.Errorf("error when write to mongo: %w", err)
	}

	return nil
}
