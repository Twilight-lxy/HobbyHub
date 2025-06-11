//申请位置权限时，需要同时申请模糊位置权限和精准位置权限

import abilityAccessCtrl,{PermissionRequestResult,Permissions} from '@ohos.abilityAccessCtrl';
import common from '@ohos.app.ability.common';
import {BusinessError} from '@ohos.base'
import { data } from '@kit.TelephonyKit';

//获取模糊位置和精准位置权限
const locationPermissions : Permissions[]=[
  "ohos.permission.APPROXIMATELY_LOCATION","ohos.permission.LOCATION"
]

//申请权限
export function reqLocationPermissions(context:common.UIAbilityContext):Promise<boolean>{
  //声明权限管理者
  let atManager:abilityAccessCtrl.AtManager = abilityAccessCtrl.createAtManager()
  return new Promise((resolve,reject)=>{
      atManager.requestPermissionsFromUser(context,locationPermissions).then(
        (data:PermissionRequestResult) =>{
          let grantStatus:Array<number> = data.authResults
          //判断精准位置权限结果
          if (grantStatus[1] == 0) {
            resolve(true)
          }else{
            resolve(false)
            return
          }
        }
      ).catch((e:BusinessError) => {
        console.log(`Failed to request permission from user,Code is ${e.code},message is ${e.message}`)
        reject(e)
      })
  })
}

