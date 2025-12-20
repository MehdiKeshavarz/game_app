package main

import (
	"encoding/json"
	"fmt"
	"game_app/repository/mysql"
	"game_app/service/authservice"
	"game_app/service/userservice"
	"io"
	"net/http"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {

	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/users/login", userLoginHandler)
	http.HandleFunc("/users/profile", userProfileHandler)

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

	authSvc := authservice.New(
		JwtSignKey,
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	_, rErr := userSvc.Register(uReq)

	if rErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, rErr.Error())))

		return
	}

	writer.Write([]byte(`{"message": "user registered"}`))
}

func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
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
	var uReq userservice.LoginRequest
	uErr := json.Unmarshal(data, &uReq)

	if uErr != nil {
		_, wErr := writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, uErr.Error())))
		if wErr != nil {
			return
		}
	}

	authSvc := authservice.New(
		JwtSignKey,
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	res, rErr := userSvc.Login(uReq)

	if rErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, rErr.Error())))

		return
	}

	jsonData, jErr := json.Marshal(res)

	if jErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, jErr.Error())))
	}

	writer.Write(jsonData)
}

func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(writer, `{"message": "everything is good!"}`)
	if err != nil {
		return
	}
}

func userProfileHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		_, err := fmt.Fprintf(writer, `{"error": "invalid http method"'}`)
		if err != nil {
			return
		}
	}

	authSvc := authservice.New(
		JwtSignKey,
		"at",
		"rt",
		time.Hour*24,
		time.Hour*24*7)

	auth := req.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(auth)

	if err != nil {
		fmt.Fprintf(writer, `{"error": "token is not valid"}`)
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	res, rErr := userSvc.GetProfile(userservice.GetProfileRequest{UserID: claims.UserID})

	if rErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, rErr.Error())))

		return
	}

	jsonData, jErr := json.Marshal(res)

	if jErr != nil {
		writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, jErr.Error())))
	}

	writer.Write(jsonData)
}
