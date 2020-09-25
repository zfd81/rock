module.exports = {
    path: "/example/utils",
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
    startsWith: function (str, prefix) {
        if (prefix == null || prefix == "") {
            return true;
        }
        return str.indexOf(prefix) == 0
    }
};