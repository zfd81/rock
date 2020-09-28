var $ = {
    KVS: {
        use: function (name) {
            var kvs = {
                get: function (key) {
                    var result = _kv_get(name, key)
                    if (result.Normal) {
                        return result.Data
                    } else {
                        throw new Error(result.Message);
                    }
                },
                set: function (key, value, ttl) {
                    var result = _kv_set(name, key, value, ttl)
                    if (!result.Normal) {
                        throw new Error(result.Message);
                    }
                }
            }
            return kvs;
        }
    },
    DB: {
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
        open: function (name) {
            var db = {
                query: function (sql, arg, pageNumber, pageSize) {
                    var result = _db_query(name, sql, arg, pageNumber, pageSize);
                    return $.DB.DBPromise.create(result);
                },
                queryOne: function (sql, arg) {
                    var result = _db_queryOne(name, sql, arg);
                    return $.DB.DBPromise.create(result);
                },
                save: function (table, arg) {
                    var result = _db_save(name, table, arg);
                    return $.DB.DBPromise.create(result);
                },
                exec: function (sql, arg) {
                    var result = _db_exec(name, sql, arg);
                    return $.DB.DBPromise.create(result);
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
    http: {
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
        }
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
    },
    left: function (str, length) {
        if (str == null)
            return null;
        if (length < 0) {
            return "";
        }
        if (length < str.length) {
            return str.substring(0, length);
        } else {
            return str;
        }
    },
    right: function (str, length) {
        if (str == null)
            return null;
        if (length < 0) {
            return "";
        }
        if (length < str.length) {
            return str.substring(str.length - length);
        } else {
            return str;
        }
    },
    startsWith: function (str, prefix) {
        if (prefix == null || prefix == "") {
            return true;
        }
        return str.indexOf(prefix) == 0
    },
    endsWith: function (str, suffix) {
        if (suffix == null || suffix == "") {
            return true;
        }
        var reg = new RegExp(suffix + "$");
        return reg.test(str);
    }
};