# Roger

[![Join the chat at https://gitter.im/senseyeio/roger](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/senseyeio/roger?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Roger is a Go [RServe](http://www.rforge.net/Rserve/) client, allowing the capabilities of [R](http://www.r-project.org/) to be used from Go applications.

The communication between Go and R is via TCP. It is thread safe and supports long running R operations synchronously or asynchronously (using channels).

Roger is a **work in progress**, currently supporting string, string array, double and double array response types. With the use of JSON, this capability can be used to transfer any serializable object.

##Example

```go
package main

import (
	"fmt"

	"github.com/senseyeio/roger"
)

func main() {
	rClient, err := roger.NewRClient("127.0.0.1", 6311)
	if err != nil {
		fmt.Println("Failed to connect")
		return
	}

	value, err := rClient.Eval("pi")
	if err != nil {
		fmt.Println("Command failed: " + err.Error())
	} else {
		fmt.Println(value) // 3.141592653589793
	}

	helloWorld, _ := rClient.Eval("as.character('Hello World')")
	fmt.Println(helloWorld) // Hello World

	arrChan := rClient.Evaluate("Sys.sleep(5); c(1,1)")
	arrResponse := <-arrChan
	arr, _ := arrResponse.GetResultObject()
	fmt.Println(arr) // [1, 1]
}
```

##Setup
Rserve should be installed and started from R:

```R
install.packages("Rserve")
require('Rserve')
Rserve()
```

More information is available on [RServe's website](https://www.rforge.net/Rserve/doc.html).

If you would like to exploit the current R environment from go, start RServe using the following command:

```R
install.packages("Rserve")
require('Rserve')
run.Rserve()
```

Install Roger using:

```
go get github.com/senseyeio/roger
```

## Contributing
Issues, pull requests and questions are welcomed. If required, assistance can be found in the project's [gitter chat room](https://gitter.im/senseyeio/roger).

### Pull Requests

 - Fork the repository
 - Make changes
 - Raise pull request
