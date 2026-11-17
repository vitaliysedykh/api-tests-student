package models

type AdvertisementResponse struct {
	Id          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       *float64        `json:"price"`
	Quantity    *int            `json:"quantity"`
	UserId      string          `json:"user_id"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Photos      []PhotoResponse `json:"photos"`
}

type PhotoResponse struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

type AdvertisementListResponse struct {
	Items      []AdvertisementResponse `json:"items"`
	HasNext    bool                    `json:"has_next"`
	NextCursor string                  `json:"next_cursor"`
	Total      int                     `json:"total"`
}

// FindById ищет объявление в списке по id
func (l AdvertisementListResponse) FindById(id string) (AdvertisementResponse, bool) {
	for _, item := range l.Items {
		if item.Id == id {
			return item, true
		}
	}
	return AdvertisementResponse{}, false
}

// PhotoIds возвращает id фотографий в порядке следования
func PhotoIds(photos []PhotoResponse) []string {
	ids := make([]string, 0, len(photos))
	for _, photo := range photos {
		ids = append(ids, photo.Id)
	}
	return ids
}
