package temporal

import (

    "go.temporal.io/sdk/client"
)

func NewClient(address string) (client.Client, error) {
    return client.Dial(client.Options{HostPort: address})
}

// CloseSafe closes the client ignoring nil.
func CloseSafe(c client.Client) {
    if c != nil {
        _ = c.Close()
    }
}
