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

func parseRequest(r *http.Request, service_name string) (env_name, branch_name string) {
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	var f interface{}
	json.Unmarshal(result, &f)
	data_map := f.(map[string]interface{})
	// object_kind, ref
	object_kind := data_map["object_kind"].(string)
	if strings.EqualFold(object_kind, "push") {
		ref_branch := strings.Replace(data_map["ref"].(string), "refs/heads/", "", 1)
		if strings.EqualFold(ref_branch, "develop") {
			env_name = "dev"
			branch_name = "develop"
		} else if strings.EqualFold(ref_branch, "master") {
			// TODO add
			//env_name = "formal"
			//branch_name = "master
		}
	} else if strings.EqualFold(object_kind, "merge_request") {
		object_attributes := data_map["object_attributes"].(map[string]interface{})
		target_branch := object_attributes["target_branch"].(string)
		merge_status := object_attributes["merge_status"].(string)
		if strings.EqualFold(merge_status, "merged") {
			if strings.EqualFold(target_branch, "develop") {
				env_name = "dev"
				branch_name = "develop"
			} else if strings.EqualFold(target_branch, "master") {
				// TODO add
				//env_name = "formal"
				//branch_name = "master
			}
		}
	} else if strings.EqualFold(object_kind, "tag_push") {
		env_name = "dev"
		branch_name = strings.Replace(data_map["ref"].(string), "refs/tags/", "", 1)
	}
	log.Println("parseRequest object_kind is " + object_kind + ", env_name is " + env_name + ", branch_name is" + branch_name)
	return env_name, branch_name
}

func sendBuildJob(req_url string)  {
	req, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.SetBasicAuth("backend", "d7961737278945b6d9a506a99c23b67e")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("err err")
	}
	defer resp.Body.Close()
}

func hookHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	service_name := ps.ByName("service")
	env_name, branch_name := parseRequest(r, service_name)
	if !strings.EqualFold(env_name, "") && !strings.EqualFold(branch_name, "") {
		req_url := fmt.Sprintf("http://publish.extantfuture.com/job/%s_%s/buildWithParameters?token=d7961737278945b6d9a506a99c23b67e&BRANCH=%s&TRIGGER=auto", env_name, service_name, branch_name)
		sendBuildJob(req_url)
	}
}

func main() {
	router := httprouter.New()
	router.POST("/hook/:service", hookHandler)
	http.ListenAndServe(":8900", router)
}
