// 后端接口的定义
export const BaseUrl = 'http://172.16.5.28:8081'

//后端接口的路径
export  const ServerPath:Record<string,string> = {
  'login': '/api/v1/login', // 登录接口
  'eventList': '/event/list', // 活动列表接口
  'register': '/user/reg', // 注册接口
  'isJoin': '/event/isJoin', // 检查是否加入活动接口
  'teamList': '/eventMember/teamList', // 获取活动团队列表接口
  'getComment': '/eventComment/listByEventId', // 获取活动评论接口
  'addComment': '/eventComment/insert', // 添加活动评论接口
  'join': '/eventMember/join', // 加入活动接口
  'eventMemberList': '/eventMember/list', // 获取活动成员列表接口
  'sendApply': '/friend/sendApply', // 发送好友申请接口
  'applyList': '/friend/applyList', // 获取好友申请列表接口
  'friendList': '/friend/getFriendList', // 获取好友列表接口
  'processApply': '/friend/processApply', // 处理好友申请接口
  'friendStatus': '/friend/status', // 获取好友状态接口
  'userInfo': '/user/info', // 获取用户信息接口
  'updatePassword': '/user/updatePassword', // 修改密码接口
  'RecommendEvent': '/event/recommend', // 推荐活动接口
  'UpdateUserInfo': '/user/update', // 更新用户信息接口（不修改头像）
  'searchList': '/event/searchEventInfoList', // 搜索活动接口
  'addEvent': '/event/insert', // 添加活动接口
  'editEvent': '/event/edit', // 编辑活动接口
  'myEventList': '/event/myEventList', // 获取我的活动接口
  'deleteEvent': '/event/delete', // 删除活动接口
  'uploadImg': '/upload/uploadFile', // 文件上传接口
  'eventInfo': '/event/info', // 获取活动详情接口
  'chatList': '/chat/byUserIdFromAndTo', // 获取聊天记录接口
  'chatAdd': '/chat/insert', // 新增聊天记录接口
  'idEventList':'/v1/activity/{id}' //根据用户id查找活动
};

// app存储的数据库名称
export const PREFERENCES_NAME = 'hobby_db';

// app存储的数据库，
export const PREFERENCES_Name = 'hobby.db'

//获取后端接口，成功响应码
export const HTTP_CODE_200:number=200

//超时时间设置
export const HTTP_READ_TIMEOUT:number = 100000