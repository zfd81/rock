$.define({
    path: "/user/multiple",
    method: "post",
    params: [{
        name: "users",
        dataType: "map[]"
    }]
})

var conf = require("/example/conf") //获得配置模块
var db = $.DB.open(conf.ds); //获取数据操作对象

var table = "das_sys_user" //表名

//通过save方法添加用户，要求user的key必须和das_sys_user的列名相同
db.save(table, users)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });

var sql = " \
    insert into das_sys_user \
        (id,name,password,full_name,phone_number,email,status,creator,created_time,modifier,modified_time) \
    values \
        {@vals[,] (:this.id,:this.name,:this.password,:this.full_name,:this.phone_number,:this.email,:this.status,:this.creator,:this.created_time,:this.modifier,:this.modified_time)}"
        
//删除用户
db.exec("delete from das_sys_user where {@vals[OR] id=:this.id}", users)

//通过exec方法添加用户，要求需要自己写SQL语句，user的key必须和SQL语句中“:变量”名相同
db.exec(sql, users)
    .then(function (data) {
        $.resp.write(data)
    })
    .catch(function (error) {
        $.err(error)
    });