import { ResultData, httpRequest } from "../utils/HttpRequest";
import {ServerPath} from "../common/Constans"; //导入服务API常量
import http from '@ohos.net.http';
//获取活动列表的方法
export function getEventList(userId:number):Promise<ResultData>{
  let url = `${ServerPath.idEventList}?userId=${userId}` //常量里的获取活动的列表
  let param : http.HttpRequestOptions = {
    method:http.RequestMethod.GET//http请求方法 GET请求不需要head
  }
  return httpRequest(url,param)
}