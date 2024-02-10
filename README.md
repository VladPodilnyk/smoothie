# Smoothie

A dead simple distributed rate limiter that smoothes your day.

### Yet another rate limiter?
Smoothie is meant to be a distributed rate limiter, unlike all other popular within go community
solutions. The project is quite small so, it doesn't make
any sense to publish it as a package. But it can be a good showcase of how to
make a distributed rate limiter using redis and Go.
If you don't need that or maybe you don't like Redis, 
then standard `time` package or `rate` package might serve you better.

### Implementation details
At the moment smoothie implements a fixed-size window strategy for limiting requests.
The api resembles the familiar `rate` package API with `Allow` and `Exec` methods.
See example below
```go
config := redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}
limiter := New(&testRedisOptions, Rate{NumberOfRequests: 1, Duration: 5 * time.Second})
isAllowed := limiter.Allow(ctx.Background(), "user-IP-address")
testEffect := func () error {
    fmt.Println("hello")
    return nil
}
err := limiter.Exec(ctx.Background(), "user-IP-adress", testEffect)
```

### Future plans
It would be nice to add the following things to the smoothie:
- other rate limiting strategies (Sliding window, token bucket)
- extend tests with test cases with concurrent access to the same limiting `key`
