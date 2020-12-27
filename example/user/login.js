$.define({
    path: "user/login",
    method: "post",
    params: [{
        name: "name",
        dataType: "string"
    }, {
        name: "pwd",
        dataType: "string"
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

//查询单用户信息
db.queryOne("select * from das_sys_user where name=:name and password=:pwd", {
        name: name,
        pwd: pwd
    })
    .then(function (data) {
        if (data == null) {
            $.resp.write({msg: "用户名或密码错误"}, {code: 400})
            return
        }
        var token = new Date().toString()
        $.resp.write(data, {
            code: 200,
            token: token,
        })
    })
    .catch(function (error) {
        $.err(error)
    });