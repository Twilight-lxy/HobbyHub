//后端接口定义
export const BaseUrl = 'http://192.168.1.9:8081' //Zhuanyong5G

//后端接口的路径
export const ServerPath: Record<string, string> = {
  'login': '/api/v1/login', //登录接口
  'register': '/api/v1/user',//注册接口
  'join':'/api/v1/activity/{id}/member',//加入活动接口
  'addEvent':'/api/v1/activity',//新增活动
  'myEventList':'/api/v1/activity', //获取我创建的活动
  'eventInfo':'/api/v1/activity/{id}',//获取活动详情
  'editEvent':'/api/v1/activity/{id}',//修改活动信息
  'deleteEvent':'/api/v1/activity/{id}',//删除活动
  'eventList':'/api/v1/activity/',//获取所有活动
  'getUserInfo':'/api/v1/user/',//获取用户信息
}

//app存储的数据库
export const PREFERENCES_NAME = 'hobby.db'

//获取后端接口，成功响应码
export const HTTP_CDE_200: number = 200

//超时时间设置
export const HTTP_READ_TIMEOUT: number = 100000