package createAdvertisement

import (
	"fmt"
	"net/http"

	"api-tests-template/internal/constants/apiErrors"
	"api-tests-template/internal/constants/limits"
	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

// Negative: Проверяем 400 ошибку при создании объявления без обязательных параметров или с некорректными параметрами
func (s *TestSuite) TestCreateAdvertisementWithInvalidParams() {
	base.Precondition(fmt.Sprintf("Готовим %d фотографии для случая с превышением лимита", limits.MaxPhotos+1))
	tooManyPhotos := utils.CreatePhotoFiles(s.T(), limits.MaxPhotos+1)

	testCases := []struct {
		name            string
		modify          func(request *models.CreateAdvertisementRequest)
		expectedMessage string
	}{
		{
			name:            "Пустое название",
			modify:          func(request *models.CreateAdvertisementRequest) { request.Title = "" },
			expectedMessage: apiErrors.MessageTitleRequired,
		},
		{
			name:            "Пустое описание",
			modify:          func(request *models.CreateAdvertisementRequest) { request.Description = "" },
			expectedMessage: apiErrors.MessageDescriptionRequired,
		},
		{
			name:            "Без фотографий",
			modify:          func(request *models.CreateAdvertisementRequest) { request.PhotoPaths = nil },
			expectedMessage: apiErrors.MessagePhotoRequired,
		},
		{
			name: "Слишком длинное название",
			modify: func(request *models.CreateAdvertisementRequest) {
				request.Title = utils.RandomString(limits.MaxTitleLength + 1)
			},
			expectedMessage: fmt.Sprintf(apiErrors.MessageTitleTooLongFormat, limits.MaxTitleLength, limits.MaxTitleLength+1),
		},
		{
			name: "Слишком длинное описание",
			modify: func(request *models.CreateAdvertisementRequest) {
				request.Description = utils.RandomString(limits.MaxDescriptionLength + 1)
			},
			expectedMessage: fmt.Sprintf(apiErrors.MessageDescriptionLengthFormat, limits.MaxDescriptionLength),
		},
		{
			name:            "Отрицательная цена",
			modify:          func(request *models.CreateAdvertisementRequest) { request.Price = utils.Ptr(-1.0) },
			expectedMessage: fmt.Sprintf(apiErrors.MessagePriceOutOfRangeFormat, limits.MaxPrice),
		},
		{
			name: "Цена выше максимальной",
			modify: func(request *models.CreateAdvertisementRequest) {
				request.Price = utils.Ptr(float64(limits.MaxPrice + 1))
			},
			expectedMessage: fmt.Sprintf(apiErrors.MessagePriceOutOfRangeFormat, limits.MaxPrice),
		},
		{
			name: "Количество меньше минимального",
			modify: func(request *models.CreateAdvertisementRequest) {
				request.Quantity = utils.Ptr(limits.MinQuantity - 1)
			},
			expectedMessage: fmt.Sprintf(apiErrors.MessageQuantityOutOfRangeFormat, limits.MinQuantity, limits.MaxQuantity),
		},
		{
			name: "Количество больше максимального",
			modify: func(request *models.CreateAdvertisementRequest) {
				request.Quantity = utils.Ptr(limits.MaxQuantity + 1)
			},
			expectedMessage: fmt.Sprintf(apiErrors.MessageQuantityOutOfRangeFormat, limits.MinQuantity, limits.MaxQuantity),
		},
		{
			name:            "Фотографий больше максимума",
			modify:          func(request *models.CreateAdvertisementRequest) { request.PhotoPaths = tooManyPhotos },
			expectedMessage: fmt.Sprintf(apiErrors.MessageTooManyPhotosFormat, limits.MaxPhotos),
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			base.Precondition("Готовим валидный запрос и портим в нём один параметр: " + testCase.name)
			request := advertisement.NewAdvertisementRequest(s.T(), 1)
			testCase.modify(&request)

			base.Step("Отправляем запрос на создание объявления и проверяем ошибку", func() {
				apiError := advertisement.CreateAdvertisementExpectError(s.T(), s.author.Token, request, http.StatusBadRequest)
				base.RequireApiError(s.T(), apiError, apiErrors.ErrorInvalidRequest, testCase.expectedMessage)
			})
		})
	}
}

// Negative: Проверяем 401 ошибку при создании объявления без токена и с неверным токеном
func (s *TestSuite) TestCreateAdvertisementWithoutAuthorization() {
	for _, testCase := range base.UnauthorizedCases() {
		s.Run(testCase.Name, func() {
			base.Precondition("Готовим валидный запрос на создание объявления")
			request := advertisement.NewAdvertisementRequest(s.T(), 1)

			base.Step("Отправляем запрос на создание объявления и проверяем ошибку авторизации", func() {
				apiError := advertisement.CreateAdvertisementExpectError(s.T(), testCase.Token, request, http.StatusUnauthorized)
				base.RequireApiError(s.T(), apiError, apiErrors.ErrorUnauthorized, testCase.ExpectedMessage)
			})
		})
	}
}
