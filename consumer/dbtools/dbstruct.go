package dbtools

type User struct {
	Id int
	Money int
	Cnt int
}

type Record struct {
	Id int
	Uid int
	Val int
	Stime int64
	Opened int
}

type Env struct {
	Id int
	Opened int
	Val int
	Stime int64
}
