$.define({
    path: "user/findAll",
    method: "get"
})

var db = $.DB.open("datasource");

$.log("遍历对象方法一：")
db.queryOne("select * from das_sys_user where id=:val", "202004171900")
    .then(function (data) {
        for (var key in data) {
            $.log(key + ":" + data[key])
        }
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("遍历对象方法二：")
db.queryOne("select * from das_sys_user where id=:val", "202004171900")
    .then(function (data) {
        Object.keys(data).forEach(function (key) {
            $.log(key + ":" + data[key])
        })
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("获得日期数据类型：")
var sql = "select \
                DATE_FORMAT(created_time, '%Y-%m-%d %H\\:%i') AS created_time \
            from \
                das_sys_user \
            where \
                id=:val"
db.queryOne(sql, "202004171900")
    .then(function (data) {
        var created_time = data["created_time"]
        $.log(created_time)
        $.err(created_time)
    })
    .catch(function (error) {
        $.log(error)
    });