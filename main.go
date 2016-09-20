package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"strings"
	"encoding/json"
)

func hookHandler(w http.ResponseWriter, r * http.Request, ps httprouter.Params)  {
	fmt.Println(r.URL.Path)
	fmt.Println(ps.ByName("service"))
	result, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var f interface{}
	json.Unmarshal(result, &f)
	m := f.(map[string]interface{})
	fmt.Println(strings.EqualFold(m["ref"].(string), "refs/heads"))

	fmt.Println(string(result))
	fmt.Println(r.Header.Get("X-Gitlab-Event"))
}

func main() {
	router := httprouter.New()
	router.POST("/hook/:service", hookHandler)
	http.ListenAndServe(":8900", router)
}
