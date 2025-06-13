import {User} from '../viewmodel/User';
import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { ServerPath } from "../common/Constans";//导入服务Api常量
import { httpRequest, ResponseData } from '../utils/HttpRequest'; //导入通用http入口

//更新用户信息的方法
export function updateUserInfo(userInfo:User):Promise<ResponseData> {
  let url = `${ServerPath.UpdateUserInfo}`
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.PUT, //http请求方法
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
    extraData: userInfo
  }
  return httpRequest(url, param)
}

//更新用户信息的方法
export function getUserInfo(id?:number,username?:string):Promise<ResponseData> {
  let url = `${ServerPath.getUserInfo}`

  if (id){
    url+=`?id=${id}`
  }
  if (username){
    url+=`?username=${username}`
  }

  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.GET, //http请求方法
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
  }
  return httpRequest(url, param)
}