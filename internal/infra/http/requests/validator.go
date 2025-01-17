package requests

import (
	"encoding/json"
	"go-rest-api/internal/infra/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

type requestType interface {
	ToDomainModel() (interface{}, error)
}

func Bind[reqType requestType, domain interface{}](r *http.Request, req reqType, targetType domain) (domain, error) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error(err)
		return targetType, err
	}

	if err := v.Struct(req); err != nil {
		logger.Logger.Error(err)
		return targetType, err
	}

	d, err := req.ToDomainModel()
	if err != nil {
		logger.Logger.Error(err)
		return targetType, err
	}

	return d.(domain), nil
}
