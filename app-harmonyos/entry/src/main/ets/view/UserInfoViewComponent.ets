import { BaseUrl } from '../common/Constans'
import { showText } from '../utils/ShowUtil'
import { loginUserAction } from '../utils/UserUtil'
import {User} from '../viewmodel/User'
import router from '@ohos.router'

//操作菜单栏
class MenuItem{
  icon:Resource
  name:string
  action:()=>void

  constructor(icon:Resource,name:string,action:()=>void) {
    this.icon = icon
    this.name = name
    this.action = () =>{
      loginUserAction(
        () => {
          action()
        },
        () => {
          showText("请先登录")
        }
      )
    }
  }
}
@Preview
@Component
export struct UserInfoViewComponent{

  //用户信息 app状态，组件之间传递
  @StorageProp("user") user:User = new User(0,"","","","","",0,0)
  @StorageProp("isLogin") isLogin : boolean = false

  build() {
    Column() {
      Row() {
        //加载头像，如果当前用户已登录，则获取用户自己的头像，如果未登录，使用alt属性加载缺省头像
        Image($r(`${BaseUrl}${this.user.headImg}`))
          .alt($r('app.media.ic_user_icon_default'))
          .width(70)
          .height(70)
          .backgroundColor(Color.White)
          .borderRadius(35) //全圆

        //判断展示登录按钮还是用户信息
        if (this.isLogin) {
          Text(this.user.name)
            .fontColor(Color.White)
            .fontSize(26)
            .fontWeight(FontWeight.Medium)
          Text(this.user.username)
            .fontColor(Color.White)
            .fontSize(16)
            .margin({ top: 5 })

        } else {
          Text("登录")
            .fontColor(Color.White)
            .fontSize(26)
            .fontWeight(FontWeight.Medium)
            .margin({ left: 20 })
            .onClick(() => {
              //点击 跳转到登陆页面
              router.pushUrl({ url: "pages/LoginPage" })
            })
        }
      }.width("100%")
      .padding({left:20})
      .margin({top:60})

      //发起活动按钮
      Button("＋发起活动")
        .width("90%")
        .borderRadius(10)//半圆角
        .fontColor("#ff3551")
        .backgroundColor(Color.White)
        .margin({top:20})
        .onClick(()=>{
          loginUserAction(
            ()=>{
              //todo 跳转到活动编辑页面
              //router.pushUrl()
            },
            ()=>{

            }
          )
        })
    }
  }
}