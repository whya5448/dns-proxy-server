### Caches
* [List Caches](#list-caches)
* [Get Cache Actual Size](#get-cache-actual-size)

### List caches

	GET /v1/caches HTTP/1.1

#### Response

	HTTP/1.1 200 OK
	Content-Type: application/json

```json
{
  "acme.com": {
    "CreationDate": "2018-05-06T19:31:40.62282457Z",
    "TimeoutDuration": 86400000000000, // nano seconds
    "Val": null
  },
	....
}
```


### Get Cache Actual Size

	GET /v1/caches/size HTTP/1.1

#### Response

	HTTP/1.1 200 OK
	Content-Type: application/json

```json
{
	"size": 1
}
```
