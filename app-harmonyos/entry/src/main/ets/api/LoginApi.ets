//用户相关操作，调用后台的类

import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { ServerPath } from "../common/Constans"; //导入服务Api常量
import { httpRequest, ResponseData } from '../utils/HttpRequest' //导入通用http入口

//获取登录信息
export function login(username: string, password: string): Promise<ResponseData> {
  let url = `${ServerPath.login}`
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.POST, //http请求方法
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
    extraData: {
      //传递的具体参数
      username: username,
      password: password
    }
  }
  return httpRequest(url, param)
}

//注册用户
export function register(username: string, password: string): Promise<ResponseData> {
  let url = `${ServerPath.register}` //常量里注册的url
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.POST, //http请求方法
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
    extraData: {
      //传递的具体参数
      username: username,
      password: password
    }
  }
  return httpRequest(url, param)
}