package favorites

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisement"
	advertisementModels "api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// nonExistentAdvertisementId валидный UUID, под которым нет объявления.
// Нулевой UUID не подходит: API считает его пустым полем и отвечает 400 вместо 404
const nonExistentAdvertisementId = "5f0c6d2e-9a41-4b8e-8f3a-2c7d1e6b9a10"

type TestSuite struct {
	suite.Suite
	// seller пользователь, зарегистрированный на время прогона, — владелец объявлений: в избранное можно
	// добавить только чужое объявление, а под TEST_LOGIN объявления не создаются, чтобы не менять его данные
	seller authModels.LoginOkResponse
	// buyer пользователь из переменных окружения, работающий с избранным
	buyer authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация покупателя с кредами из переменных окружения")
	s.buyer = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Регистрация нового продавца на время прогона")
	s.seller = auth.RegisterNewUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// createSellerAdvertisement создаёт объявление продавца, которое удаляется по завершении теста
func (s *TestSuite) createSellerAdvertisement() advertisementModels.AdvertisementResponse {
	base.Precondition("Продавец создаёт объявление")
	return advertisement.CreateAdvertisementWithCleanup(s.T(), s.seller.Token,
		advertisement.NewAdvertisementRequest(s.T(), 1))
}

// requireInFavorites проверяет, что объявление есть в избранном пользователя
func (s *TestSuite) requireInFavorites(token string, ad advertisementModels.AdvertisementResponse) {
	item, ok := favorites.GetFavorites(s.T(), token).FindByAdvertisementId(ad.Id)
	require.Truef(s.T(), ok, "Объявление %s отсутствует в избранном", ad.Id)
	require.Equal(s.T(), ad.Title, item.Advertisement.Title, "Название объявления в избранном не совпадает")
	base.Result("Объявление " + ad.Id + " есть в избранном")
}

// requireNotInFavorites проверяет, что объявления нет в избранном пользователя
func (s *TestSuite) requireNotInFavorites(token, advertisementId string) {
	_, ok := favorites.GetFavorites(s.T(), token).FindByAdvertisementId(advertisementId)
	require.Falsef(s.T(), ok, "Объявление %s всё ещё в избранном", advertisementId)
	base.Result("Объявления " + advertisementId + " нет в избранном")
}
