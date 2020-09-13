# throt
Golang IO throttling 

Usage:
```go
package main
 
 import "github.com/slavash/throt"
 
 // in the server start
 // set bandwidth limit per server
 globalLimiter := throt.NewLimiter(globalRateLimit)

 ...
 // in connection handler
 
 fd, err := os.Open(fileName)
 
 if err != nil {
     return err
 }
 
 // set bandwidth limit per connection
 connLimiter := throt.NewLimiter(connRateLimit)

 
 reader := throt.NewReader(ctx, fd)
 reader.ApplyLimits(connLimiter, globalLimiter)
 sent, err = io.Copy(c, reader)
 
 ...
 
 // The same may be done with writer:
 
 writer := throt.NewWriter(ctx, c)
 writer.ApplyLimits(connLimiter, globalLimiter)
 sent, err = io.Copy(writer, fd)
 
 sent, err = io.Copy(c, reader)

 ...
```