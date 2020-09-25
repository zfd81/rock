$.define({
    path: "/user/findAll/num/{num}/size/{size}",
    method: "get"
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

//分页查询用户信息
db.query("select * from das_sys_user", null, num, size)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });