package main

// 首先得完成系统设计
// 考虑树洞的验证
// 每个树洞的匿名性怎么办???
// 这有一个设计原则在于同一个树洞下的帖子得是同一个Id, THU Hole是怎么维护的呢????
// 然后一个树洞需要维护的信息又哪些
// 考虑先设计几个接口

// 如何保证同一个人在同一个树洞里面id是一致的呢????
//
//
// 可以考虑利用hash pid + user ====> 生成一个值,
func SaveUser(username, password string) error {

}

func FindUser(username, password string) (int, string, error) {

}
