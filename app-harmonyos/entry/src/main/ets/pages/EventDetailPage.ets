//活动详情


import { router } from '@kit.ArkUI'
import { Event } from '../viewmodel/Event'
import { getEventInfo } from '../api/EventApi'
import { RouterParam } from '../viewmodel/RouterParam'
import { showText } from '../utils/ShowUtil'
import { BusinessError } from '@kit.BasicServicesKit'
import { BaseUrl } from '../common/Constans'
import { joinEvent } from '../api/EventMemberApi'
import { loginUserAction } from '../utils/UserUtil'
import { User } from '../viewmodel/User'

@Entry
@Component
struct EventDetailPage {
  @State event: Event = new Event(0, 0, "", "", "", "", "", 0, 0, 0)
  @State currentUserId: number = 0
  
  aboutToAppear() {
    // 获取路由参数中的活动ID
    let param = router.getParams() as RouterParam<number>
    let eventId = param.data
    
    // 加载活动详情
    this.loadEventDetail(eventId)
    
    // 获取当前登录用户ID
    loginUserAction((user: User) => {
      this.currentUserId = user.id
    })
  }
  
  // 加载活动详情
  loadEventDetail(eventId: number) {
    getEventInfo(eventId)
      .then(res => {
        if (res.success) {
          this.event = res.data as Event
        } else {
          showText(res.errorMessage)
        }
      })
      .catch((e: BusinessError) => {
        showText(e.message)
      })
  }
  
  // 加入活动
  joinEvent() {
    if (this.currentUserId <= 0) {
      showText('请先登录')
      return
    }
    
    joinEvent(this.event.id, this.currentUserId)
      .then(res => {
        if (res.success) {
          showText('加入成功')
        } else {
          showText(res.errorMessage)
        }
      })
      .catch((e: BusinessError) => {
        showText(e.message)
      })
  }
  
  build() {
    Stack({ alignContent: Alignment.Bottom }) {
      Column() {
        // 顶部导航栏
        Row() {
          Image($r('app.media.ic_back'))
            .width(30)
            .height(30)
            .onClick(() => {
              router.back()
            })
          Text('活动详情')
            .fontSize(20)
            .fontWeight(FontWeight.Medium)
          Blank()
        }
        .width('100%')
        .height(60)
        .padding({ left: 16, right: 16 })
        .backgroundColor('#FFFFFF')
        
        // 活动内容区域
        Scroll() {
          Column() {
            // 活动图片
            Image(this.event.headImg ? `${BaseUrl}${this.event.headImg}` : $r('app.media.ic_swiper_1'))
              .width('100%')
              .height(200)
              .objectFit(ImageFit.Cover)
              
            // 白色背景区域
            Column() {
              // 活动标题
              Text(this.event.name || "户外休闲游")
                .fontSize(22)
                .fontWeight(FontWeight.Bold)
                .margin({ top: 16, bottom: 16 })
                .textAlign(TextAlign.Start)
                .width('100%')
              
              // 活动地点
              Row() {
                Image($r('app.media.user_position'))
                  .width(20)
                  .height(20)
                Text(this.event.addr || "山东省青岛市城阳区")
                  .fontSize(16)
                  .margin({ left: 8 })
              }
              .margin({ bottom: 16 })
              .width('100%')
              .alignItems(VerticalAlign.Center)
              
              // 活动内容
              Text(this.event.intro || "欢迎参加我们的户外休闲游活动！在这个活动中，您将有机会探索大自然的美丽景色，体验户外运动的乐趣和挑战。我们将组织各种活动，如徒步旅行、露营、烧烤聚会、山地骑行等，让您尽情放松身心，远离城市喧嚣。无论是与家人朋友一起出游，还是独自探险，都能在户外休闲游中找到乐趣和满足，让我们一起享受户外生活的美好，感受大自然的魅力！")
                .fontSize(16)
                .margin({ bottom: 20 })
                .textAlign(TextAlign.Start)
                .width('100%')
                .lineHeight(24)
            }
            .width('100%')
            .padding({ left: 16, right: 16, top: 16, bottom: 100 })
            .backgroundColor(Color.White)
            .borderRadius({ topLeft: 16, topRight: 16 })
            .margin({ top: -20 })
          }
          .width('100%')
        }
        .layoutWeight(1)
      }
      .width('100%')
      .height('100%')
      
      // 底部操作栏
      Row() {
        Row() {
          Image($r('app.media.ic_share'))
            .width(24)
            .height(24)
          Text('分享')
            .fontSize(16)
            .margin({ left: 4 })
        }
        .layoutWeight(1)
        .justifyContent(FlexAlign.Center)
        
        Row() {
          Image($r('app.media.ic_comment'))
            .width(24)
            .height(24)
          Text('评论')
            .fontSize(16)
            .margin({ left: 4 })
        }
        .layoutWeight(1)
        .justifyContent(FlexAlign.Center)
        
        Button('+ 加入')
          .backgroundColor('#FF3551')
          .borderRadius(24)
          .fontColor(Color.White)
          .width(120)
          .height(48)
          .fontSize(18)
          .onClick(() => {
            this.joinEvent()
          })
      }
      .width('100%')
      .height(70)
      .padding({ left: 16, right: 16 })
      .backgroundColor(Color.White)
      .justifyContent(FlexAlign.SpaceAround)
      .shadow({ radius: 5, color: '#E0E0E0', offsetY: -2 })
    }
    .width('100%')
    .height('100%')
  }
} 