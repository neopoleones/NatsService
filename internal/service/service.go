package service

import "context"

type Service interface {
	Subscribe(context.Context) error
	Run() error
	Utilize()
}
