import { BaseUrl } from '../common/Constans'
import { showText } from '../utils/ShowUtil'
import { loginUserAction } from '../utils/UserUtil'
import { User } from '../viewmodel/User'
import { router } from '@kit.ArkUI'

//操作菜单栏
class MenuItem {
  icon: Resource
  name: string
  action: () => void

  constructor(icon: Resource, name: string, action: () => void) {
    this.icon = icon
    this.name = name
    this.action = () => {
      //判断状态值
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
export struct UserInfoViewComponent {
  //用户信息 app状态，组件之间传递
  @StorageProp("user") user: User = new User(0, "", "", "", "", "", 0, 0)
  @StorageProp("isLogin") isLogin: boolean = false
  //操作列表定义
  menuList: MenuItem[] = [
    new MenuItem($r('app.media.ic_user_menu_event'), "管理活动", () => {
      router.pushUrl({ url: 'pages/MyEventPage' })
    }),
    new MenuItem($r('app.media.ic_user_menu_setting'), "用户管理", () => {
      router.pushUrl({ url: 'pages/UserEditPage' })
    }),
    new MenuItem($r('app.media.ic_user_menu_security'), "修改密码", () => {
      router.pushUrl({ url: 'pages/ResetPasswordPage' })
    }),
    new MenuItem($r('app.media.ic_user_menu_about'), "关于我们", () => {
      showText("隐私声明：\t\n" +
        "我们非常重视您的隐私权和个人信息保护。我们兴趣活动APP承诺对用户的个人信息进行合理的收集、使用和保护，严格遵守相关法律法规，绝不会泄露、出售或分享用户的个人信息给第三方。\t\n" +
        "具体而言，我们会收集用户的个人信息包括但不限于姓名、联系方式、兴趣爱好等用于在APP上提供个性化推荐、活动安排等服务。我们会采取合理的安全措施保护用户的个人信息，并仅在用户同意的情况下使用用户的个人信息。\t\n" +
        "用户在使用我们的兴趣活动APP时，即表示用户同意我们按照本隐私声明的规定收集、使用和保护用户的个人信息。如果用户对我们的隐私声明有任何疑问或意见，可随时联系我们的客服人员进行解答和沟通。感谢您对我们的信任和支持！")
    }),
  ]

  build() {
    Column() {
      Row() {
        //加载头像，如果当前用户未登录，则获取用户自己的头像；如果未登录则使用alt属性加载缺省头像
        Image($r(`${BaseUrl}${this.user.headImg}`))
          .alt($r('app.media.ic_user_icon_default'))
          .width(70)
          .height(70)
          .backgroundColor(Color.White)
          .borderRadius(35) //圆形

        //判断展示登录按钮还是用户信息
        if (this.isLogin) {
          Column() {
            Text(this.user.name)
              .fontColor(Color.White)
              .fontSize(26)
              .fontWeight(FontWeight.Medium)
            Text(this.user.username)
              .fontColor(Color.White)
              .fontSize(16)
              .margin({ top: 5 })
          }.alignItems(HorizontalAlign.Start) //设置文字排版为竖列
          .margin({ left: 10 })
        } else {
          Text("登录")
            .fontColor(Color.White)
            .fontSize(26)
            .fontWeight(FontWeight.Medium)
            .margin({ left: 20 })
            .onClick(() => {
              //点击跳转到登录页面
              router.pushUrl({ url: 'pages/LoginPage' })
            })
        }
      }.width('100%')
      .padding({ left: 20 })
      .margin({ top: 60 })

      //发起活动按钮
      Button("⊕发起活动")
        .width('90%')
        .type(ButtonType.Normal)
        .borderRadius(10)//半圆角边框
        .fontColor('#ff3551')
        .backgroundColor(Color.White)
        .margin({ top: 20 })
        .onClick(() => {
          loginUserAction(
            () => {
              router.pushUrl({url:"pages/EventEditPage",params:{message:'add'}})
            },
            () => {
              showText("请先登录")
            }
          )
        })

      //用户管理的相关操作列表
      Column() {
        ForEach(this.menuList, (item: MenuItem) => {
          Row() {
            Image(item.icon)
              .width(20)
              .margin({ left: 20 })
            Text(item.name)
              .margin({ left: 20 })
          }.height(60)
          .width('100%')
          .onClick(() => {
            item.action()
          })
        })
      }
      .width('90%')
      .backgroundColor(Color.White)
      .borderRadius(15)
      .margin({ top: 20 })
      .shadow(ShadowStyle.OUTER_FLOATING_SM)

      //退出按钮
      Text("退出登录")
        .fontSize(18)
        .margin({ top: 20 })
        .fontColor('#999999')
        .onClick(() => {
          if (this.isLogin) {
            // 删除当前登录用户记录
            AppStorage.delete("user")
            // 设置当前是否登录为否
            AppStorage.set("isLogin", false)
            // 显示退出成功提示
            showText("退出登录成功", () => {
              // 延迟一秒跳转到首页
              setTimeout(() => {
                router.replaceUrl({ url: 'pages/Index' })
              }, 1000)
            })
          } else {
            showText("您尚未登录")
          }
        })
    }.width('100%')
    .height('100%')
    .backgroundImage($r('app.media.bg_userinfo'))
    .backgroundImageSize({
      height: '40%',
      width: '100%'
    })
  }
}