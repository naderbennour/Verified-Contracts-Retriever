package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiJson struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Sourcecode string `json:"SourceCode"`
	} `json:"result"`
}

func CollectSourceCode(apiKey, contractAddr string) (string, error) {
	// Construct the url for the api call using the contract address and api key
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s",
		contractAddr, apiKey)
	// Make a GET request to the url
	resp, err := http.Get(url)
	if err != nil {
		// Return an error if the GET request fails
		return "", err
	}
	// Declare an ApiJson struct to hold the json response
	apiJson := ApiJson{}
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Return an error if reading the response body fails
		return "", err
	}
	// Close the response body
	resp.Body.Close()
	// Unmarshal the json response into the ApiJson struct
	err = json.Unmarshal([]byte(body), &apiJson)
	if err != nil {
		// Return an error if unmarshalling the json response fails
		return "", err
	}
	if len(apiJson.Result) == 0 {
		// Return an error if the result is empty
		return "", errors.New("collect result empty error")
	}
	if len(apiJson.Result) > 1 {
		// Print a warning if the result has more than one source
		fmt.Printf("Warning: Result of has 2 or more srcs, address = \"%s\"\n", contractAddr)
	}
	// Return the source code of the contract
	return apiJson.Result[0].Sourcecode, nil
}

func UnmarshalSourceCode(sourceCodeJson string) ([]string, error) {
	// Check if the source code json string starts with "{{"
	if !strings.HasPrefix(sourceCodeJson, "{{") {
		// If it does not, return the json string as a single element slice
		return []string{sourceCodeJson}, nil
	}
	// Declare a variable to hold the unmarshaled json
	var src interface{}
	// Remove the first and last two characters from the json string
	sourceCodeJson = sourceCodeJson[1 : len(sourceCodeJson)-1]
	// Unmarshal the json string into the src variable
	err := json.Unmarshal([]byte(sourceCodeJson), &src)
	if err != nil {
		// Return an error if json unmarshal fails
		return nil, err
	}
	// Assert src type to map[string]interface{} and access the "sources" key
	src = src.(map[string]interface{})["sources"]
	// Assert src type to map[string]interface{} again
	src = src.(map[string]interface{})
	var sols []string
	// Assert src type to map[string]interface{}
	srcMap, ok := src.(map[string]interface{})
	if !ok {
		// Return an error if assertion fails
		return nil, errors.New("Error while asserting src type to map[string]interface{}")
	}
	// Iterate through the srcMap and append the "content" value to the sols slice
	for _, v := range srcMap {
		sols = append(sols, v.(map[string]interface{})["content"].(string))
	}
	return sols, nil
}

func main() {
	apiKey := "Your Etherscan API"
	contractAddr := "Address you want to get the contracts for"

	// Collect source code JSON of the contract address
	sourceCodeJson, err := CollectSourceCode(apiKey, contractAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Unmarshal the source code JSON string of a contract
	sols, err := UnmarshalSourceCode(sourceCodeJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save all the solidity files in the contract as individual files
	for i, sol := range sols {
		fileName := fmt.Sprintf("contract_%d.sol", i)
		err := ioutil.WriteFile(fileName, []byte(sol), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Saved %s\n", fileName)
	}
}
