$.define({
    path: "/httpclient/serv",
    method: "get"
})

var conf = require("/example/conf") //获得配置模块
var utils = require("/example/utils")

$.log(utils.left("abcde",2))
var param = {
    name: "zfd",
    pwd: "123456"
}
$.http.post(conf.serv("/user/login"), param)
    .then(function (data, header) {
        if (header.Code == 200) {
            $.log("===================");
            $.resp.write(data.full_name + "登陆成功")
            return
        }
        $.log("===================");
        $.resp.write("登陆失败，错误码为："+header.Code)
    })
    .catch(function (status, data, header) {
        $.resp.write(data)
    });