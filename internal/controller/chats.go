package controller

import (
	"ReviewTask/pkg/logger"
	"fmt"
	"net/http"
)

type ChatsService interface {
	create()
	get_messages()
	get()
}

type ChatsHander struct {
	logger  logger.Interface
	service ChatsService
}

func NewChatsController(logger logger.Interface, service ChatsService) ChatsHander {
	return ChatsHander{
		logger:  logger,
		service: service,
	}
}

func (h ChatsHander) hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func (h ChatsHander) headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
