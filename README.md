## FCFS & RR Servers

This lab implements two simple M/G/k server farms with two queuing disciplines - FCFS and RR. The M/G/k system consists of a single queue, served by k servers, where jobs arrive according to a Poisson process with rate
λ and have generally distributed i.i.d. job service requirements (normally in this case) with mean 1/μ.

### Usage
Most of the parameters, such as arrival/service rates, HTTP post or number of servers in the farm k, can be modified with flags. For this it is required to compile the binary and launch it with flags. Example usage:
```
$ go build post.go
$ ./post -n 3 -u 4
```

If parameters are not specified, binary will launch with defaults. General usage and all available flags are available via
```
$ ./main --help
$ ./post --help
```
