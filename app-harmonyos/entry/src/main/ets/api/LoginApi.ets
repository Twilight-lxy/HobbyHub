// 用户相关操作，调用后台的类

import http from '@ohos.net.http';
import {ServerPath} from "../common/Constans"; //导入服务API常量
import {httpRequest,ResultData} from '../utils/HttpRequest' //导入通用http入口
//获取登录信息
export function login(username:string,password:string):Promise<ResultData>{
  let url = `${ServerPath.login}`
  let param:http.HttpRequestOptions = {
    method:http.RequestMethod.POST,//http请求方法
    header:{//参数传递使用json数据传递
      'Content-Type':'application/json'//默认form表单,但这里不是
    },
    extraData:{//传递具体参数
      username:username,
      password:password
    }
  }
  return httpRequest(url,param)
}