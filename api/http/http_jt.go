package apihttp

import (
	"bytes"
	"io/ioutil"
	"log"
	"mantap2/config"
	"net/http"
)

func SendToHttpAPI(IDReader string, SNRFID string) {
	var err error
	var client = &http.Client{}
	var jsonStr = []byte(`{"IDReader":"` + IDReader + `", "SNRFID":"` + SNRFID + `"}`)

	//var dataResponse string
	// var payload = bytes.NewBufferString(param.Encode())
	var payload = bytes.NewBuffer(jsonStr)

	log.Println(config.Api_server + "/api/trx")

	// req, err := http.NewRequest("POST", api_server+"/api/trx", bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "MANTAP2")
	// req.Header.Set("Content-Type", "application/json")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	request, err := http.NewRequest("POST", config.Api_server+"/api/trx", payload)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	request.Header.Set("X-Custom-Header", "MANTAP2")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer response.Body.Close()

	log.Println("response Status:", response.Status)
	log.Println("response Headers:", response.Header)

	body, _ := ioutil.ReadAll(response.Body)
	log.Println("response Body:", string(body))

	// err = json.NewDecoder(response.Body).Decode(&dataResponse)
	// if err != nil {
	// 	log.Printf("Error: %s\n", err)
	// 	return
	// }
}
