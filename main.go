package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cep           = "01153000"
	BrasilAPI_URL = "https://brasilapi.com.br/api/cep/v1/" + cep
	ViaCEP_URL    = "https://viacep.com.br/ws/" + cep + "/json"
)

func main() {
	channelBrasilAPI := make(chan BrasilAPIStruct)
	channelViaCEP := make(chan ViaCEPStruct)

	go getBrasilAPI(channelBrasilAPI)
	go getViaCEP(channelViaCEP)

	select {
	case response := <-channelBrasilAPI:
		fmt.Printf("BrasilAPI: %+v\n", response)
	case response := <-channelViaCEP:
		fmt.Printf("ViaCEP: %+v\n", response)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}

type BrasilAPIStruct struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func getBrasilAPI(channel chan<- BrasilAPIStruct) {
	response, err := fetchAPI(BrasilAPI_URL)
	if err != nil {
		panic(err)
	}
	data := response.(map[string]interface{})
	channel <- BrasilAPIStruct{
		Cep:          data["cep"].(string),
		State:        data["state"].(string),
		City:         data["city"].(string),
		Neighborhood: data["neighborhood"].(string),
		Street:       data["street"].(string),
		Service:      data["service"].(string),
	}
}

type ViaCEPStruct struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func getViaCEP(channel chan<- ViaCEPStruct) {
	response, err := fetchAPI(ViaCEP_URL)
	if err != nil {
		panic(err)
	}
	data := response.(map[string]interface{})
	channel <- ViaCEPStruct{
		Cep:         data["cep"].(string),
		Logradouro:  data["logradouro"].(string),
		Complemento: data["complemento"].(string),
		Bairro:      data["bairro"].(string),
		Localidade:  data["localidade"].(string),
		Uf:          data["uf"].(string),
		Ibge:        data["ibge"].(string),
		Gia:         data["gia"].(string),
		Ddd:         data["ddd"].(string),
		Siafi:       data["siafi"].(string),
	}
}

func fetchAPI(url string) (response interface{}, err error) {
	request, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
