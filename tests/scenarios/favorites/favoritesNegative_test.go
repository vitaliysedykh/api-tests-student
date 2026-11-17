package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/apiErrors"
	"api-tests-template/internal/constants/limits"
	"api-tests-template/internal/managers/auth"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// invalidAdvertisementId строка, не являющаяся UUID
const invalidAdvertisementId = "not-a-uuid"

// Negative: Проверяем 404 ошибку при добавлении в избранное несуществующего объявления
func (s *TestSuite) TestAddNonExistentAdvertisementToFavorites() {
	base.Step("Добавляем в избранное объявление с несуществующим id", func() {
		apiError := favorites.AddFavoriteExpectError(s.T(), s.buyer.Token, nonExistentAdvertisementId, http.StatusNotFound)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorNotFound, apiErrors.MessageAdvertisementNotFound)
	})
}

// Negative: Проверяем 404 ошибку при удалении из избранного несуществующего объявления
func (s *TestSuite) TestRemoveNonExistentAdvertisementFromFavorites() {
	base.Step("Удаляем из избранного объявление с несуществующим id", func() {
		apiError := favorites.RemoveFavoriteExpectError(s.T(), s.buyer.Token, nonExistentAdvertisementId, http.StatusNotFound)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorNotFound, apiErrors.MessageFavoriteNotFound)
	})
}

// Negative: Проверяем 401 ошибку на всех ручках избранного без токена и с неверным токеном
func (s *TestSuite) TestFavoritesWithoutAuthorization() {
	endpoints := []struct {
		name    string
		request func(t *testing.T, token string) commonModels.ErrorResponse
	}{
		{
			name: "GET /favorites",
			request: func(t *testing.T, token string) commonModels.ErrorResponse {
				return favorites.GetFavoritesExpectError(t, token, http.StatusUnauthorized)
			},
		},
		{
			name: "POST /favorites",
			request: func(t *testing.T, token string) commonModels.ErrorResponse {
				return favorites.AddFavoriteExpectError(t, token, nonExistentAdvertisementId, http.StatusUnauthorized)
			},
		},
		{
			name: "DELETE /favorites/{ad_id}",
			request: func(t *testing.T, token string) commonModels.ErrorResponse {
				return favorites.RemoveFavoriteExpectError(t, token, nonExistentAdvertisementId, http.StatusUnauthorized)
			},
		},
	}

	for _, endpoint := range endpoints {
		for _, authCase := range base.UnauthorizedCases() {
			s.Run(endpoint.name+": "+authCase.Name, func() {
				base.Step(fmt.Sprintf("Вызываем %s: %s", endpoint.name, authCase.Name), func() {
					apiError := endpoint.request(s.T(), authCase.Token)
					base.RequireApiError(s.T(), apiError, apiErrors.ErrorUnauthorized, authCase.ExpectedMessage)
				})
			})
		}
	}
}

// Negative: Проверяем 409 ошибку при повторном добавлении объявления в избранное
func (s *TestSuite) TestAddFavoriteTwice() {
	ad := s.createSellerAdvertisement()

	base.Step("Покупатель добавляет объявление в избранное первый раз", func() {
		favorites.AddFavorite(s.T(), s.buyer.Token, ad.Id)
	})

	base.Step("Покупатель добавляет то же объявление повторно", func() {
		apiError := favorites.AddFavoriteExpectError(s.T(), s.buyer.Token, ad.Id, http.StatusConflict)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorAlreadyExists, apiErrors.MessageFavoriteAlreadyExists)
	})
}

// Negative: Проверяем 403 ошибку при добавлении в избранное собственного объявления
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	ad := s.createSellerAdvertisement()

	base.Step("Продавец добавляет в избранное своё объявление", func() {
		apiError := favorites.AddFavoriteExpectError(s.T(), s.seller.Token, ad.Id, http.StatusForbidden)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorForbidden, apiErrors.MessageCannotFavoriteOwn)
	})
}

// Negative: Проверяем 400 ошибку при добавлении в избранное по id, который не является UUID
func (s *TestSuite) TestAddFavoriteWithInvalidAdvertisementId() {
	base.Step("Добавляем в избранное объявление с id не в формате UUID", func() {
		apiError := favorites.AddFavoriteExpectError(s.T(), s.buyer.Token, invalidAdvertisementId, http.StatusBadRequest)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorInvalidRequest,
			fmt.Sprintf(apiErrors.MessageInvalidUuidLengthFormat, len(invalidAdvertisementId)))
	})
}

// Negative: Проверяем 400 ошибку при добавлении в избранное сверх лимита в 50 объявлений
func (s *TestSuite) TestAddFavoriteOverLimit() {
	base.Precondition("Регистрируем отдельного покупателя: избранное, заполненное до лимита, не должно остаться у пользователя из TEST_LOGIN")
	buyer := auth.RegisterNewUser(s.T())
	ad := s.createSellerAdvertisement()

	base.Precondition(fmt.Sprintf("Заполняем избранное покупателя до лимита: %d объявлений каталога", limits.FavoritesLimit))
	favorites.FillFavoritesToLimit(s.T(), buyer.Token, limits.FavoritesLimit)

	base.Step("Проверяем, что избранное заполнено до лимита", func() {
		total := favorites.GetFavorites(s.T(), buyer.Token).Total
		require.Equalf(s.T(), limits.FavoritesLimit, total,
			"В избранном %d объявлений вместо %d: часть объявлений удалили владельцы во время заполнения", total, limits.FavoritesLimit)
		base.Result(fmt.Sprintf("В избранном %d объявлений", total))
	})

	base.Step("Добавляем в избранное объявление сверх лимита", func() {
		apiError := favorites.AddFavoriteExpectError(s.T(), buyer.Token, ad.Id, http.StatusBadRequest)
		base.RequireApiError(s.T(), apiError, apiErrors.ErrorLimitReached,
			fmt.Sprintf(apiErrors.MessageFavoritesLimitFormat, limits.FavoritesLimit))
	})
}
