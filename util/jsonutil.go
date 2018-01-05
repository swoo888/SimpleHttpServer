package util

import (
	"encoding/json"
	"net/http"
	"io"
	"log"
)

// DecodeJSONBodyOrBadRequest is a json decode util function to decode request body into v
// v is a pointer to some data
// if decoding failed, it will write StatusBadRequest to response
func DecodeJSONBodyOrBadRequest(resp http.ResponseWriter, reqBody io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(reqBody)
	defer closeReqBody(reqBody)
	err := decoder.Decode(v)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
	}
	return err
}

// WriteJSONRespOrInternalServerError writes json data of v to response
// if json marshaling of v failed, it will write StatusInternalServerError to response
func WriteJSONRespOrInternalServerError(resp http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	_, err = resp.Write(js)
	if err != nil {
		log.Println("error writing json model.", err)
	}
}

func closeReqBody(reqBody io.ReadCloser) {
	err := reqBody.Close()
	if err != nil {
		log.Println("error closing request body.", err)
	}
}
