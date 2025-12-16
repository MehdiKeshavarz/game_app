package main

import (
	"fmt"
	"game_app/entity"
	"game_app/repository/mysql"
)

func main() {

}

func testUserMySqlRepo() {

	mysqlRepo := mysql.New()
	user, err := mysqlRepo.Register(entity.User{
		ID:          0,
		Name:        "fardinKeshavarz",
		PhoneNumber: "09373068746",
	})
	if err != nil {
		return
	}

	fmt.Println(user)

	isUnique, uErr := mysqlRepo.IsPhoneNumberUnique(user.PhoneNumber)

	if uErr != nil {
		fmt.Println(uErr)
	} else {
		fmt.Println("phone_number is unique:", isUnique)
	}

}
