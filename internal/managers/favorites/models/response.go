package models

import (
	advertisementModels "api-tests-template/internal/managers/advertisement/models"
)

type FavoriteItem struct {
	FavoriteId    string                                    `json:"favorite_id"`
	AddedAt       string                                    `json:"added_at"`
	Advertisement advertisementModels.AdvertisementResponse `json:"advertisement"`
}

type FavoriteListResponse struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

// FindByAdvertisementId ищет запись избранного по id объявления
func (l FavoriteListResponse) FindByAdvertisementId(advertisementId string) (FavoriteItem, bool) {
	for _, item := range l.Items {
		if item.Advertisement.Id == advertisementId {
			return item, true
		}
	}
	return FavoriteItem{}, false
}
