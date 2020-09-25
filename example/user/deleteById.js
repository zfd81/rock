$.define({
    path: "/user/id/{id}",
    method: "delete"
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象


//删除用户信息
db.exec("delete from das_sys_user where id=:val", id)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });