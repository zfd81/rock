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
        Result: {
            create: function () {
                var result = {
                    status: false,
                    data: undefined,
                    msg: "",
                    setVal: function (val) {
                        this.data = val;
                        this.status = true;
                    },
                    setMsg: function (msg) {
                        this.msg = msg;
                    }
                }
                return result
            }
        },
        DBPromise: {
            create: function (result) {
                var promise = {
                    then: function (func) {
                        if (result.status) {
                            func(result.data);
                        }
                        return this
                    },
                    catch: function (func) {
                        if (!result.status) {
                            func(result.msg);
                        }
                    }
                }
                return promise
            }
        },
        open: function (name) {
            var db = {
                query: function (sql, arg, pageNumber, pageSize) {
                    var result = $.DB.Result.create();
                    try {
                        if (sql == undefined || sql == null || sql == "") {
                            result.setMsg("SQL statement cannot be empty")
                            console.log("SQL statement cannot be empty")
                        } else {
                            if (arg == undefined) {
                                arg = null;
                            }
                            if (pageNumber == undefined) {
                                pageNumber = 0;
                            } else {
                                if (typeof pageNumber != "number") {
                                    pageNumber = parseInt(pageNumber)
                                }
                            }
                            if (pageSize == undefined) {
                                pageSize = 10;
                            } else {
                                if (typeof pageSize != "number") {
                                    pageSize = parseInt(pageSize)
                                }
                            }
                            result.setVal(_db_query(name, sql, arg, pageNumber, pageSize));
                        }
                    } catch (err) {
                        result.setMsg(err);
                        console.log(err);
                    }
                    return $.DB.DBPromise.create(result);
                },
                queryOne: function (sql, arg) {
                    var result = $.DB.Result.create();
                    try {
                        if (sql == undefined || sql == null || sql == "") {
                            result.setMsg("SQL statement cannot be empty")
                            console.log("SQL statement cannot be empty")
                        } else {
                            if (arg == undefined) {
                                arg = null;
                            }
                            result.setVal(_db_queryOne(name, sql, arg));
                        }
                    } catch (err) {
                        result.setMsg(err);
                        console.log(err);
                    }
                    return $.DB.DBPromise.create(result);
                },
                save: function (table, arg) {
                    var result = $.DB.Result.create();
                    try {
                        if (table == undefined || table == null || table == "") {
                            result.setMsg("Table name cannot be empty")
                            console.log("Table name cannot be empty")
                        } else if (arg == undefined || typeof arg != "object") {
                            result.setMsg("Parameter data type error")
                            console.log("Parameter data type error")
                        } else {
                            result.setVal(_db_save(name, table, arg));
                        }
                    } catch (err) {
                        result.setMsg(err);
                        console.log(err);
                    }
                    return $.DB.DBPromise.create(result);
                },
                exec: function (sql, arg) {
                    var result = $.DB.Result.create();
                    try {
                        if (sql == undefined || sql == null || sql == "") {
                            result.setMsg("SQL statement cannot be empty")
                            console.log("SQL statement cannot be empty")
                        } else {
                            if (arg == undefined) {
                                arg = {};
                            }
                            result.setVal(_db_exec(name, sql, arg));
                        }
                    } catch (err) {
                        result.setMsg(err);
                        console.log(err);
                    }
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