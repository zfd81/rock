var _ = {
    Promise: {
        create: function () {
            var promise = {
                response: {},
                then: function (func) {
                    if (_.status(this.response) == 200) {
                        func(_.content(this.response), _.header(this.response));
                    }
                    return this
                },
                catch: function (func) {
                    if (_.status(this.response) != 200) {
                        func(_.status(this.response), _.content(this.response), _.header(this.response));
                    }
                }
            }
            return promise
        }
    },
    log: function (msg) {
        console.log(msg);
    },
    status: function (resp) {
        return resp["StatusCode"];
    },
    body: function (resp) {
        return resp["Content"];
    },
    content: function (resp) {
        var body = this.body(resp)
        if (typeof body == "undefined" || body == null || body == "") {
            return {}
        }
        return JSON.parse(body);
    },
    header: function (resp) {
        return resp["Header"];
    },
    get: function (url, param, header) {
        var resp = _http_get(url, param, header);
        var promise = this.Promise.create();
        promise.response = resp;
        return promise;
    },
    post: function (url, param, header) {
        var resp = _http_post(url, param, header);
        var promise = this.Promise.create();
        promise.response = resp;
        return promise;
    }
};