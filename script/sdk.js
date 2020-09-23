var $ = {
    util: {},
    HttpPromise: {
        create: function (response) {
            var promise = {
                status: $.resp.status(response) == 200,
                then: function (func) {
                    if (this.status) {
                        func($.resp.data(response), $.resp.header(response));
                    }
                    return this
                },
                catch: function (func) {
                    if (!this.status) {
                        func($.resp.status(response), $.resp.data(response), $.resp.header(response));
                    }
                }
            }
            return promise
        }
    },
    DBPromise: {
        create: function (result) {
            var promise = {
                status: result["StatusCode"] == 200,
                then: function (func) {
                    if (this.status) {
                        func(result["Data"]);
                    }
                    return this
                },
                catch: function (func) {
                    if (!this.status) {
                        func(result["Message"]);
                    }
                }
            }
            return promise
        }
    },
    DB: {
        open: function (name) {
            var db = {
                query: function (sql, arg, pageNumber, pageSize) {
                    var result = _db_query(name, sql, arg, pageNumber, pageSize);
                    return $.DBPromise.create(result);
                },
                queryOne: function (sql, arg) {
                    var result = _db_queryOne(name, sql, arg);
                    return $.DBPromise.create(result);
                },
                save: function (table, arg) {
                    var result = _db_save(name, table, arg);
                    return $.DBPromise.create(result);
                },
                exec: function (sql, arg) {
                    var result = _db_exec(name, sql, arg);
                    return $.DBPromise.create(result);
                },
            }
            return db
        }
    },
    log: function (msg) {
        _sys_log(msg);
    },
    err: function (msg) {
        _sys_err(msg);
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
            try {
                return JSON.parse(body);
            } catch (e) {
                return {}
            }
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
        var promise = this.HttpPromise.create(resp);
        return promise;
    },
    post: function (url, param, header) {
        var resp = _http_post(url, param, header);
        var promise = this.HttpPromise.create(resp);
        return promise;
    },
    delete: function (url, param, header) {
        var resp = _http_delete(url, param, header);
        var promise = this.HttpPromise.create(resp);
        return promise;
    },
    put: function (url, param, header) {
        var resp = _http_put(url, param, header);
        var promise = this.HttpPromise.create(resp);
        return promise;
    },
    define: function (definition) {
        __serv_definition = definition;
    },
    date: function () {
        var date = new Date()
        return date.toLocaleDateString()
    },
    time: function () {
        var date = new Date()
        return date.toLocaleDateString() + " " + date.toLocaleTimeString()
    }
};