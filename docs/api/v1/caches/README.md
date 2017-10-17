### Caches
* [List Caches](#list-caches)
* [Get Cache Actual Size](#Get-Cache-Actual-Size)

### List caches

	GET /v1/caches HTTP/1.1

#### Response

	HTTP/1.1 200 OK
	Content-Type: application/json

```json
{
	"key1": "value1"
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
