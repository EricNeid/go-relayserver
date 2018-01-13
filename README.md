# go - openweatherapi

This Repo contains golang library to query OpenWetherMaps (http://openweathermap.org/) for weather information.

* current weather: http://openweathermap.org/current
* 5 days forecast: http://openweathermap.org/forecast5

## Install


```sh
 $ go get github.com/CombatMage/openweatherapi
```

## Documentation

Is available on ``godoc``:

 https://godoc.org/github.com/CombatMage/openweatherapi

## Examples

Consuming the library:
```go
// create a query
q := NewQueryForCity(readAPIKey(), "Berlin,de")

// obtain data
resp, err := q.Weather()

// enjoy
fmt.Println(resp.Name) // Berlin
fmt.Println(resp.Weather[0].Description) // mis
fmt.Println(resp.Main.Temp) // 1
```

See the test files for more example.

A simple client for testing is also included:
```sh
 $ go build cmd/openweatherclient
 $ openweatherclient -key <OpenWeather API Key> -city Berlin,de
```
