package api_runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-tests-template/internal/utils"
)

const (
	maxRateLimitRetries = 5
	rateLimitRetryDelay = 500 * time.Millisecond
)

// rateLimitRetryTransport повторяет запрос с нарастающей паузой, если сервер ответил 503:
// так TestBoard отклоняет запросы сверх лимита частоты (около 10 в секунду с одного адреса)
type rateLimitRetryTransport struct {
	next http.RoundTripper
}

// newRateLimitRetryClient возвращает http-клиент, переживающий ограничение частоты запросов TestBoard
func newRateLimitRetryClient() *http.Client {
	return &http.Client{Transport: rateLimitRetryTransport{next: http.DefaultTransport}}
}

func (rt rateLimitRetryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	// Тело читается заранее: apitest не передаёт GetBody, без копии повтор ушёл бы с пустым телом
	body, err := readRequestBody(request)
	if err != nil {
		return nil, err
	}

	for attempt := 1; ; attempt++ {
		attemptRequest := request.Clone(request.Context())
		if request.Body != nil {
			attemptRequest.Body = io.NopCloser(bytes.NewReader(body))
		}

		response, err := rt.next.RoundTrip(attemptRequest)
		if err != nil || response.StatusCode != http.StatusServiceUnavailable || attempt > maxRateLimitRetries {
			return response, err
		}
		_ = response.Body.Close()

		delay := rateLimitRetryDelay * time.Duration(attempt)
		utils.LogWithLabelAndTimeStamp("RateLimit", fmt.Sprintf("%s %s: сервер ответил 503, повтор %d из %d через %s",
			request.Method, request.URL.Path, attempt, maxRateLimitRetries, delay))
		time.Sleep(delay)
	}
}

// readRequestBody читает и закрывает тело запроса, чтобы его можно было отправить повторно
func readRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	defer request.Body.Close()
	return io.ReadAll(request.Body)
}
