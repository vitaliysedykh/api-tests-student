package createAdvertisement

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	// author пользователь, зарегистрированный на время прогона. Объявления не создаются под TEST_LOGIN,
	// чтобы не менять его данные: их параллельно проверяет сценарий myAdvertisement
	author models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация нового пользователя — автора объявлений")
	s.author = auth.RegisterNewUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}
