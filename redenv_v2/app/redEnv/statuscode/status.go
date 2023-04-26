package statuscode

var OK int = 0      //OK
var NoSuchUser = -1 //没有这名用户
var TooManyEnv = -2 //获取红包超过上限
var Thankyou = -3	//没抢到 只能说运气并不是很好
var NoThisEnv = -4  //没有这条红包记录
var AlreadyOpened = -5	//该红包已经被打开过了
var FlowOverrun = -6
