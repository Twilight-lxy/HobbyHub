import { login } from '../api/LoginApi'
import { loginUserAction } from '../utils/UserUtil'
import {Event} from '../viewmodel/Event'
import {getMyEventList,getEventInfo,editEvent,deleteEvent} from '../api/EventApi'
import { BaseUrl } from '../common/Constans'
import { router } from '@kit.ArkUI'
import { showText } from '../utils/ShowUtil'

@Entry
@Component
struct MyEventPage{

  @State eventList:Event[] = []

  //arkTs的生命周期
  //页面显示时加载
  onPageShow(): void {
    this.getEventList()
  }
  //获取我的活动列表
  getEventList(){
    loginUserAction((user) =>{
      getMyEventList(user.id).then(res=>{
        if (res.success) {
          this.eventList =  res.data as Event[]
        }
      })
    })
  }
  //删除活动
  deleteEvent(eventId:number){
    AlertDialog.show({
      message:"您确定删除这个活动吗",
      alignment:DialogAlignment.Center,
      primaryButton:{
        value:"确定",
        action:()=>{
          deleteEvent(eventId).then(res =>{
            if (res.success) {
              showText("删除成功")
            }else{
              showText('删除失败')
            }
            //刷新列表
            this.getEventList()
          })
        }
      },
      secondaryButton:{
        value:"取消",
        action:()=>{}
      }
    })
  }

  //活动组件
  @Builder
  MyEvent(item:Event,index:number){
    Row(){
      Image(`${BaseUrl}${item.headImg}`)
        .height(80)
        .width(140)
        .borderStyle(10)

      Column(){
        Text(item.name)
          .maxLines(1) //只允许出现一行maxlines与textoverflow需要同时使用才能生效
          .textOverflow({overflow:TextOverflow.Clip}) //overflow:TextOverflow.Clip表示超过有的截取处理
          .constraintSize({maxWidth:120})

        Row(){
          Image($r('app.media.event_position'))
            .width(14)
            .height(14)
          Text(item.addr)
            .maxLines(1) //只允许出现一行maxlines与textoverflow需要同时使用才能生效
            .textOverflow({overflow:TextOverflow.Ellipsis}) //overflow:TextOverflow.Ellipsis表示超过的用...处理
            .constraintSize({maxWidth:120})
            .fontSize(14)
            .fontColor("#666666")
        }.padding(2)
        .margin({top:5})
        .borderStyle(2)
        .backgroundColor("#f5f5f5f5")

        Row(){
          Text("编辑")
            .fontSize(14)
            .onClick(()=>{
              router.pushUrl({url:"pages/EventEditPage",params:{message:"edit",data:item}})
            })
          Text("删除")
            .fontSize(14)
            .margin({left:10})
            .fontColor(Color.Red)
            .onClick(()=>{
              this.deleteEvent(item.id)
            })
        }.padding(2)
        .margin({top:10})
        .borderStyle(2)
      }.alignItems(HorizontalAlign.Start)
      .height(80)
      .margin({left:10,top:5})
    }.height(100)
    .width("90%")
    .backgroundColor(Color.White)
    .borderStyle(5)
    .margin(10)
    .padding(10)
  }

  build() {
    Column(){
      Text("活动管理")
        .width("100%")
        .backgroundColor("#ff3551")
        .fontColor(Color.White)
        .fontSize(24)
        .textAlign(TextAlign.Center)
        .padding(10)
        .height("7%")
        .shadow(ShadowStyle.OUTER_DEFAULT_XS)

      List(){
        ForEach(this.eventList,(item:Event,index) =>{
          ListItem(){
            //单个event组件
            this.MyEvent(item,index)
          }
        })

        // ListItem(){
        //   Row(){
        //     Image('app.media.user_icon')
        //       .height(80)
        //       .width(140)
        //       .borderStyle(10)
        //
        //     Column(){
        //       Text("item.name")
        //         .maxLines(1) //只允许出现一行maxlines与textoverflow需要同时使用才能生效
        //         .textOverflow({overflow:TextOverflow.Clip}) //overflow:TextOverflow.Clip表示超过有的截取处理
        //         .constraintSize({maxWidth:120})
        //
        //       Row(){
        //         Image($r('app.media.event_position'))
        //           .width(14)
        //           .height(14)
        //         Text("item.addr")
        //           .maxLines(1) //只允许出现一行maxlines与textoverflow需要同时使用才能生效
        //           .textOverflow({overflow:TextOverflow.Ellipsis}) //overflow:TextOverflow.Ellipsis表示超过的用...处理
        //           .constraintSize({maxWidth:120})
        //           .fontSize(14)
        //           .fontColor("#666666")
        //       }.padding(2)
        //       .margin({top:5})
        //       .borderStyle(2)
        //       .backgroundColor("#f5f5f5f5")
        //
        //       Row(){
        //         Text("编辑")
        //           .fontSize(14)
        //           .onClick(()=>{
        //             //router.pushUrl({url:"pages/EventEditPage",params:{message:"edit",data:item}})
        //           })
        //         Text("删除")
        //           .fontSize(14)
        //           .margin({left:10})
        //           .fontColor(Color.Red)
        //           .onClick(()=>{
        //             // todo 删除事件的方法
        //           })
        //       }.padding(2)
        //       .margin({top:10})
        //       .borderStyle(2)
        //     }.alignItems(HorizontalAlign.Start)
        //     .height(80)
        //     .margin({left:10,top:5})
        //   }.height(100)
        //   .width("90%")
        //   .backgroundColor(Color.White)
        //   .borderStyle(5)
        //   .margin(10)
        //   .padding(10)
        //
        // }

      }.alignListItem(ListItemAlign.Center)
      .width("100%")
      .height('90%')
    }.backgroundColor("#f8f8f8")
    .width("100%")
    .height("100%")
  }
}