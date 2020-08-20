
function _log(msg){
    console.log(msg);
}

function _status(resp){
    return resp["StatusCode"];
}

function __body(resp){
    return resp["Content"];
}

function _body(resp){
    return JSON.parse(resp["Content"]);
}

function _header(resp){
    return resp["Header"];
}