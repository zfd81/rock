module.exports = {
    path: "/comm/util",
    message: 'hello world',
    left: function (str, length) {
        if (str == null)
            return null;
        if (length < 0) {
            return "";
        }
        if (length < str.length) {
            return str.substring(0, length);
        } else {
            return str;
        }
    },
    right: function (str, length) {
        if (str == null)
            return null;
        if (length < 0) {
            return "";
        }
        if (length < str.length) {
            return str.substring(str.length - length);
        } else {
            return str;
        }
    },
    print: function (msg) {
        $.log(msg);
    },
    add: function (a, b) {
        return a + b;
    }
};