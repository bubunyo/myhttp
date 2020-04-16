# myhttp - fetch and hash http url responses 

`myhttp` fetches a url over the network and returns an `md5` hash of the response gotten via the network call.

by Bubunyo Nyavor.

# Dependencies
- Go +1.31.1 - The golang toolchain needs to be installed to run this project successfuly.
To install the, the golang toolchain go [here](https://golang.org/doc/install). 

# Run Application

```bash
./myhttp [-parallel 3] example.com example2.com
```

- The `parallel` flag sets the number of processes to run in parallel. It defaults to `10` if the flag is not provided.
- The default value for maximum parallel processes is `30`. 
You can reset this value by setting `MAX_PROC` in your environmental variables to a suitable value

## Run Tests
```
go test
```
