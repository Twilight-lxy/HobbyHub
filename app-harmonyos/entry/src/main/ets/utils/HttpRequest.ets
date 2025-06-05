import http from '@ohos.net.http';//申请http请求权限的依赖库

import {HTTP_CODE_200,BaseUrl} from "../common/Constans"; //导入服务地址和状态码
import {BusinessError} from "@ohos.base" //导入鸿蒙异常处理
import {showText} from '../utils/ShowUtil' //导入通用提示类

//定义返回的封装类型
export interface ResultData{
  success:boolean
  data:string
}

export interface ErrorResponse{
  errorMessage:string
}
//实现一个通用的http请求入口
//url 表示后端的请求地址
//option 是发送请求参数的封装
export function httpRequest(url:string,option:http.HttpRequestOptions){
  //定义http请求
  let httpRequest = http.createHttp()
  let responseResult = httpRequest.request(`${BaseUrl}${url}`,option) //向后台发起请求

  //返回类型
  let serverData:ResultData = {
    success : false,
    data:""
  }
  //返回后端数据，进行封装后供前台使用
  return responseResult.then((value) => {//正常获取后台数据
    if (value.responseCode === HTTP_CODE_200) {//判断状态码是否为200，是则封装数据给前台使用
      //封装后台数据
      serverData.success = true
      serverData.data = `${value.result}`
    }
    else
    {
      let returnjson = JSON.parse(`${value.result}`) as ErrorResponse
      serverData.success = false
      serverData.data = returnjson.errorMessage
    }
    return serverData
  }).catch((error:BusinessError) =>{ //非正常获取后台数据时
    showText('调用系统错误')
    serverData.success = false
    return serverData
  })
}