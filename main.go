package main

import (
	"encoding/json"
	"fmt"
	"game_app/entity"
	"game_app/repository/mysql"
	"game_app/service/userservice"
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/health-check", healthCheckHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		_, err := fmt.Fprintf(writer, `{"error": "invalid http method"'}`)
		if err != nil {
			return
		}

	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
	}

	var uReq userservice.RegisterRequest
	uErr := json.Unmarshal(data, &uReq)
	if uErr != nil {
		_, wErr := writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, uErr.Error())))
		if wErr != nil {
			return
		}
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(mysqlRepo)

	_, rErr := userSvc.Register(uReq)

	if rErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, rErr.Error())))

		return
	}

	writer.Write([]byte(`{"message": "user registered"}`))
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(writer, `{"message": "everything is good!"}`)
	if err != nil {
		return
	}
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
