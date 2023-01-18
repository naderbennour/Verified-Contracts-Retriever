# Verified Contracts Retriever

This Go script uses the Etherscan API to retrieve the source code of a smart contract by its address and unmarshals the JSON response. 

## Usage

1. Assign your Etherscan API key to the `apiKey` variable and the contract address you want to retrieve the source code for to the `contractAddr` variable in the `main` function.
2. Run the script with `go run script_name.go`
3. The source code and any errors will be returned. 

## Functions

```
CollectSourceCode(apiKey, contractAddr string) (string, error)
```

This function takes in an API key and a contract address as arguments, and makes a GET request to the Etherscan API to retrieve the source code. The response is then parsed into a struct and the source code is returned.

```
UnmarshalSourceCode(sourceCodeJson string) ([]string, error)
```

This function takes in the source code json string and checks if it starts with "{{" (this is done to handle a special case where the json string returned is wrapped in double curly braces). If it does, it removes the first and last two characters, unmarshals the json string and asserts the type to a map. It then iterates through the map and appends the "content" value to a slice. 

