$.define({
    path: "user/findAll",
    method: "get"
})

var db = $.DB.open("datasource");


$.log("删除记录")
db.exec("delete from das_sys_user where id=:val", 25)
    .then(function (data) {
        $.log(data)
        $.log("删除成功")
    })
    .catch(function (error) {
        $.err(error)
    });


$.log("增加记录")
var sql = "insert into \
            das_sys_user \
            (id,name,password,full_name,phone_number,email,status,creator,created_time,modifier,modified_time) \
            values \
            (:Id,:Name,:Password,:full_name,:phone_number,:email,:status,:Creator,:created_time,:Modifier,:modified_time)"

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
db.exec(sql, info)
    .then(function (data) {
        $.log(data)
        $.log("保存成功")
    })
    .catch(function (error) {
        $.err(error)
    });

$.log("修改记录")
var newInfo = info
newInfo.Name = "user25_new"
newInfo.Password = "pwd25_new"
newInfo.full_name = "用户25_new"
db.exec("update das_sys_user set name=:Name, password=:Password, full_name=:full_name where id=:Id", newInfo)
    .then(function (data) {
        $.log(data)
        $.log("修改成功")
    })
    .catch(function (error) {
        $.err(error)
    });