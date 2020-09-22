
$.define({
    path: "zfd/{name}/{password}/",
    method: "get"
})

var util = require("comm/util")
$.log(util.add(2,3))
$.log(util.left("helloworld",3))
$.log(util.right("helloworld",3))
util.print("hello world!!!")
$.post("http://localhost:8080/login", {
    name: name,
    password: password
}, {})
    .then(function (data, header) {
        $.get("http://localhost:8080/auth/das/test/{name}", {
            name: data.name
        }, {
            zxcvb: header["Atv"]
        })
            .then(function (data) {
                $.log("aaa")
                var obj = {name:"zfd",pwd:"password"}
                $.resp.write(data,obj)
            });
    })
    .catch(function (status, data) {
        $.log(status)
        $.log(JSON.stringify(data))
    });