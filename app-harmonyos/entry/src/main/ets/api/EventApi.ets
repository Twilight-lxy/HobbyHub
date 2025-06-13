

import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { ServerPath } from "../common/Constans"; //导入服务Api常量
import { httpRequest, ResultData } from '../utils/HttpRequest' //导入通用http入口
//import { Event } from '@kit.ArkTS';
import {Event} from '../viewmodel/Event'


//获取活动列表的方法
export function getEventList(userId:number,kw:string = ''):Promise<ResultData>{
  let url = `${ServerPath.eventList}?userId = ${userId} & kw = ${encodeURI(kw)}` //常量里获取的活动列表
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.GET //http请求方法
  }
  return httpRequest(url, param)

}

//添加活动
export function addEvent(ent:Event):Promise<ResultData>{
  let url = `${ServerPath.addEvent}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.POST,
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
    extraData: ent
  }
  return httpRequest(url, option)
}

//获取我创建的活动
export function getMyEventList(userId:number):Promise<ResultData>{
  let url = `${ServerPath.myEventList}?userId=${userId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.GET
  }
  return httpRequest(url, option)
}
//获取单个活动信息
export function getEventInfo(eventId:number):Promise<ResultData>{
  let url = `${ServerPath.eventInfo}?eventId=${eventId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.GET
  }
  return httpRequest(url, option)
}
//编辑活动
export function editEvent(ent:Event):Promise<ResultData>{
  let url = `${ServerPath.eventInfo}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.PUT,
    header:{ //参数传递使用json数据传递
      'Content-Type':'application/json'
  },
    extraData:ent
  }
  return httpRequest(url, option)
}

//删除活动
export function deleteEvent(eventId:number):Promise<ResultData>{
  let url = `${ServerPath.deleteEvent}?eventId=${eventId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.DELETE
  }
  return httpRequest(url, option)
}




