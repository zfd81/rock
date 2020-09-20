$.define({
    path: "user/findAll",
    method: "get"
})

var db = $.DB.open("datasource");

var table = "das_sys_user" //表名

//一条用户信息
var user = {
    Id: 25,
    Name: "user25",
    Password: "pwd25",
    full_name: "用户25",
    phone_number: "188888888",
    email: "25@qq.com",
    status: 1,
    Creator: 1,
    created_time: $.time(), //当前时间
    Modifier: 1,
    modified_time: $.time() //当前时间
}

db.exec("delete from das_sys_user where id=:val", 25)
    .then(function (data) {
        $.log(data)
        $.log("删除成功")
    })
    .catch(function (error) {
        $.err(error)
    });
db.save(table, user) //添加一条记录
    .then(function (data) {
        $.log(data)
        $.log("保存成功")
    })
    .catch(function (error) {
        $.err(error)
    });

db.exec("delete from das_sys_user where id=26 or id=27")
    .then(function (data) {
        $.log(data)
        $.log("删除成功")
    })
    .catch(function (error) {
        $.err(error)
    });

//多条用户信息
var users = [{
    Id: 26,
    Name: "user26",
    Password: "pwd26",
    full_name: "用户26",
    phone_number: "188888888",
    email: "26@qq.com",
    status: 1,
    Creator: 1,
    created_time: $.time(), //当前时间
    Modifier: 1,
    modified_time: $.time() //当前时间
}, {
    Id: 27,
    Name: "user27",
    Password: "pwd27",
    full_name: "用户27",
    phone_number: "188888888",
    email: "27@qq.com",
    status: 1,
    Creator: 1,
    created_time: $.time(), //当前时间
    Modifier: 1,
    modified_time: $.time() //当前时间
}]
db.save(table, users) //添加多条记录
    .then(function (data) {
        $.log(data)
        $.log("保存成功")
    })
    .catch(function (error) {
        $.err(error)
    });