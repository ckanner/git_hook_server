package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
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
		fmt.Println("send job request")
		// send job request
		u, _ := url.Parse("https://ci.office.extantfuture.com/job/dev_java_common/buildWithParameters")
		query_params := u.Query()
		query_params.Set("token", "af100519383a99866be4bead138c081c")
		query_params.Set("ENV_NAME", "dev")
		u.RawQuery = query_params.Encode()
		res, _ := http.Get(u.String())
		result, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		fmt.Println(result)
	}
}

func main() {
	router := httprouter.New()
	router.POST("/hook/:service", hookHandler)
	http.ListenAndServe(":8900", router)
}
