import { BusinessError } from "@ohos.base" //导入鸿蒙异常处理
import { showText } from '../utils/ShowUtil'; //导入通用提示类
import http from '@ohos.net.http'; //导入申请http请求权限的依赖库
import { HTTP_CDE_200, BaseUrl } from "../common/Constans"; //导入服务地址和状态码
import { JSON } from '@kit.ArkTS';
import { data } from "@kit.TelephonyKit";

export interface ResponseData {
  success:boolean
  data:object
  errorMessage:string
}

//实现一个通用的http请求入口
//url：表示后端的请求地址
//option：发送请求参数的封装
export function httpRequest(url: string, option: http.HttpRequestOptions):Promise<ResponseData>{
  //定义http请求
  let httpRequest = http.createHttp()
  const token = AppStorage.get<string>('auth_token');
  if (token) {
    if (!option.header) {
      option.header = {};
    }
    option.header['Authorization'] = token.toString();
  }
  //向后台发起请求
  let responseResult = httpRequest.request(`${BaseUrl}${url}`, option)

  let resultData:ResponseData={
    success:false,
    data:Object,
    errorMessage:""
  }

  //返回后端数据，进行封装后供前台使用
  return responseResult.then((value) => {
    //正常获取后台数据时
    if (value.responseCode === HTTP_CDE_200) { //判断状态码，是否为200，试着封装数据给前台使用
      try{
        //封装后台数据
        let result = `${value.result}`
        let resultJson = JSON.parse(result)
        // Check if response contains token
        if ((resultJson as object)?.["token"]) {
          // Save token to AppStorage for future requests
          AppStorage.setOrCreate<string>('auth_token', (resultJson as object)?.["token"]);
        }
        resultData.success=true
        resultData.data=(resultJson as object)
        console.log(result)
        resultData.errorMessage=""
      }catch (e) {
        console.error('Error parsing JSON response:', e);
      }
    } else {
      console.error('Error parsing JSON response',value.result);
    }
    return resultData;
  }).catch((error: BusinessError) => { //非正常获取后台数据时
    showText('调用系统错误')
    resultData.errorMessage="调用系统错误"
    return resultData
  });
}