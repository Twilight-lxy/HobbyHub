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
  @State event: Event = new Event(0, 0, "", "", "", "", "", 0, 0, 0, "", 0 ,"")
  @State currentUserId: number = 0
  
  aboutToAppear() {
    // 获取路由参数
    let param = router.getParams() as RouterParam<Event>
    
    if (param && param.data) {
      // 如果是从活动列表点击进入，直接使用传递的活动数据
      this.event = param.data
    } else if (param && param.message === 'id' && typeof param.data === 'number') {
      // 如果只传递了ID，则加载活动详情
      this.loadEventDetail(param.data as number)
    }
    
    // 获取当前登录用户ID
    loginUserAction((user: User) => {
      this.currentUserId = user.id
    })
  }
  
  // 加载活动详情
  loadEventDetail(eventId: number) {
    getEventInfo(eventId)
      .then(res => {
        if (res.success == true) {
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
        if (res.success == true) {
          showText('加入成功')
        } else {
          showText(res.errorMessage)
        }
      })
      .catch((e: BusinessError) => {
        showText(e.message)
      })
  }
  
  // 根据图片路径获取对应的本地资源
  private getImageResource(path: string): Resource | string {
    if (!path) return $r('app.media.ic_swiper_1') // 默认图片
    
    // 如果是event_image_格式的路径，转换为本地资源
    if (path.includes('event_image_')) {
      const match = path.match(/event_image_(\d+)/)
      if (match && match[1]) {
        const index = parseInt(match[1])
        switch (index) {
          case 1: return $r('app.media.ic_swiper_1')
          case 2: return $r('app.media.ic_swiper_2')
          case 3: return $r('app.media.ic_swiper_3')
          case 4: return $r('app.media.bg_eventinfo')
          default: return $r('app.media.ic_swiper_1')
        }
      }
    }
    
    // 其他情况返回完整的远程URL路径
    return `${BaseUrl}${path}`
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
            Image(this.getImageResource(this.event.headImg))
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