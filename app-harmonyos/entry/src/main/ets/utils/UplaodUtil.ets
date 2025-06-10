import common from '@ohos.app.ability.common';
import picker from '@ohos.file.picker';
import fs from '@ohos.file.fs'
import request from '@ohos.request';
import {BaseUrl,ServerPath} from '../common/Constans'
import { JSON } from '@kit.ArkTS';
import { ResultData } from './HttpRequest';
import { BusinessError } from '@ohos.base';

//选择图库文件并复制到后台程序，后台程序缓存目录中
export async function selectImage(context:common.UIAbilityContext):Promise<string>{
  let photoSelectOptions = new picker.PhotoSelectOptions();
  photoSelectOptions.MIMEType = picker.PhotoViewMIMETypes.IMAGE_TYPE
  photoSelectOptions.maxSelectNumber = 1

  let photoPicker = new picker.PhotoViewPicker();
  let photoSelectResult = await photoPicker.select(photoSelectOptions) //选择的结果
  // 从图库选择图片后，返回图片uri
  let uri = photoSelectResult.photoUris[0]

  //读取上面返回的uri
  let file = fs.openSync(uri,fs.OpenMode.CREATE)

  //将上述图片复制到缓存目录下
  fs.copyFileSync(file.fd,context.cacheDir+'test.jpg')

  //上传到服务器上
  return uploadImage(context)
}

//上传到后台服务器
function uploadImage(context:common.UIAbilityContext):Promise<string>{
  //上传任务配置
  let uploadConfig:request.UploadConfig = {
    url : `${BaseUrl}/${ServerPath.uploadImg}`,
    header:{'Content-Type':'multipart/form-data'},
    method:"POST",
    files:[
      {filename:'text.jpg',name:'file',uri:'internal://cache/test.jpg',type:'jpg'} //部署到手机上后代码稍微改一改
    ],
    data:[]
  }
  //将本地应用文件上传到网络服务器
  return new Promise((resolve,reject)=> {
    try {
      request.uploadFile(context,uploadConfig)
        .then((uploadTask) => {
          //监听响应数据，并返回图片上传后的URL
          uploadTask.on('headerReceive', (header) => {
            let bodyStr = header['body'] as string
            let body = JSON.parse(bodyStr) as ResultData
            let url = body.obj.toString()
            resolve(url)
          });
        }).catch((err: BusinessError) => {
        console.error(`xx Invoke uploadFile failed,code is ${err.code},message is ${err.message}`);
        reject(err)
      })
    } catch (err) {
      console.error(`xx Invoke uploadFile failed,code is ${err.code},message is ${err.message}`);
      reject(err)
    }
  })
}