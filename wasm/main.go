package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall/js"
)

const (
	stockAPIURI = "/v1/assets"
)

type assetDetail struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func main() {
	//js.Global().Set("getAssetDetail", js.FuncOf(getAssetDetail))
	getAssetDetail()
	<-make(chan bool)
}

func getAssetDetail() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in 'getAssetDetail':", r)
		}
	}()

	document := js.Global().Get("document")

	input := document.Call("querySelector", "#asset-symbol")
	output := document.Call("querySelector", "#asset-details")
	submitBtn := document.Call("querySelector", "#asset-symbol-button")

	submitBtn.Call("addEventListener", "click", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			symbol := input.Get("value").String()
			if symbol == "" {
				output.Set("innerHTML", "Please enter a symbol")
				return nil
			}

			output.Set("innerHTML", "Loading...")

			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered in 'getAssetDetail::click':", r)
					}
				}()

				var asset assetDetail

				response, err := http.Get(fmt.Sprintf("%s/%s", stockAPIURI, symbol))
				if err != nil {
					output.Set("innerHTML", err.Error())
					return
				}
				defer response.Body.Close()

				switch response.StatusCode {
				case http.StatusOK:
					err = json.NewDecoder(response.Body).Decode(&asset)
					if err != nil {
						output.Set("innerHTML", err.Error())
						return
					}
					output.Set("innerHTML", fmt.Sprintf("%s (%s)", asset.Name, asset.Symbol))
				case http.StatusNotFound:
					output.Set("innerHTML", "Invalid symbol")
				default:
					output.Set("innerHTML", fmt.Sprintf("Error: %d", response.StatusCode))
				}
			}()

			return nil
		}),
	)
}
