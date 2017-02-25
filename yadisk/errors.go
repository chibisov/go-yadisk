package yadisk

import "fmt"

// APIError error can occur if the request was formed incorrectly,
// the specified resource does not exist on the server,
// the server is not working, and so on.
// All errors are returned with the corresponding HTTP response codes.
// All the possible response codes and explanations
// are given in Polygon https://tech.yandex.ru/disk/poligon/
// Errors are additionally described by a JSON object
// https://tech.yandex.com/disk/api/reference/response-objects-docpage/#error
type APIError struct {
	Description string `json:"description"`
	Code        string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Code == "" && e.Description == "" {
		return "Yandex.Disk API error"
	}
	return fmt.Sprintf(
		"Yandex.Disk API error. Code: %s. Description: %s.",
		e.Code,
		e.Description,
	)
}
