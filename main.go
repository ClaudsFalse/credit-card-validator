// main.go
package main

import (
	"encoding/json"
	"net/http"
)

// Response is a struct for representing the JSON response.
type Response struct {
	Valid bool `json:"valid"` // Valid field indicates whether the card number is valid.
}

// creditCardValidator handles the credit card validation logic and JSON response.
func creditCardValidator(writer http.ResponseWriter, request *http.Request) {
	// Check if the request method is POST.
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Create a struct to hold the incoming JSON payload.
	var cardNumber struct {
		Number string `json:"number"` // Number field holds the credit card number.
	}

	// Decode the JSON payload from the request body into the cardNumber struct.
	err := json.NewDecoder(request.Body).Decode(&cardNumber)
	if err != nil {
		http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate the credit card number using the Luhn algorithm.
	isValid := luhnAlgorithm(cardNumber.Number)

	// Create a response struct with the validation result.
	response := Response{Valid: isValid}

	// Marshal the response struct into JSON format.
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, "Error creating response", http.StatusInternalServerError)
		return
	}

	// Set the content type header to indicate JSON response.
	writer.Header().Set("Content-Type", "application/json")

	// Write the JSON response back to the client.
	writer.Write(jsonResponse)
}

func main() {
	// Register the creditCardValidator function to handle requests at the root ("/") path.
	http.HandleFunc("/", creditCardValidator)

	// Start an HTTP server listening on port 8080.
	http.ListenAndServe(":8080", nil)
}
