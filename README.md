# throt
Golang IO throttling 

Usage:
```go
package main
 
 import "github.com/slavash/throt"
 
 // set bandwidth limit per server
 globalLimiter := throt.NewLimiter(globalRateLimit, burst)

 ...
 // connection handler - serving the file
 fd, err := os.Open(fileName)
 
 if err != nil {
     return err
 }
 
 // set bandwidth limit per connection
 connLimiter := throt.NewLimiter(int(connLimit), burst)
 ...

 // decorating the reader
 r1 := throt.NewReader(ctx, fd)
 r1.ApplyLimits(connLimiter)

 // decorating the reader again...
 r2 = throt.NewReader(ctx, r1)
 r2.ApplyLimits(globalLimiter)

 sent, err = io.Copy(c, r2)
 
 // The same may be done with io.Writer
 ...
```