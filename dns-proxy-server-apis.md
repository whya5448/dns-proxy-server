```javascript
// Get ENVs
fetch('/env/', {
  headers: {
    'Accept': 'application/json'
  }
})

/* Returns
 * [
 *   {
 *     "name": "",
 *     "hostnames": [
 *       {
 *         "id": "1560631033461458993",
 *         "hostname": "docker.dns",
 *         "ip": [
 *           0,
 *           0,
 *           0,
 *           0
 *         ],
 *         "target": "localhost",
 *         "ttl": 60,
 *         "type": "CNAME"
 *       }
 *     ]
 *   },
 *   {
 *     "name": "some-env"
 *   }
 * ] */

// Create ENV
fetch('/env/', {
  method: 'POST',
  body: JSON.stringify({
    name: "some-env.2"
  }),
  headers: {
    "Content-Type": "application/json"
  }
})
```
