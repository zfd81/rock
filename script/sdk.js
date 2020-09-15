var $ = {
    HttpPromise: {
        create: function () {
            var promise = {
                response: {},
                then: function (func) {
                    if ($.resp.status(this.response) == 200) {
                        func($.resp.data(this.response), $.resp.header(this.response));
                    }
                    return this
                },
                catch: function (func) {
                    if ($.resp.status(this.response) != 200) {
                        func($.resp.status(this.response), $.resp.data(this.response), $.resp.header(this.response));
                    }
                }
            }
            return promise
        }
    },
    log: function (msg) {
        _sys_log(msg);
    },
    resp: {
        status: function (resp) {
            return resp["StatusCode"];
        },
        content: function (resp) {
            return resp["Content"];
        },
        data: function (resp) {
            var body = this.content(resp)
            if (typeof body == "undefined" || body == null || body == "") {
                return {}
            }
            return JSON.parse(body);
        },
        header: function (resp) {
            return resp["Header"];
        },
        write: function (data, header) {
            _resp_write(data, header)
        }
    },
    get: function (url, param, header) {
        var resp = _http_get(url, param, header);
        var promise = this.HttpPromise.create();
        promise.response = resp;
        return promise;
    },
    post: function (url, param, header) {
        var resp = _http_post(url, param, header);
        var promise = this.HttpPromise.create();
        promise.response = resp;
        return promise;
    },
    delete: function (url, param, header) {
        var resp = _http_delete(url, param, header);
        var promise = this.HttpPromise.create();
        promise.response = resp;
        return promise;
    },
    put: function (url, param, header) {
        var resp = _http_put(url, param, header);
        var promise = this.HttpPromise.create();
        promise.response = resp;
        return promise;
    },
    define: function (definition) {
        __serv_definition = definition;
    }
};