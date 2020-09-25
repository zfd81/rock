$.define({
    path: "/user",
    method: "post",
    params: [{
        name: "user",
        dataType: "map"
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

var table = "das_sys_user" //表名

//通过save方法添加用户，要求user的key必须和das_sys_user的列名相同
db.save(table, user)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });

//删除用户
db.exec("delete from das_sys_user where id=:val", user.id)

var sql = " \
    insert into das_sys_user \
        (id,name,password,full_name,phone_number,email,status,creator,created_time,modifier,modified_time) \
    values \
        (:id,:name,:password,:full_name,:phone_number,:email,:status,:creator,:created_time,:modifier,:modified_time)"

//通过exec方法添加用户，要求需要自己写SQL语句，user的key必须和SQL语句中“:变量”名相同
db.exec(sql, user)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });