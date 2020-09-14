## Example
This is straightforward implementation to demonstrate the library only.  
Please don't consider it as a code example.  

#### Usage:

Start server:
```bash
$ go run example.go
```

Run client:  

```bash
echo "get file.exe" | nc localhost 7777 > /dev/null
```
assuming the file.exe exists in the local directory

Change limit (bytes/s):
```bash
echo "setl 500000" | nc localhost 7777
```
