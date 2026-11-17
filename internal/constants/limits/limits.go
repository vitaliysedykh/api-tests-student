package limits

// Бизнес-ограничения TestBoard API, подтверждённые ответами сервера
const (
	MaxTitleLength       = 50
	MaxDescriptionLength = 500
	MaxPrice             = 1000000
	MinQuantity          = 1
	MaxQuantity          = 100
	MaxPhotos            = 3

	// FavoritesLimit максимальное число объявлений в избранном одного пользователя
	FavoritesLimit = 50
	// FavoritesMaxPageSize наибольший допустимый page_size для GET /favorites
	FavoritesMaxPageSize = 50
	// CatalogMaxPageSize наибольший допустимый limit для GET /advertisements
	CatalogMaxPageSize = 50
)
