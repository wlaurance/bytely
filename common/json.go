package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ReadJSONRequest(req *http.Request, body interface{}) error {
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawBody, body)
	if err != nil {
		return err
	}

	return nil
}

/*func WriteJSONResponse(rw http.ResponseWriter, response Responder) error {
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return err
	}
	rw.WriteHeader(response.StatusCode())
	rw.Write(res)

	return nil
}
*/
