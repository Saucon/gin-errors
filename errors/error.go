package errors

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ErrorJson struct {
	Code string            `json:"code"`
	Desc string            `json:"desc"`
	Data []ErrorSingleJson `json:"errors,omitempty"`
}

type ErrorSingleJson struct {
	Code string `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}

var errContract map[string]string

func getErrorContract() error {
	file, err := ioutil.ReadFile("errors/errorContract.json")
	json.Unmarshal([]byte(file), &errContract)

	return err
}

func ErrorMessage(err error) (httpCode int, errJson ErrorJson) {
	errGet := getErrorContract()
	if errGet != nil {
		log.Printf("error : %v", errGet)
		errJson.Code = "99"
		errJson.Desc = "General Error"
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: errGet.Error(),
		})
		httpCode = 500
		return
	}

	contract := errContract[err.Error()]
	if contract == "" {
		log.Printf("error : %v", errGet)
		errJson.Code = "99"
		errJson.Desc = "General Error"
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: err.Error(),
		})
		httpCode = 500
		return
	}

	errJson.Code = err.Error()
	errJson.Desc = errContract[err.Error()]

	return httpCode, errJson
}

func ErrorMessageArray(code string, err []error) (httpCode int, errJson ErrorJson) {
	errGet := getErrorContract()
	if errGet != nil {
		log.Printf("error : %v", errGet)
		errJson.Code = "99"
		errJson.Desc = "General Error"
		httpCode = 500
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: errGet.Error(),
		})
		return
	}

	errJson.Code = code
	errJson.Desc = "Invalid request"
	httpCode = 400
	for e := range err {
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Code: err[e].Error(),
			Desc: errContract[err[e].Error()],
		})
	}

	return httpCode, errJson
}
