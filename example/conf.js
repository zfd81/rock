var utils = require("/example/utils") //获得工具模块
module.exports = {
    path: "/example/conf",
    ds: "datasource",
    serv: function (url) {
        if (utils.startsWith(url, "/")) {
            return "http://localhost:8081" + url;
        }
        return "http://localhost:8081/" + url;
    }
};