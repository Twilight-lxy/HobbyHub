import { Event } from '../viewmodel/Event'
import { BaseUrl } from '../common/Constans'

@Component
export struct EventComponent {
  //父组件向子组件传递值
  @Prop item: Event

  build() {
    Column() {
      //活动海报
      Image(`${BaseUrl}${this.item.headImg}`)
        .constraintSize({ maxHeight: 240 })
        .borderRadius(10)
        .objectFit(ImageFit.Cover)
        .onClick(() => {
          //跳转到活动详情页面
        })

      //文字描述：活动地点
      Text(this.item.name)
        .maxLines(2)//文字最多显示两行
        .textOverflow({ overflow: TextOverflow.Clip })//超过两行截取
        .margin({ top: 5 })
        .fontColor('#333333')
        .fontSize(16)
        .width('90%')
        .textAlign(TextAlign.Start) //靠左

      //活动描述：活动地点
      Text(this.item.name)
        .maxLines(1)//文字最多显示两行
        .textOverflow({ overflow: TextOverflow.Clip })//超过两行截取
        .margin({ top: 5 })
        .fontColor('#333333')
        .fontSize(14)
        .width('90%')
        .textAlign(TextAlign.Start) //靠左
    }
    .borderRadius(10)
    .width('45%')
  }
}