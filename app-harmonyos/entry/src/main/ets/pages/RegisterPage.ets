//注册页面

import router from '@ohos.router'; //导入arkts路由功能
import { register } from '../api/LoginApi' //导入注册api
import { showText } from '../utils/ShowUtil' //导入公用提示框
import { url } from '@kit.ArkTS';

@Entry
@Component
struct RegisterPage {
  //注册用户名，密码和重复密码
  @State username: string = ''
  @State password: string = ''
  @State repassword: string = ''

  //发起用户注册请求
  async register() {
    //参数校验
    if (this.username == '' || this.password == '' || this.repassword == '') {
      showText("填写内容不能为空")
    } else if (this.password.length < 6) {
      showText("密码长度至少6位")
    } else if (this.password != this.repassword) {
      showText("两次输入的密码不一致")
    } else {
      //发起注册请求
      let result = await register(this.username, this.password)
      if (result.success) {
        //提示注册成功跳转到登录页面
        showText("注册成功，即将跳转到登录页面", () => {
          router.back({ url: 'pages/LoginPage' })
        })
      } else {
        showText(result.info)
      }
    }
  }

  build() {
    Column() {
      Row() {
        Blank() //空行
        Text("随便看看")
          .fontSize(16)
          .fontColor('#999999')
          .onClick(e => {
            //todo 不登录跳转到个人首页
          })
      }.height('10%')
      .width('90%')
      .margin({ top: 10 })

      //兴趣小队文字区域
      Text("欢迎来到兴趣小队")
        .width('80%')
        .fontSize(22)
        .fontWeight(FontWeight.Medium)
        .padding(10)
        .textAlign(TextAlign.End) //靠右

      Column() {
        //用户名输入框
        TextInput({ placeholder: '用户名', text: '' })
          .onChange((value: string) => {
            this.username = value
          })
          .shadow(ShadowStyle.OUTER_DEFAULT_SM)
          .backgroundColor("#ffffff")
          .height(45)
        Blank()

        //密码输入框
        TextInput({ placeholder: '密码', text: '' })
          .onChange((value: string) => {
            this.password = value
          })
          .type(InputType.Password)
          .shadow(ShadowStyle.OUTER_DEFAULT_SM)
          .backgroundColor("#ffffff")
          .height(45)
        Blank()

        //确认密码输入框
        TextInput({ placeholder: '确认密码', text: '' })
          .onChange((value: string) => {
            this.repassword = value
          })
          .type(InputType.Password)
          .shadow(ShadowStyle.OUTER_DEFAULT_SM)
          .backgroundColor("#ffffff")
          .height(45)
      }.height('25%') //增加高度
      .width('90%')
      .padding(10) //行间距
      .borderRadius(10) //文本框边缘设置为圆形

      //注册按钮
      Button("注册")
        .width('75%')
        .height(40)
        .margin({ top: 12 })//上边距
        .borderRadius(25)//圆角边框
        .backgroundColor("#ff3551")//单击事件
        .onClick(() => {
          //发起注册请求
          this.register()
        })

      Row() {
        Text("<-返回登录界面")
          .fontSize(16)
          .fontColor('#666666')
          .onClick(e => {
            //arkTs页面间跳转,返回
            router.back({ url: 'pages/LoginPage' })
          })
      }.height('6%')
      .margin({ top: 10 })
      .width('80%')
      .justifyContent(FlexAlign.End)

    }.width('100%')
    .height('100%')
    .backgroundImage($r('app.media.ic_bg_login'))
    .backgroundImageSize({
      height: '20%',
      width: '100%'
    })
  }
}


