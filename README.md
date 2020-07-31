# Sesame API client for Go

sesame-client-go is the unofficial Sesame API client for Go programming language.

This client supports the version 3 API.

- [CANDY HOUSE Developer Reference(Official)](https://docs.candyhouse.co/)

Installation
------------

```
go get -u github.com/tukaelu/sesame-client-go
```

Example
-------

```

func main() {
    ctx := context.Background()
    cli := opensesame.NewSesameAPI("YOUR_API_KEY")

    // Get Sesame list
	devices, err := cli.GetList(ctx)
	if err != nil {
        log.Fatal(err)
		return
	}

    // Get Sesame status
    stat, err := api.GetStatus(ctx, "DEVICE_ID")
	if err != nil {
        log.Fatal(err)
		return
	}
    fmt.Printf("Battery: %d%", stat.Battery) // Battery: 80%

    // Control Sesame
    ctrl, err := api.Control(ctx, "DEVICE_ID", "lock")
	if err != nil {
        log.Fatal(err)
		return
	}
    fmt.Printf("Task ID: %s", ctrl.TaskID) // Task ID: 01234567-890a-bcde-f012-34567890abcd

    // Query Execution Result
    result, err := api.GetExecutionResult(ctx, ctrl.TaskID)
	if err != nil {
        log.Fatal(err)
		return
	}
    fmt.Printf("Status: %s", ctrl.Status) // Status: processing
}

```
