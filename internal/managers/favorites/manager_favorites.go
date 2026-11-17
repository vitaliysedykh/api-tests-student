package favorites

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/constants/limits"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisement"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

// AddFavorite добавляет объявление в избранное и возвращает запись избранного, ожидает 201
func AddFavorite(t *testing.T, token, advertisementId string) models.FavoriteItem {
	return response.DecodeJSON[models.FavoriteItem](t,
		favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId}), http.StatusCreated)
}

// AddFavoriteExpectError пытается добавить объявление в избранное и возвращает ошибку API, ожидает expectedStatusCode
func AddFavoriteExpectError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	return response.DecodeJSON[commonModels.ErrorResponse](t,
		favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId}), expectedStatusCode)
}

// GetFavorites возвращает первую страницу избранного максимального размера, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoriteListResponse {
	return response.DecodeJSON[models.FavoriteListResponse](t,
		favorites.HttpGetFavorites(t, token, limits.FavoritesMaxPageSize), http.StatusOK)
}

// GetFavoritesExpectError пытается получить избранное и возвращает ошибку API, ожидает expectedStatusCode
func GetFavoritesExpectError(t *testing.T, token string, expectedStatusCode int) commonModels.ErrorResponse {
	return response.DecodeJSON[commonModels.ErrorResponse](t,
		favorites.HttpGetFavorites(t, token, limits.FavoritesMaxPageSize), expectedStatusCode)
}

// RemoveFavorite удаляет объявление из избранного, ожидает 204
func RemoveFavorite(t *testing.T, token, advertisementId string) {
	response.RequireStatusCode(t, favorites.HttpDeleteFavorite(t, token, advertisementId), http.StatusNoContent)
}

// RemoveFavoriteExpectError пытается удалить объявление из избранного и возвращает ошибку API, ожидает expectedStatusCode
func RemoveFavoriteExpectError(t *testing.T, token, advertisementId string, expectedStatusCode int) commonModels.ErrorResponse {
	return response.DecodeJSON[commonModels.ErrorResponse](t,
		favorites.HttpDeleteFavorite(t, token, advertisementId), expectedStatusCode)
}

// FillFavoritesToLimit добавляет в избранное объявления каталога, начиная с самых старых, пока их не станет limit.
// Каталог общий для всех: объявление могут удалить между чтением страницы и добавлением, такие (404) пропускаются
func FillFavoritesToLimit(t *testing.T, token string, limit int) {
	added := 0
	cursor := ""
	for {
		page := advertisement.GetAdvertisements(t, oldestFirstCatalogParams(cursor))
		for _, ad := range page.Items {
			if added == limit {
				return
			}
			if addFavoriteSkippingDeleted(t, token, ad.Id) {
				added++
			}
		}
		if added == limit {
			return
		}
		require.Truef(t, page.HasNext, "В каталоге не хватило объявлений: в избранное добавлено %d из %d", added, limit)
		cursor = page.NextCursor
	}
}

// oldestFirstCatalogParams параметры страницы каталога от старых объявлений к новым: старые реже удаляют
func oldestFirstCatalogParams(cursor string) map[string]string {
	params := map[string]string{
		"sort_by":    "created_at",
		"sort_order": "asc",
		"limit":      strconv.Itoa(limits.CatalogMaxPageSize),
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	return params
}

// addFavoriteSkippingDeleted добавляет объявление в избранное; возвращает false, если объявление уже удалено (404)
func addFavoriteSkippingDeleted(t *testing.T, token, advertisementId string) bool {
	rawResponse := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId})
	body := response.ReadBody(t, rawResponse)

	if rawResponse.StatusCode == http.StatusNotFound {
		utils.LogWithLabelAndTimeStamp("Precondition",
			fmt.Sprintf("Объявление %s удалено владельцем до добавления, пропускаем", advertisementId))
		return false
	}

	require.Equalf(t, http.StatusCreated, rawResponse.StatusCode,
		"Не удалось добавить объявление %s в избранное. Тело ответа: %s", advertisementId, body)
	return true
}
