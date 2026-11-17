package favorites

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: Проверяем, что добавленное объявление появляется в избранном, а удалённое из избранного — пропадает
func (s *TestSuite) TestAddAndRemoveFavorite() {
	ad := s.createSellerAdvertisement()

	base.Step("Покупатель добавляет объявление в избранное через POST /favorites", func() {
		favorite := favorites.AddFavorite(s.T(), s.buyer.Token, ad.Id)
		require.Equal(s.T(), ad.Id, favorite.Advertisement.Id, "В избранное добавлено другое объявление")
	})

	base.Step("Проверяем, что объявление появилось в списке избранного GET /favorites", func() {
		s.requireInFavorites(s.buyer.Token, ad)
	})

	base.Step("Покупатель удаляет объявление из избранного через DELETE /favorites/{ad_id}", func() {
		favorites.RemoveFavorite(s.T(), s.buyer.Token, ad.Id)
	})

	base.Step("Проверяем, что объявление пропало из списка избранного", func() {
		s.requireNotInFavorites(s.buyer.Token, ad.Id)
	})
}

// Positive: Проверяем, что объявление, удалённое владельцем, пропадает из избранного покупателя
func (s *TestSuite) TestDeletedAdvertisementDisappearsFromFavorites() {
	base.Precondition("Продавец создаёт объявление без автоудаления: удаление — шаг этого теста")
	ad := advertisement.CreateAdvertisement(s.T(), s.seller.Token, advertisement.NewAdvertisementRequest(s.T(), 1))

	base.Step("Покупатель добавляет объявление в избранное", func() {
		favorites.AddFavorite(s.T(), s.buyer.Token, ad.Id)
		s.requireInFavorites(s.buyer.Token, ad)
	})

	base.Step("Продавец удаляет объявление через DELETE /advertisement", func() {
		advertisement.DeleteAdvertisement(s.T(), s.seller.Token, ad.Id)
	})

	base.Step("Проверяем, что объявление пропало из избранного покупателя", func() {
		s.requireNotInFavorites(s.buyer.Token, ad.Id)
	})
}
