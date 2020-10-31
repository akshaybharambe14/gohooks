# GoHooks

GoHooks make it easy to send and consume secured web-hooks from a Go application.

## Installation

Add `github.com/averageflow/gohooks` to your `go.mod` file and then import it into where you want to be using the package by using:

```go
import (
    "github.com/averageflow/gohooks/gohooks"
)
```

## Usage

Here I will list the most basic usage for the GoHooks. If you desire more customization please read the section below for more options.

#### Sending

The most basic usage for sending is:

```go
// Data can be any type, accepts interface{}
data := []int{1, 2, 3, 4} 
// String sent in the GoHook that helps identify actions to take with data
resource := "int-list-example"
// Secret string that should be common to sender and receiver
// in order to validate the GoHook signature
saltSecret := "0014716e-392c-4120-609e-555e295faff5"

hook := &gohooks.GoHook{}
hook.Create(data, resource, salt)

// Will return *http.Response and error
resp, err := hook.Send("www.example.com/hooks")
```

#### Receiving

The most basic usage for receiving is:

```go
type MyWebhook struct {
    Resource string `json:"resource"`
    Data []int `json:"data"`
}

var request MyWebhook
// Assuming you use Gin Gonic, otherwise unmarshall JSON yourself.
_ = c.ShouldBindJSON(&request)

// Shared secret with sender
saltSecret := "0014716e-392c-4120-609e-555e295faff5"
// Assuming you use Gin Gonic, obtain signature header value
receivedSignature := c.GetHeader(gohooks.DefaultSignatureHeader)

// Verify validity of GoHook
isValid := gohooks.IsGoHookValid(requestBody, receivedSignature, saltSecret)
// Decide what to do if GoHook is valid or not.
```

## Customization

GoHooks use the custom header `X-GoHooks-Verification` to send the encrypted SHA string. You can customize this header by initializing the GoHook struct with the custom option `SignatureHeader`. 

Example: 
```go
hook := &gohooks.GoHook{ SignatureHeader: "X-Example-Custom-Header" }
```

GoHooks are by default not verifying the receiver's SSL certificate validity. If you desire this behaviour then enable it by initializing the GoHook struct with the custom option `IsSecure`.

Example: 
```go
hook := &gohooks.GoHook{ IsSecure: true }
```

GoHooks will by default be sent via a `POST` request. If you desire to use a different HTTP method, amongst the allowed `POST`, `PUT`, `PATCH`, `DELETE`, then feel free to pass that option when initializing the GoHook struct, with `PreferredMethod`. Any other value will make the GoHook default to a `POST` request.

Example: 
```go
hook := &gohooks.GoHook{ PreferredMethod: http.MethodDelete }
```