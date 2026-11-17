package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	registeredUserEmailFormat = "autotest_%s@example.com"
	registeredUserFirstName   = "Autotest"
	// registeredUserPasswordFormat префикс гарантирует буквы обоих регистров и цифру — этого требует API
	registeredUserPasswordFormat = "Autotest1_%s"
	registeredUserRandomPart     = 12
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// RegisterNewUser регистрирует нового пользователя со случайными email и паролем и возвращает его токен и данные, ожидает 201
func RegisterNewUser(t *testing.T) models.LoginOkResponse {
	request := models.RegisterRequest{
		Email:     fmt.Sprintf(registeredUserEmailFormat, strings.ToLower(utils.RandomString(registeredUserRandomPart))),
		Password:  fmt.Sprintf(registeredUserPasswordFormat, utils.RandomString(registeredUserRandomPart)),
		FirstName: registeredUserFirstName,
	}

	return response.DecodeJSON[models.LoginOkResponse](t, auth.HttpPostAuthRegister(t, request), http.StatusCreated)
}
