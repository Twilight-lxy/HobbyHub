import { BaseUrl } from "../common/Constans"
import { Event } from "../viewmodel/Event"

@Component
export struct EventComponent{
  //父组件向子组件传递值
  @Prop item:Event
  build() {
    Column(){
      Image(`${BaseUrl}${this.item.headImg}`)
        .constraintSize({maxHeight:240})
        .borderRadius(10)
        .objectFit(ImageFit.Cover)
        .onClick(()=>{
          //todo 跳转活动详情页面
        })

      Text(this.item.name)
        .maxLines(2)//文字显示最大行数
        .textOverflow({overflow:TextOverflow.Clip})//超过2行。截取
        .margin({top:5})
        .fontColor('#333333')
        .fontSize(16)
        .width("90%")
        .textAlign(TextAlign.Start)

      Text(this.item.name)
        .maxLines(1)//文字显示最大行数
        .textOverflow({overflow:TextOverflow.Clip})//超过2行。截取
        .margin({top:5})
        .fontColor('#333333')
        .fontSize(14)
        .width("90%")
        .textAlign(TextAlign.Start)
    }
    .borderRadius(10)
    .width('45%')
  }
}