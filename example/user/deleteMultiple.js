$.define({
    path: "/user/multiple",
    method: "delete",
    params: [{
        name: "ids",
        dataType: "string[]"
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

//删除多个用户
db.exec("delete from das_sys_user where {@vals[OR] id=:this.val}", ids)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });