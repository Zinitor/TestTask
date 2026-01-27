package messages

import (
	"ReviewTask/pkg/logger"
	"ReviewTask/pkg/postgres"
)

type MessagesRepository struct {
	logger logger.Interface
	db     *postgres.DB
}

func NewRepository(logger logger.Interface, db *postgres.DB) MessagesRepository {
	return MessagesRepository{
		logger: logger,
		db:     db,
	}
}
func (r MessagesRepository) Add() {

}

func (r MessagesRepository) Get() {
}

func (r MessagesRepository) Delete() {
}
