//后端接口定义
export const BaseUrl = 'http://192.168.1.9:8081' //Zhuanyong5G

//后端接口的路径
export const ServerPath: Record<string, string> = {
  'login': '/usr/login', //登录接口
  'register': '/usr/reg',//注册接口
  'join':'/eventMember/join',//加入活动接口
  'addEvent':'/event/insert',//新增活动
  'myEventList':'/event/myEventList', //获取我创建的活动
  'eventInfo':'/event/info',//获取活动详情
  'editEvent':'/event/edit',//修改活动信息
  'deleteEvent':'/event/delete',//删除活动
}

//app存储的数据库
export const PREFERENCES_NAME = 'hobby.db'

//获取后端接口，成功响应码
export const HTTP_CDE_200: number = 200

//超时时间设置
export const HTTP_READ_TIMEOUT: number = 100000