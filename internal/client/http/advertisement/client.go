package advertisement

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/steinfletcher/apitest"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	advertisementModels "api-tests-template/internal/managers/advertisement/models"
)

// HttpPostAdvertisement выполняет POST-запрос создания объявления: поля и фотографии передаются как multipart/form-data
func HttpPostAdvertisement(t *testing.T, token string, request advertisementModels.CreateAdvertisementRequest) *http.Response {
	apiRequest := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description)

	if request.Price != nil {
		apiRequest = apiRequest.MultipartFormData("price", strconv.FormatFloat(*request.Price, 'f', -1, 64))
	}
	if request.Quantity != nil {
		apiRequest = apiRequest.MultipartFormData("quantity", strconv.Itoa(*request.Quantity))
	}
	if len(request.PhotoPaths) > 0 {
		apiRequest = apiRequest.MultipartFile("photos", request.PhotoPaths...)
	}

	return apiRequest.Expect(t).End().Response
}

// HttpGetAdvertisement выполняет GET-запрос получения объявления по id
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisements выполняет GET-запрос списка объявлений каталога с параметрами поиска, сортировки и пагинации
func HttpGetAdvertisements(t *testing.T, params map[string]string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		QueryParams(params).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос списка фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).
		Expect(t).
		End().Response
}

// HttpGetPhotoFile скачивает файл фотографии по абсолютному url хранилища; хост не из API_URL, поэтому без общего раннера
func HttpGetPhotoFile(t *testing.T, photoUrl string) *http.Response {
	return apitest.New().EnableNetworking().Get(photoUrl).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос удаления объявления по id
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
