

import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { ServerPath } from "../common/Constans"; //导入服务Api常量
import { httpRequest, ResultData } from '../utils/HttpRequest' //导入通用http入口


//获取活动列表的方法
export function getEventList(userId:number,kw:string = ''):Promise<ResultData>{
  let url = `${ServerPath.eventList}?userId = ${userId} & kw = ${encodeURI(kw)}` //常量里获取的活动列表
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.GET //http请求方法
  }
  return httpRequest(url, param)

}