$.define({
    path: "/user/jwt",
    method: "get",
})

var data = {
        code: "d11",
        name: "dzfdde",
        status: "daaa"
    }

var token = $.jwt.create(data, "zfd")
$.log(token)
var sleep = function (time) {
    var startTime = new Date().getTime() + parseInt(time, 10);
    while (new Date().getTime() < startTime) {
    }
};

sleep(7000); // 延时函数，单位ms

try {
    var d = $.jwt.parse(token, "zfd")
    $.log(d.name)
} catch (e) {
    $.err(e)
}
