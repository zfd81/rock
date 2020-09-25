$.define({
    path: "/user",
    method: "put",
    params: [{
        name: "user",
        dataType: "map"
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象


//修改用户信息
db.exec("update das_sys_user set name=:name, password=:password, full_name=:full_name where id=:id", user)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });