
package beans


type UserContext struct {
	UserId			string // 用户id
	UpdateTime     int64    // 记录更新时间
	Plateform    	string // 登陆平台
	Role         	int    // 角色
	Ip       	    int    // 登陆ip地址
	WeipanToken		string // 微盘token
	WeipanUserId	int    // 微盘用户 恒大用户 id
	WeipanTokenCreateTime	int64    // 微盘token创建时间
	ManageRoomIds	string    // 管理的房间列表
	MemberId		int64    // 会员id
}