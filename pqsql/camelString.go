// package pqsql

// import (
// 	"strings"
// )

// func StringToCamel(str string) string{
// 	temp := strings.Split(str, "-")

// 	var temp=str.split("-");
// 	for(var i=1;i<temp.length;i++){
// 		 temp[i]=temp[i][0].toUpperCase()+temp[i].slice(1);
// 	 }
// 	  return temp.join("");
//   }

//   func StringToCamelCase(str string) string{
// 	var reg=/-(\w)/g;//子项()表示子项

// 	return str.replace(reg,function($0,$1){//$0代表正则整体，replace（）方法中的第二个参数若是回调函数，那么这个回调函数中的参数就是匹配成功后的结果
// 		//若回调函数中有多个参数时，第一个参数代表整个正则匹配结果，第二个参数代表第一个子项
// 		alert($0);//-b
// 		alert($1);//b
// 		return $1.toUpperCase();
// 	});
// }