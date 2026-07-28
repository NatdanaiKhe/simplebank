// simplebank/service/container.go
package service

import (
	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
)

type ServiceContainer struct {
	AccountService  AccountService
	TransferService TransferService
	UserService     UserService
}

func NewServiceContainer(store db.Store) *ServiceContainer {
	return &ServiceContainer{
		AccountService:  NewAccountService(store),
		TransferService: NewTransferService(store),
		UserService:     NewUserService(store),
	}
}
