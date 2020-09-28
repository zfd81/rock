$.define({
    path: "/kvs/set",
    method: "get"
})

var kvs = $.KVS.use("test")
kvs.set("name", "zhang", 15)
kvs.set("age", 32, 15)
kvs.set("boy", true)
kvs.set("user", {
    name: "zhang",
    age: 32,
    boy: true
})
kvs.set("friend", ["zhang", "wang", "li"], 70)
kvs.set("tels", [12, 22, 32], 70)
kvs.set("users", [{
    name: "zhang",
    age: 32,
    boy: true
}, {
    name: "wang",
    age: 30,
    boy: false
}])