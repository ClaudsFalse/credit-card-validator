## Introduction 🪩

In this tutorial, we will build a credit card validator using Go and the Luhn Algorithm.

The Luhn algorithm is a mathematical formula used to validate credit card numbers. We are going to write a function that implements it and use it to validate a card number.  
To achieve this, we will build a very simple server in Go which processes POST requests, extract the JSON payload with the card number, and returns a JSON response if the number is valid or not.

**⚙️What you'll need**

* Your IDE of choice
    
* Postman (alternatively you can use Curl)
    

### Let's start coding 🚀

**Project set up**

1. Create a new directory for your project, for example, `credit-card-validation`.
    
2. Open a terminal and navigate to the project directory: `cd path/to/credit-card-validation`
    
3. Inside your project directory, create a Go file where we will implement the Luhn algorithm: `touch luhn_algorithm.go`
    

**Implement the Luhn Algorithm**

To implement the algorithm, I've translated into code the rules of the formula, taken from the internet. So, there isn't much to add there but I will post the code for the function with inline comments to give context.

```go
package main

func luhnAlgorithm(cardNumber string) bool {
	// this function implements the luhn algorithm
	// it takes as argument a cardnumber of type string
	// and it returns a boolean (true or false) if the
	// card number is valid or not

	// initialise a variable to keep track of the total sum of digits
	total := 0
	// Initialize a flag to track whether the current digit is the second digit from the right.
	isSecondDigit := false

	// iterate through the card number digits in reverse order
	for i := len(cardNumber) - 1; i >= 0; i-- {
		// conver the digit character to an integer
		digit := int(cardNumber[i] - '0')

		if isSecondDigit {
			// double the digit for each second digit from the right
			digit *= 2
			if digit > 9 {
				// If doubling the digit results in a two-digit number,
				//subtract 9 to get the sum of digits.
				digit -= 9
			}
		}

		// Add the current digit to the total sum
		total += digit

		//Toggle the flag for the next iteration.
		isSecondDigit = !isSecondDigit
	}

	// return whether the total sum is divisible by 10
	// making it a valid luhn number
	return total%10 == 0
}
```

**Build your server**

We will now set up a simple server that you can run locally. Create a new Go file: `touch main.go`

We're only going to write one route, and for a bit of extra flare, we'll pass the port number as a command line argument, instead of hard-coding it.

```go
func main() {
	args := os.Args
	port := args[1]

	// Register the creditCardValidator function to handle requests at the root ("/") path.
	http.HandleFunc("/", creditCardValidator)
	fmt.Println("Listening on port:", port) 
	// Start an HTTP server listening on the specified port.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error:", err) // Print an error message if the server fails to start.
	}
}
```

In the code, the have put a placeholder function **creditCardValidator** to handle the requests that arrive at route "/" - we will now implement this function in the same file.

```go
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)
type Response struct {
	Valid bool `json:"valid"`
}

// creditCardValidator handles the credit card validation logic and JSON response.
func creditCardValidator(writer http.ResponseWriter, request *http.Request) {
	// Check if the request method is POST.
	if request.Method != http.MethodPost {
        // if not, throw an error
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
```

### **Deep dive** 🤿

**The & symbol**

Note that when we are decoding the JSON, we are using the symbol `&`before the variable `cardNumber.`

<div data-node-type="callout">
<div data-node-type="callout-emoji">💡</div>
<div data-node-type="callout-text">In Go, the <code>&amp;</code> symbol is used as the "address-of" operator. It's used to obtain the memory address of a variable. When you use <code>&amp;</code> in front of a variable, it returns a pointer to that variable's memory location.</div>
</div>

In the context of our code snippet, `cardNumber` is a struct variable that holds the credit card number extracted from the JSON payload. The `Decode` function of the `json.NewDecoder` reads JSON data from an input source (in this case, the request body `r.Body`) and tries to populate the fields of the provided struct (in this case, `cardNumber`) with the corresponding JSON values.

The `&cardNumber` part is passing a pointer to the `cardNumber` struct to the `Decode` function. This allows the `Decode` function to directly modify the fields of the `cardNumber` struct using the memory address of the variable, rather than making a copy of the struct. This is more memory-efficient and allows you to work with the actual struct instance rather than a copy!

**To "marshal"**

In Go, to "marshal" refers to the process of converting a Go data structure (such as a struct, map, or array) into its JSON representation. In other words, it's the process of encoding the Go data into a JSON format that can be sent over the network or stored in a file. That's what we are doing when using `json.marshal(response)`.

**Initialize Go modules**

1. Initialize Go modules to manage dependencies:
    
    ```go
    go mod init credit-card-validation
    ```
    

**Finally, run the project**

1. Run the project using the Go compiler. I am using port 8080, but you can use whichever port you prefer.
    
    ```go
    go run main.go luhn_algorithm.go 8080
    ```
    

To test if it's working, use Postman or Curl to send a POST request to the server at route "/", adding a credit card number to the request body:

1. Set the **headers**: Content-Type: application/json
    
2. Send a POST request to [http://localhost:8080/](http://localhost:8080/)
    
3. Set the request body to `{ "number": "4003600000000014" }`
    

If all goes well, you should receive this response `{"valid": true}`

### Conclusion 🏁

Congratulations on building a credit card validator!

Let me know, what else should I implement? And did you enjoy this blog post? Leave a comment with your thoughts, and see you next time 👋

<div data-node-type="callout">
<div data-node-type="callout-emoji">💙</div>
<div data-node-type="callout-text"><strong>You can find the full codebase on my Github. Link on my profile</strong></div>
</div>