# rock


## Example

Performing a `GET` request

```js
_.get("http://localhost:8080/auth/das/user/{name}", {
        name: "zfd" //路径参数name
    }, {
        token: "1234567890" //添加请求头信息
    })
    .then(function (data) {
        _.log(data["full_name"]) //日志打印响应中的full_name
    })
    .catch(function (status) {
        _.log(status) //日志打印响应状态码(status)
    });
```

Performing a `POST` request

```js
_.post("http://localhost:8080/login", {
        name: "zfd", //参数name
        password: "123456" //参数password
    })
    .then(function (data, header) {
        _.log(data["full_name"]) //日志打印响应中的full_name
        _.log(header["token"]) //日志打印响应头信息中的token
    })
    .catch(function (status, data, header) {
        _.log(status) //日志打印响应状态码(status)
        _.log(data.err)
        _.log(header.token)
    });
```

Performing nested requests

```js
_.post("http://localhost:8080/login", {
        name: "zfd",
        password: "123456"
    })
    .then(function (data, header) {
        _.get("http://localhost:8080/auth/das/user/{name}", {
                name: data.name
            }, {
                token: header["token"]
            })
            .then(function (data) {
                _.log(data["name"])
            });
    })
    .catch(function (status) {
        _.log(status)
    });
```

## rock API

##### get(url[, param, header])

```js
// Send a GET request
_.get("http://localhost:8080/user/id/20200904");
```


##### post(url[, param, header])

```js
// Send a POST request
_.post("http://localhost:8080/login", {
        name: "zfd",
        password: "123456"
    })
```


##### delete(url[, param, header])

```js
// Send a DELETE request
_.delete("http://localhost:8080/user/name/{name}", {
        name: "zfd"
    })
```


##### put(url[, param, header])

```js
// Send a PUT request
_.put("http://localhost:8080/user/id/{id}", {
        id: "20200904",
        name: "dada",
        password: "88888888"
    })
```

