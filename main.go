package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type chatRequest struct {
	Message string `json:"message"`
}

type chatResponse struct {
	Text string `json:"text"`
}

// API endpoint
const API_URL = "https://api.cohere.com/v1/chat"

func getGeneratedResponse(prompt string) (string, error) {
	//Set a prompt to chatRequest type
	chatReq := &chatRequest{
		Message: prompt,
	}

	////Marshal Go struct into Json
	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		log.Printf("Failed to Marshal: %v", err)
		return "", err
	}

	//Create Http request struct with request method, endpoint and request body
	req, err := http.NewRequestWithContext(context.Background(), "POST", API_URL, bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("Failed to create http request struct: %v", err)
		return "", err
	}

	//Add necessary headers, including the API key for authorization
	apiKey := os.Getenv("COHERE_API_KEY")
	if apiKey == "" {
		log.Printf("Failed to get API KEY: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Execute http request to cohere and get response
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get http response: %v", err)
		return "", err
	}

	//Check if http status code is ok
	if res.StatusCode != http.StatusOK {
		err := errors.New("Unexpected status code")
		log.Printf("Failed to get expected status code: %v :%d", err, res.StatusCode)
		return "", err
	}

	//Read http response body
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		return "", err
	}

	//Unmarshal Json response into Go struct
	chatRes := &chatResponse{}
	err = json.Unmarshal(body, chatRes)
	if err != nil {
		log.Printf("Failed to unmarshal: %v", err)
		return "", err
	}

	return chatRes.Text, nil
}

func main() {
	words := [3]string{"nonchalant", "reckon", "appalled"}
	prompt := fmt.Sprintf("Please create an English example sentence using following words: %s, %s, %s",
		words[0], words[1], words[2])

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Prompt ++++++")
	fmt.Println(prompt)

	fmt.Println("")
	fmt.Println("")

	fmt.Println("++++++ Generated response ++++++")
	response, err := getGeneratedResponse(prompt)
	if err != nil {
		log.Fatalf("Failed to get generated response from Cohere API: %v", err)
	}
	fmt.Println(response)

	fmt.Println("")
	fmt.Println("")
}
