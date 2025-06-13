

import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { ServerPath } from "../common/Constans"; //导入服务Api常量
import { httpRequest , ResponseData} from '../utils/HttpRequest' //导入通用http入口
//import { Event } from '@kit.ArkTS';
import {Event} from '../viewmodel/Event'


//获取活动列表的方法
export function getEventList():Promise<ResponseData>{
  let url = `${ServerPath.eventList}` //常量里获取的活动列表
  let param: http.HttpRequestOptions = {
    method: http.RequestMethod.GET //http请求方法
  }
  return httpRequest(url, param)

}

//添加活动
export function addEvent(ent:Event):Promise<ResponseData>{
  let url = `${ServerPath.addEvent}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.PUT,
    header: {
      //参数传递使用json数据传递
      'Content-Type': 'application/json'
    },
    extraData: ent
  }
  return httpRequest(url, option)
}

//获取我创建的活动
export function getMyEventList(userId:number):Promise<ResponseData>{
  let url = `${ServerPath.myEventList}?userId=${userId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.GET
  }
  return httpRequest(url, option)
}
//获取单个活动信息
export function getEventInfo(eventId:number):Promise<ResponseData>{
  let url = `${ServerPath.eventInfo}?eventId=${eventId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.GET
  }
  return httpRequest(url, option)
}
//编辑活动
export function editEvent(ent:Event):Promise<ResponseData>{
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
export function deleteEvent(eventId:number):Promise<ResponseData>{
  let url = `${ServerPath.deleteEvent}?eventId=${eventId}`
  let option:http.HttpRequestOptions={
    method:http.RequestMethod.DELETE
  }
  return httpRequest(url, option)
}




