import {router} from '@kit.ArkUI'
import {Event} from '../viewmodel/Event'
import {addEvent} from '../api/EventApi'
import { showText } from '../utils/ShowUtil'
import {BusinessError} from "@kit.BasicServicesKit"
import {selectImage} from '../utils/UploadUtil'
import {common} from '@kit.AbilityKit';
import { BaseUrl } from '../common/Constans'

@Entry
@Component
struct EventEditPage{
  @State type:string='add'
  @State event : Event = new Event(0,0,"","","","","",120.39,36.30,0)
  @State showImage : boolean = false
  //是否显示选择组件
  @State selectShow: boolean= false
  //添加及活动信息
  addEvent(){
    let event = this.event
    //检查活动必填信息
    if (event.name != ''&&event.intro!='' && event.addr!="") {
      //发起后端创建活动请求
      addEvent(event).then(res => {
        if (res.success) {
          showText("添加成功",()=>{
            router.back()
          })
        }else{
          showText(res.info)
        }
      })
        .catch((e:BusinessError)=>{
          showText(e.message)
        })
    }else {
      showText('请填写完整活动信息')
    }
  }

  //头像选择组件


  build() {
    Stack(){
      Column(){
        Row(){
          Image($r(`app.media.ic_back`))
            .width(30)
            .height(30)
            .onClick(()=>{
              router.back()
            })
          Blank()
          Button(this.type == 'add'?'发布':'保存')
            .backgroundColor("#ff3551")
            .onClick(()=>{
              //todo 判断是否新增，新增则向后台发起添加活动
              if (this.type == 'add') {
                this.addEvent()
              }
            })
        }
        .padding(10)
        .width('100%')

        //标题
        TextInput({text:this.event.name,placeholder:'请填写活动标题'})
          .borderStyle(0)
          .backgroundColor(Color.White)
          .placeholderFont({size:20,weight:FontWeight.Medium})
          .onChange((val)=>{
            this.event.name = val
          })
        Divider().height(2) //分割线
        //活动内容
        TextArea({text:this.event.intro,placeholder:'请填写活动描述'})
          .width('100%')
          .height("30%")
          .borderStyle(0)
          .backgroundColor(Color.White)
          .onChange((value)=>{
            this.event.intro = value
          })
        Divider().height(2) //分割线
        //位置信息
        Row(){
          Image($r('app.media.event_position'))
            .width(25)
          Text('活动地点')
            .fontWeight(FontWeight.Medium)
            .width('25%')
            .margin({left:5})
          TextInput({text:this.event.addr,placeholder:'请填写活动地点'})
            .backgroundColor(Color.White)
            .textAlign(TextAlign.End)
            .width("70%")
            .onChange((value)=>{
              this.event.addr = value
            })
        }.padding(10)
        .width('100%')

        Divider().height(2)

        Text("活动图片")
          .fontWeight(FontWeight.Medium)
          .width('90%')
          .padding(10)

        Row(){
          if (this.event.headImg != ''){
            Image($r(`${BaseUrl}${this.event.headImg}`))
              .width(80)
              .height(80)
                .margin({right:20})
              .onClick(()=>{
                //因为是模拟器，无法调用Picker，通过半模态转场模拟活动图片
                this.showImage = true //状态值
                //通过Picker选择图片
                // selectImage((getContext(this)) as common.UIAbilityContext)
                //   .then(url =>{
                //     this.event.headImg = url
                //   })
                //   .then()
              })
          }
          Image($r('app.media.ic_img_add'))
            .width(80)
            .height(80)
            .onClick(()=>{
              //因为是模拟器，无法调用Picker，通过半模态转场模拟活动图片
              this.showImage = true //状态值
              //通过Picker选择图片
              // selectImage((getContext(this)) as common.UIAbilityContext)
              //   .then(url =>{
              //     this.event.headImg = url
              //   })
              //   .then()
            })
        }.padding(10)
        .width('100%')
        // .bindSheet($$this.selectShow,this.ImageListComponent(),{
        //   height:500,
        //   backgroundColor:Color.White
        // })
      }.width('100%')
      .height('100%')
    }

  }
}