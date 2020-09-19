$.define({
    path: "user/findAll",
    method: "get"
})

var db = $.DB.open("datasource");

var table = "das_sys_user" //表名

//记录信息
var info = {
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
$.log($.date()) //输出当前日期
$.log($.time()) //输出当前时间
db.save(table, info)
    .then(function (data) {
        $.log(data)
        $.log("保存成功")
    })
    .catch(function (error) {
        $.err(error)
    });