// Package roger implements a RServe client, allowing R commands to be
// evaluted and the resulting data made available to Go applications.
//
// Connect:
//      rClient, err := roger.NewRClient("127.0.0.1", 6311)
//
// Evaluation:
//      value, err := rClient.Eval("pi")
//      fmt.Println(value) // 3.141592653589793
//
package roger
