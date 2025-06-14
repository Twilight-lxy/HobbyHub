import { Event } from '../viewmodel/Event'
import { BaseUrl } from '../common/Constans'
import router from '@ohos.router'
import { RouterParam } from '../viewmodel/RouterParam'

@Component
export struct EventComponent {
  //父组件向子组件传递值
  @Prop item: Event
  
  // 根据图片路径获取对应的本地资源
  private getImageResource(path: string): Resource | string {
    if (!path) return $r('app.media.ic_img_add') // 默认图片
    
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
          default: return $r('app.media.ic_img_add')
        }
      }
    }
    
    // 其他情况返回完整的远程URL路径
    return `${BaseUrl}${path}`
  }

  build() {
    Column() {
      //活动海报
      Image(this.getImageResource(this.item.headImg))
        .width('100%')
        .height(160)
        .borderRadius(10)
        .objectFit(ImageFit.Cover)
        .onClick(() => {
          //跳转到活动详情页面
          const params: RouterParam<Event> = new RouterParam<Event>(this.item, 'detail');
          router.pushUrl({
            url: 'pages/EventDetailPage',
            params: params
          })
        })

      //活动标题
      Text(this.item.name)
        .maxLines(2)//文字最多显示两行
        .textOverflow({ overflow: TextOverflow.Clip })//超过两行截取
        .margin({ top: 5 })
        .fontColor('#333333')
        .fontSize(16)
        .width('95%')
        .textAlign(TextAlign.Start) //靠左

      //活动描述
      Text(this.item.intro)
        .maxLines(1)//文字最多显示一行
        .textOverflow({ overflow: TextOverflow.Clip })//超过一行截取
        .margin({ top: 5 })
        .fontColor('#333333')
        .fontSize(14)
        .width('95%')
        .textAlign(TextAlign.Start) //靠左
        
      //活动地点
      Row() {
        Image($r('app.media.event_position'))
          .width(16)
          .height(16)
        Text(this.item.addr)
          .maxLines(1)
          .textOverflow({ overflow: TextOverflow.Ellipsis })
          .margin({ left: 5 })
          .fontColor('#666666')
          .fontSize(14)
          .width('85%')
      }
      .width('95%')
      .margin({ top: 5 })
      .alignItems(VerticalAlign.Center)
    }
    .borderRadius(10)
    .width('48%') // 增加卡片宽度
    .backgroundColor(Color.White)
    .padding(8)
  }
}