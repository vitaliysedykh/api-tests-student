package apiErrors

// Значения поля error в ответах API
const (
	ErrorInvalidRequest = "invalid_request"
	ErrorUnauthorized   = "unauthorized"
	ErrorForbidden      = "forbidden"
	ErrorNotFound       = "not_found"
	ErrorAlreadyExists  = "already_exists"
	ErrorLimitReached   = "limit_reached"
)

// Значения поля message в ответах API
const (
	MessageAuthorizationRequired = "Authorization header required"
	MessageAuthorizationInvalid  = "Invalid or expired token"

	MessageTitleRequired            = "title is required"
	MessageTitleTooLongFormat       = "title must not exceed %d characters provided %d"
	MessageDescriptionRequired      = "description is required"
	MessageDescriptionLengthFormat  = "description must be between 1 and %d characters"
	MessagePriceOutOfRangeFormat    = "price must be between 0 and %d inclusive"
	MessageQuantityOutOfRangeFormat = "quantity must be between %d and %d inclusive"
	MessagePhotoRequired            = "at least 1 photo is required"
	MessageTooManyPhotosFormat      = "maximum %d photos allowed"

	MessageAdvertisementNotFound   = "Advertisement not found"
	MessageFavoriteNotFound        = "Favorite entry not found"
	MessageFavoriteAlreadyExists   = "Advertisement is already in favorites"
	MessageCannotFavoriteOwn       = "Cannot add own advertisement to favorites"
	MessageInvalidUuidLengthFormat = "invalid UUID length: %d"
	MessageFavoritesLimitFormat    = "Favorites limit of %d reached"
)
