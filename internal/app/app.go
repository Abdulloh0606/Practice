package app

import (
	"fmt"
	"miniTrello/pkg/database/postgresql"
)

func Run() error{
	_, err := postgresql.NewPostgresPgxCon()
	if err != nil {
		return fmt.Errorf("db connection error: %v", err)
	}
	fmt.Println("CONNECTION SUCCES!")
	return nil
}
