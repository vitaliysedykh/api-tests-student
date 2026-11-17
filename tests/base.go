package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/apiErrors"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/utils"
)

// InvalidToken заведомо неверный авторизационный токен
const InvalidToken = "incorrect_token"

// UnauthorizedCase вариант запроса без валидной авторизации и ожидаемое сообщение об ошибке
type UnauthorizedCase struct {
	Name            string
	Token           string
	ExpectedMessage string
}

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Step выводит название шага с меткой [Step] и выполняет шаг
func Step(name string, step func()) {
	utils.LogWithLabelAndTimeStamp("Step", name)
	step()
}

// Result выводит результат успешной проверки с меткой [Result]
func Result(text string) {
	utils.LogWithLabelAndTimeStamp("Result", text)
}

// UnauthorizedCases варианты неавторизованного запроса: без токена и с неверным токеном
func UnauthorizedCases() []UnauthorizedCase {
	return []UnauthorizedCase{
		{Name: "Запрос без токена", Token: "", ExpectedMessage: apiErrors.MessageAuthorizationRequired},
		{Name: "Запрос с неверным токеном", Token: InvalidToken, ExpectedMessage: apiErrors.MessageAuthorizationInvalid},
	}
}

// RequireApiError проверяет поля error и message в ответе API с ошибкой и логирует результат
func RequireApiError(t *testing.T, actual commonModels.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, actual.Message, "Некорректное значение message параметра")
	Result(fmt.Sprintf("Получена ожидаемая ошибка: %s — %s", expectedError, expectedMessage))
}
