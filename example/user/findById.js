$.define({
    path: "/user/id/{id}",
    method: "get"
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

//定义查询SQL
var sql = "SELECT \
                id, \
                name, \
                full_name, \
                phone_number, \
                email, \
                creator, \
                DATE_FORMAT(created_time, '%Y-%m-%d %H\\:%i\\:%s') AS created_time, \
                modifier, \
                DATE_FORMAT(modified_time, '%Y-%m-%d %H\\:%i\\:%s') AS modified_time \
            FROM \
                das_sys_user \
            WHERE \
                id = :val";
//查询单用户信息
db.queryOne(sql, id)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });