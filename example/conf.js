module.exports = {
    path: "/example/conf",
    ds: "datasource",
    serv: function (url) {
        if ($.startsWith(url, "/")) {
            return "http://localhost:8081" + url;
        }
        return "http://localhost:8081/" + url;
    }
};