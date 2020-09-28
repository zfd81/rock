$.define({
    path: "/kvs/get",
    method: "get"
})

var kvs = $.KVS.use("test")
$.log(kvs.get("name"))
$.log(kvs.get("age"))
$.log(kvs.get("boy"))
var user = kvs.get("user")
if (user != null) {
    $.log(user.name)
}
$.log(kvs.get("friend"))
$.log(kvs.get("tels"))
var users = kvs.get("users")
if (users != null) {
    $.log(users[1].name)
}