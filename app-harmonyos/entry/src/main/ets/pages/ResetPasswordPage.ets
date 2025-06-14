import { router } from '@kit.ArkUI'
import http from '@ohos.net.http'
import { ServerPath } from '../common/Constans'
import { httpRequest } from '../utils/HttpRequest'
import { showText } from '../utils/ShowUtil'
import { loginUserAction } from '../utils/UserUtil'
import { User } from '../viewmodel/User'

@Entry
@Component
struct ResetPasswordPage {
  @State oldPassword: string = ''
  @State newPassword: string = ''
  @State confirmPassword: string = ''
  @State userId: number = 0
  @State username: string = ''

  aboutToAppear() {
    // 获取当前登录用户信息
    loginUserAction((user: User) => {
      this.userId = user.id
      this.username = user.username
    }, () => {
      showText('请先登录')
      router.back()
    })
  }

  // 修改密码API
  resetPassword(): Promise<any> {
    let url = `${ServerPath.resetPassword}` // 使用正确的API路径
    let param: http.HttpRequestOptions = {
      method: http.RequestMethod.POST,
      header: {
        'Content-Type': 'application/json'
      },
      extraData: {
        userId: this.userId,
        username: this.username,
        oldPassword: this.oldPassword,
        newPassword: this.newPassword
      }
    }
    return httpRequest(url, param)
  }

  // 提交修改密码
  submitReset() {
    // 表单验证
    if (!this.oldPassword) {
      showText('请输入原密码')
      return
    }
    if (!this.newPassword) {
      showText('请输入新密码')
      return
    }
    if (this.newPassword !== this.confirmPassword) {
      showText('两次输入的新密码不一致')
      return
    }
    if (this.newPassword.length < 6) {
      showText('密码长度不能少于6位')
      return
    }

    // 调用API修改密码
    this.resetPassword()
      .then(res => {
        if (res.success) {
          showText('密码修改成功', () => {
            router.back()
          })
        } else {
          showText(res.errorMessage || '密码修改失败')
        }
      })
      .catch(err => {
        showText('网络错误，请稍后重试')
        console.error('Reset password error:', err)
      })
  }

  build() {
    Column() {
      // 顶部导航栏
      Row() {
        Image($r('app.media.ic_back'))
          .width(30)
          .height(30)
          .onClick(() => {
            router.back()
          })
        Text('修改密码')
          .fontSize(20)
          .fontWeight(FontWeight.Medium)
          .margin({ left: 16 })
        Blank()
      }
      .width('100%')
      .height(60)
      .padding({ left: 16, right: 16 })
      .backgroundColor(Color.White)

      // 表单内容
      Column() {
        // 原密码
        TextInput({ placeholder: '请输入原密码' })
          .type(InputType.Password)
          .placeholderColor('#999999')
          .placeholderFont({ size: 16 })
          .backgroundColor(Color.White)
          .height(50)
          .width('100%')
          .margin({ top: 20 })
          .onChange((value) => {
            this.oldPassword = value
          })

        // 新密码
        TextInput({ placeholder: '请输入新密码' })
          .type(InputType.Password)
          .placeholderColor('#999999')
          .placeholderFont({ size: 16 })
          .backgroundColor(Color.White)
          .height(50)
          .width('100%')
          .margin({ top: 16 })
          .onChange((value) => {
            this.newPassword = value
          })

        // 确认新密码
        TextInput({ placeholder: '请再次输入新密码' })
          .type(InputType.Password)
          .placeholderColor('#999999')
          .placeholderFont({ size: 16 })
          .backgroundColor(Color.White)
          .height(50)
          .width('100%')
          .margin({ top: 16 })
          .onChange((value) => {
            this.confirmPassword = value
          })

        // 提交按钮
        Button('确认修改')
          .width('100%')
          .height(50)
          .backgroundColor('#ff3551')
          .fontColor(Color.White)
          .fontSize(18)
          .margin({ top: 40 })
          .onClick(() => {
            this.submitReset()
          })

        // 密码规则提示
        Text('密码长度不少于6位，建议使用字母、数字和符号的组合')
          .fontSize(14)
          .fontColor('#999999')
          .margin({ top: 16 })
          .width('100%')
          .textAlign(TextAlign.Center)
      }
      .width('90%')
      .padding(16)
      .backgroundColor(Color.White)
      .borderRadius(10)
      .margin({ top: 20 })
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#f5f5f5')
    .padding({ top: 10 })
  }
} 