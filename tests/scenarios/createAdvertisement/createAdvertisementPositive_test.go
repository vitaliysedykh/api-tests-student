package createAdvertisement

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/limits"
	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"

	base "api-tests-template/tests"
)

const pngContentType = "image/png"

// Positive: Создаём объявление со всеми параметрами и проверяем, что оно доступно по id, находится поиском и фотографии загружены
func (s *TestSuite) TestCreateAdvertisementWithAllParams() {
	t := s.T()
	token := s.author.Token

	base.Precondition(fmt.Sprintf("Готовим объявление со всеми параметрами: название, описание, цена, количество и %d фотографии",
		limits.MaxPhotos))
	request := advertisement.NewAdvertisementRequest(t, limits.MaxPhotos)

	var created models.AdvertisementResponse
	base.Step("Создаём объявление через POST /advertisement", func() {
		created = advertisement.CreateAdvertisementWithCleanup(t, token, request)
		require.NotEmpty(t, created.Id, "В ответе на создание пустой id объявления")
		require.Len(t, created.Photos, len(request.PhotoPaths), "В ответе на создание не все фотографии")
		base.Result("Объявление создано, id = " + created.Id)
	})

	base.Step("Получаем объявление через GET /advertisement и сверяем параметры с переданными при создании", func() {
		received := advertisement.GetAdvertisement(t, created.Id)

		require.Equal(t, created.Id, received.Id, "id не совпадает")
		require.Equal(t, request.Title, received.Title, "Название не совпадает")
		require.Equal(t, request.Description, received.Description, "Описание не совпадает")
		require.Equal(t, request.Price, received.Price, "Цена не совпадает")
		require.Equal(t, request.Quantity, received.Quantity, "Количество не совпадает")
		require.Equal(t, s.author.User.Id, received.UserId, "Объявление принадлежит другому пользователю")
		require.Equal(t, models.PhotoIds(created.Photos), models.PhotoIds(received.Photos),
			"Фотографии объявления не совпадают с загруженными при создании")
		base.Result("Объявление по id содержит те же параметры, что переданы при создании")
	})

	base.Step("Ищем объявление по полному названию через GET /advertisements?search", func() {
		found, ok := advertisement.SearchAdvertisements(t, request.Title).FindById(created.Id)
		require.Truef(t, ok, "Объявление %s не найдено поиском по названию %q", created.Id, request.Title)
		require.Equal(t, request.Title, found.Title, "Название в результатах поиска не совпадает")
		base.Result("Объявление найдено поиском по названию")
	})

	base.Step("Проверяем фотографии через GET /advertisements/{id}/photos и скачиваем их из хранилища", func() {
		photos := advertisement.GetAdvertisementPhotos(t, token, created.Id)
		require.Equal(t, models.PhotoIds(created.Photos), models.PhotoIds(photos),
			"Список фотографий не совпадает с загруженными при создании")

		for _, photo := range photos {
			_, err := url.ParseRequestURI(photo.Url)
			require.NoErrorf(t, err, "У фотографии %s некорректный url %q", photo.Id, photo.Url)

			content := advertisement.DownloadPhoto(t, photo.Url)
			require.Equalf(t, pngContentType, http.DetectContentType(content),
				"Файл фотографии %s не является PNG-изображением", photo.Id)
		}
		base.Result(fmt.Sprintf("Все %d фотографии доступны по ручке и скачиваются из хранилища", len(photos)))
	})
}
