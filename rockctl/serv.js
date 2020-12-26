
$.define({
    path: "zfd/{name}/{password}/1",
    method: "get",
    params: [{
        name: "user",
        dataType: "map"
    }]
})

var m = {
    name:name,
    pwd:password
}
$.log(123)
$.log("hello")
$.log(m)

                var obj = {name:"zfd",pwd:"password"}
                $.resp.write(m,user)
