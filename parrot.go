package parrot

import (
	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()
	vm.Run(`
var  str = '{"a":"a1","b":"b1"}';
  var obj1 = eval('('+xyzzy+')'); //使用eval函数

   console.log(obj1.msg); // 4
`)
}
