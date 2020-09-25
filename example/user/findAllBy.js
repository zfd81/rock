$.define({
    path: "/user/findAll",
    method: "post",
    params: [{
        name: "status",
        dataType: "string",
    }, {
        name: "info",
        dataType: "map",
    }, {
        name: "pageNumber",
        dataType: "int",
    }, {
        name: "pageSize",
        dataType: "int",
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

//分页查询用户信息
db.query("select * from das_sys_user where status=:status {and name like CONCAT(:name,'%')}", {
        status: status,
        name: info.name
    }, pageNumber, pageSize)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });