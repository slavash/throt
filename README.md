# throt
Golang IO throttling 

Usage:
```go
package main
 
 import "github.com/slavash/throt"
 
 ...
 
 fd, err := os.Open(fileName)
 
 if err != nil {
     return err
 }
 
 // set bandwidth limit per connection
 connLimiter := rate.NewLimiter(rate.Limit(rateLimit), int(rateLimit))
 // set bandwidth limit per server
 globalLimiter := ctx.Value("rateLimit").(*rate.Limiter)
 
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