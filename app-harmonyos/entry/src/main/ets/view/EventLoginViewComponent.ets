import { router } from '@kit.ArkUI'

@Component
export struct EventLoginViewComponent {
  build() {
    Column() {
      // 未登录状态图标
      Image($r('app.media.ic_user_icon_default'))
        .width(80)
        .height(80)
        .margin({ top: '30%' })
      
      // 未登录状态提示
      Text("未登录")
        .fontSize(24)
        .fontWeight(FontWeight.Medium)
        .margin({ top: 20 })
      
      Text("请先登录！")
        .fontSize(16)
        .fontColor('#999999')
        .margin({ top: 10 })
      
      // 登录按钮
      Button("去登录")
        .width('80%')
        .height(50)
        .margin({ top: 40 })
        .borderRadius(25)
        .backgroundColor("#ff3551")
        .onClick(() => {
          router.pushUrl({ url: 'pages/LoginPage' })
        })
        
      // 底部导航栏占位，确保内容不会被底部导航栏遮挡
      Blank()
        .height('10%')
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#f5f5f5')
    .justifyContent(FlexAlign.Start)
    .alignItems(HorizontalAlign.Center)
  }
} 