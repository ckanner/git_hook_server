package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"bytes"
	"github.com/golang/go/src/fmt"
)

func hookHandler(w http.ResponseWriter, r * http.Request, ps httprouter.Params)  {
	// parse body and params
	service_name := ps.ByName("service")
	fmt.Println("service_name=", service_name)
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var f interface{}
	json.Unmarshal(result, &f)
	data_map := f.(map[string]interface{})
	// object_kind, ref
	object_kind := data_map["object_kind"].(string)
	ref := data_map["ref"].(string)
	fmt.Println("ref=", ref)

	if object_kind == "push" || object_kind == "merge_request" {
		// send job request
		url_params := url.Values{}
		url_params.Set("ENV_NAME", "dev")
		req, _ := http.NewRequest("POST", "https://ci.office.extantfuture.com/job/dev_java_common/build", bytes.NewBufferString(url_params.Encode()))
		req.SetBasicAuth("backend", "backend20166")
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
	}
}

func main() {
	router := httprouter.New()
	router.POST("/hook/:service", hookHandler)
	http.ListenAndServe(":8900", router)
}
