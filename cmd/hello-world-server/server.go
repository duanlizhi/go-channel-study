package main

import (
	"net/http"
	"fmt"
)

/**
 * <p>Description: (helloworld server入口文件) </p>
 * @author lizhi_duan
 * @date 2018/10/28 10:27
 */
func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer,"Hello world! %s", request.FormValue("name"))
	})
	http.ListenAndServe(":8888",nil)
}
