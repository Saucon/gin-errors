package errors

import (
	"encoding/json"
	"io/ioutil"
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

func getErrorContract() (err error) {
	var file []byte

	if errContract == nil {
		file, err = ioutil.ReadFile("errors/errorContract.json")
		json.Unmarshal([]byte(file), &errContract)
	}

	return err
}

func ErrorMessage(custHttpCode int, err error) (httpCode int, errJson ErrorJson) {
	errGet := getErrorContract()
	if errGet != nil {
		errJson.Code = "99"
		errJson.Desc = "Unexpected Error"
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: errGet.Error(),
		})
		httpCode = 500
		return
	}

	contract := errContract[err.Error()]
	if contract == "" {
		errJson.Code = "99"
		errJson.Desc = "Unexpected Error"
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: err.Error(),
		})
		httpCode = 500
		return
	}

	errJson.Code = err.Error()
	errJson.Desc = errContract[err.Error()]
	httpCode = custHttpCode

	return httpCode, errJson
}

func ErrorMessageArray(custHttpCode int, code string, err []error) (httpCode int, errJson ErrorJson) {
	errGet := getErrorContract()
	if errGet != nil {
		errJson.Code = "99"
		errJson.Desc = "Unexpected Error"
		httpCode = 500
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Desc: errGet.Error(),
		})
		return
	}

	contract := errContract[code]
	if contract == "" {
		errJson.Code = "99"
		errJson.Desc = "Unexpected Error"
		for e := range err {
			errJson.Data = append(errJson.Data, ErrorSingleJson{
				Code: err[e].Error(),
				Desc: errContract[err[e].Error()],
			})
		}
		httpCode = 500
		return
	}

	errJson.Code = code
	errJson.Desc = errContract[code]
	httpCode = custHttpCode
	for e := range err {
		errJson.Data = append(errJson.Data, ErrorSingleJson{
			Code: err[e].Error(),
			Desc: errContract[err[e].Error()],
		})
	}

	return httpCode, errJson
}
