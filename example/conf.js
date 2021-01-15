var u = require("/example/utils") //获得配置模块

exports = {
    define:{
        path: "/example/conf",
    },
    dd: u.left("hello",2),
    ds: "datasource",
    serv: function (url) {
        if ($.startsWith(url, "/")) {
            return "http://localhost:8081" + url;
        }
        return "http://localhost:8081/" + url;
    }
};