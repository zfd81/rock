$.define({
    path: "user/findAll",
    method: "get"
})

var query = $.DB.open("datasource");

$.log("测试SQL查询：")
query.query("select * from das_sys_user")
    .then(function (data) {
        //循环方式一
        for (i = 0, len = data.length; i < len; i++) {
            $.log(data[i].name)
        }
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("测试单一参数的SQL查询：")
query.query("select * from das_sys_user where status=:val","1")
    .then(function (data) {
        //循环方式二
        data.forEach(function (item,index,array) {
            $.log(item.name)
        })
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("测试多参数的SQL查询：")
query.query("select * from das_sys_user where status=:sts and creator=:ctr",{
    sts:"1",
    ctr:"202004171900"
})
    .then(function (data) {
        //循环方式一
        for (i = 0, len = data.length; i < len; i++) {
            $.log(data[i].name)
        }
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("测试分页SQL查询：")
query.query("select * from das_sys_user where status=:val","1",1,2)
    .then(function (data) {
        //循环方式一
        for (i = 0, len = data.length; i < len; i++) {
            $.log(data[i].name)
        }
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("测试like参数的SQL查询(方法一)：")
query.query("select * from das_sys_user where name like :val","%d")
    .then(function (data) {
        //循环方式一
        for (i = 0, len = data.length; i < len; i++) {
            $.log(data[i].name)
        }
    })
    .catch(function (error) {
        $.log(error)
    });

$.log("测试like参数的SQL查询(方法二)：")
query.query("select * from das_sys_user where name like CONCAT('%',:val)","%dm%")
    .then(function (data) {
        //循环方式一
        for (i = 0, len = data.length; i < len; i++) {
            $.log(data[i].name)
        }
    })
    .catch(function (error) {
        $.log(error)
    });