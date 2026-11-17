package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadBody читает и закрывает тело ответа
func ReadBody(t *testing.T, rawResponse *http.Response) []byte {
	t.Helper()
	defer rawResponse.Body.Close()

	body, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	return body
}

// RequireStatusCode проверяет статус-код ответа и возвращает тело; при несовпадении в сообщение попадает тело ответа
func RequireStatusCode(t *testing.T, rawResponse *http.Response, expectedStatusCode int) []byte {
	t.Helper()
	body := ReadBody(t, rawResponse)

	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d, получен %d. Тело ответа: %s", expectedStatusCode, rawResponse.StatusCode, body)

	return body
}

// DecodeJSON проверяет статус-код ответа и разбирает JSON-тело в структуру T
func DecodeJSON[T any](t *testing.T, rawResponse *http.Response, expectedStatusCode int) T {
	t.Helper()
	body := RequireStatusCode(t, rawResponse, expectedStatusCode)

	var result T
	require.NoErrorf(t, json.Unmarshal(body, &result), "Тело ответа не разбирается как JSON: %s", body)

	return result
}
