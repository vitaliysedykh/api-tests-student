package advertisement

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisement/models"
	commonModels "api-tests-template/internal/managers/common/models"
	"api-tests-template/internal/utils"
)

const (
	advertisementTitlePrefix     = "Autotest "
	advertisementTitleRandomPart = 16
	advertisementDescription     = "Объявление создано API-автотестом"
	advertisementPrice           = 1999.99
	advertisementQuantity        = 5
)

// NewAdvertisementRequest возвращает валидный запрос создания объявления с уникальным названием и photosCount фотографиями
func NewAdvertisementRequest(t *testing.T, photosCount int) models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       advertisementTitlePrefix + utils.RandomString(advertisementTitleRandomPart),
		Description: advertisementDescription,
		Price:       utils.Ptr(advertisementPrice),
		Quantity:    utils.Ptr(advertisementQuantity),
		PhotoPaths:  utils.CreatePhotoFiles(t, photosCount),
	}
}

// CreateAdvertisement создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.AdvertisementResponse {
	return response.DecodeJSON[models.AdvertisementResponse](t,
		advertisement.HttpPostAdvertisement(t, token, request), http.StatusCreated)
}

// CreateAdvertisementWithCleanup создаёт объявление и удаляет его по завершении теста, ожидает 201
func CreateAdvertisementWithCleanup(t *testing.T, token string, request models.CreateAdvertisementRequest) models.AdvertisementResponse {
	created := CreateAdvertisement(t, token, request)
	t.Cleanup(func() {
		utils.LogWithLabelAndTimeStamp("Cleanup", "Удаляем созданное в тесте объявление "+created.Id)
		DeleteAdvertisement(t, token, created.Id)
	})
	return created
}

// CreateAdvertisementExpectError пытается создать объявление и возвращает ошибку API, ожидает expectedStatusCode
func CreateAdvertisementExpectError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) commonModels.ErrorResponse {
	return response.DecodeJSON[commonModels.ErrorResponse](t,
		advertisement.HttpPostAdvertisement(t, token, request), expectedStatusCode)
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.AdvertisementResponse {
	return response.DecodeJSON[models.AdvertisementResponse](t,
		advertisement.HttpGetAdvertisement(t, id), http.StatusOK)
}

// GetAdvertisements возвращает страницу каталога объявлений по параметрам запроса, ожидает 200
func GetAdvertisements(t *testing.T, params map[string]string) models.AdvertisementListResponse {
	return response.DecodeJSON[models.AdvertisementListResponse](t,
		advertisement.HttpGetAdvertisements(t, params), http.StatusOK)
}

// SearchAdvertisements ищет объявления каталога по вхождению строки в название, ожидает 200
func SearchAdvertisements(t *testing.T, search string) models.AdvertisementListResponse {
	return GetAdvertisements(t, map[string]string{"search": search})
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.PhotoResponse {
	return response.DecodeJSON[[]models.PhotoResponse](t,
		advertisement.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}

// DownloadPhoto скачивает файл фотографии по url и возвращает его содержимое, ожидает 200
func DownloadPhoto(t *testing.T, photoUrl string) []byte {
	return response.RequireStatusCode(t, advertisement.HttpGetPhotoFile(t, photoUrl), http.StatusOK)
}

// DeleteAdvertisement удаляет объявление, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	response.RequireStatusCode(t, advertisement.HttpDeleteAdvertisement(t, token, id), http.StatusNoContent)
}
