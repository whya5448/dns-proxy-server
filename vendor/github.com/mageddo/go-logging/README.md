A really simple and fast logger implementation independent

Compability 

* Compatible with go 1.9
* Works with go 1.7 but line numbers will be not accurate

Example

```go
import "github.com/mageddo/go-logging"
...
logging.Debugf("hey %s", "Mark")
logging.Infof("hey %s", "Mark")
logging.Warnf("hey %s", "Mark")
logging.Errorf("hey %s", "Mark")
```

Out
```
2018/05/05 22:36:27 DEBUG f=main.go:8 pkg=main m=main hey Mark
2018/05/05 22:36:27 INFO f=main.go:9 pkg=main m=main hey Mark
2018/05/05 22:36:27 WARNING f=main.go:10 pkg=main m=main hey Mark
2018/05/05 22:36:27 ERROR f=main.go:12 pkg=main m=main hey Mark
```

Testing it

```
docker-compose up --abort-on-container-exit ci-build-test
```