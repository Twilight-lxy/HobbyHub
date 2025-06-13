import common from '@ohos.app.ability.common';
import { picker } from '@kit.CoreFileKit';
import fs from '@ohos.file.fs';
import { BusinessError, request } from '@kit.BasicServicesKit';
import { BaseUrl } from '../common/Constans';
import { ServerPath } from '../common/Constans';
import { ResponseData } from './HttpRequest';


//选择图库文件并复制到后台程序缓存目录中
export async function selectImage(context: common.UIAbilityContext): Promise<string> {
  let photoSelectOptions = new picker.PhotoSelectOptions();
  photoSelectOptions.MIMEType = picker.PhotoViewMIMETypes.IMAGE_TYPE
  photoSelectOptions.maxSelectNumber = 1

  let photoPicker = new picker.PhotoViewPicker();
  let photoSelectResult = await photoPicker.select(photoSelectOptions)
  //从图库选择图片后，返回图片uri
  let uri = photoSelectResult.photoUris[0]

  //读取与上面返回的uri
  let file = fs.openSync(uri, fs.OpenMode.CREATE)

  //将上述图片复制到缓存目录下
  fs.copyFileSync(file.fd, context.cacheDir + '/test.jpg')

  //上传到服务器上
  return uploadImage(context)
}

//上传到后台服务器
function uploadImage(context: common.UIAbilityContext): Promise<string> {
  //上传任务配置
  let uploadConfig: request.UploadConfig = {
    url: `${BaseUrl}/${ServerPath.uploadImg}`,
    header: { 'Content-Type': 'multipart/form-data' },
    method: 'post',
    files: [
      {
        filename: 'text.jpg',
        name: 'file',
        uri: 'internal://cache/text.jpg',
        type: 'jpg'
      }
    ],
    data: []
  }
  //将本地应用文件上传到网络服务器
  return new Promise((resolve, reject) => {
    try {
      request.uploadFile(context, uploadConfig)
        .then((uploadTask) => {
          //监听响应数据，并返回图片上传后的URL
          uploadTask.on('headerReceive', (header) => {
            let bodyStr = header['body'] as string
            let body = JSON.parse(bodyStr) as ResponseData
            let url:string  = body.data?.toString() || ""
            resolve(url)
          });
        }).catch((err: BusinessError) => {
        console.error(`xx Invoke uploadFile failed, code is ${err.code},message is ${err.message}`);
        reject(err)
      })
    } catch (err) {
      console.error(`xx Invoke uploadFile failed, code is ${err.code},message is ${err.message}`);
      reject(err)
    }
  })
}