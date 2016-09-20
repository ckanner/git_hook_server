package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"log"
)

func hookHandler(w http.ResponseWriter, r * http.Request, ps httprouter.Params)  {
	// parse body and params
	service_name := ps.ByName("service")
	log.Println("hookHandler service_name is " + service_name)
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var f interface{}
	json.Unmarshal(result, &f)
	data_map := f.(map[string]interface{})
	// object_kind, ref
	object_kind := data_map["object_kind"].(string)
	log.Println("hookHandler object_kind is " + object_kind)

	var env_name = ""
	if strings.EqualFold(object_kind, "push") {
		ref_branch := data_map["ref"].(string)
		if strings.EqualFold(ref_branch, "refs/heads/feature/docker") {
			env_name = "dev"
		} else if strings.EqualFold(ref_branch, "refs/heads/master") {
			env_name = "formal"
		}
	} else if strings.EqualFold(object_kind, "merge_request") {
		object_attributes := data_map["object_attributes"].(map[string]interface{})
		target_branch := object_attributes["target_branch"].(string)
		merge_status := object_attributes["merge_status"].(string)
		if strings.EqualFold(merge_status, "merged") {
			if strings.EqualFold(target_branch, "feature/docker") {
				env_name = "dev"
			} else if strings.EqualFold(target_branch, "master") {
				env_name = "formal"
			}
		}
	}

	log.Println("hookHandler env_name is " + env_name)
	if !strings.EqualFold(env_name, "") {
		req_url := fmt.Sprintf("https://ci.office.extantfuture.com/job/%s_%s/buildWithParameters?token=af100519383a99866be4bead138c081c&ENV_NAME=%s", env_name, service_name, env_name)
		req, err := http.NewRequest("GET", req_url, nil)
		if err != nil {
			fmt.Println(err)
		}
		req.SetBasicAuth("backend", "af100519383a99866be4bead138c081c")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("err err")
		}
		defer resp.Body.Close()
	}
}

func main() {
	router := httprouter.New()
	router.POST("/hook/:service", hookHandler)
	http.ListenAndServe(":8900", router)
}
