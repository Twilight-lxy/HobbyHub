import {User} from '../viewmodel/User'
import {ServerPath} from "../common/Constans"; //导入服务API常量
import {httpRequest,ResultData} from '../utils/HttpRequest' //导入通用http入口
import http from '@ohos.net.http';
//更新用户信息的方法
export function updateUserInfo(userInfo:User):Promise<ResultData>{
  //修改用户信息
  let url = `${ServerPath.UpdateUserInfo}`
  let param:http.HttpRequestOptions = {
    method:http.RequestMethod.PUT,//http请求方法
    header:{//参数传递使用json数据传递
      'Content-Type':'application/json'//默认form表单,但这里不是
    },
    extraData:userInfo
  }
  return httpRequest(url,param)
}