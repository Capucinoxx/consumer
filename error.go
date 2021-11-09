package consumer

import (
	"encoding/json"
	"fmt"
)

// consumerError est une représentation d'erreur utilisée
// par le router
type consumerError struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

// Error retournes une chaine de caractères représentant l'erreur
func (e *consumerError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

// ResponseBody retourne le corps de la réponse JSON.
func (e *consumerError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

// ResponseHeaders retourne le code http ainsi que les headers
func (e *consumerError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

// Error fait la construction de l'objet `consumerError` et retourne
// le contenu dans le type `error`
func Error(err error, status int, detail string) error {
	return &consumerError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}
