import {router} from '@kit.ArkUI'
import {Event} from '../viewmodel/Event'
import {addEvent,editEvent} from '../api/EventApi'
import { showText } from '../utils/ShowUtil'
import {BusinessError} from "@kit.BasicServicesKit"
import {selectImage} from '../utils/UploadUtil'
import {common} from '@kit.AbilityKit';
import { BaseUrl } from '../common/Constans'
import {RouterParam} from '../viewmodel/RouterParam'
import { reqLocationPermissions } from '../utils/LocaltionUtil'
import { geoLocationManager } from '@kit.LocationKit'
import { JSON } from '@kit.ArkTS'
import { loginUserAction } from '../utils/UserUtil'
import { User } from '../viewmodel/User'
import { data } from '@kit.TelephonyKit'

@Entry
@Component
struct EventEditPage{
  @State type:string='add'
  @State event : Event = new Event(0,0,"","","","","",120.39,36.30,0,"",0,"")
  @State showImage : boolean = false
  @State selectedImageIndex: number = -1

  autoLocation:boolean = false
  @State imageList:Resource[] = [
    $r('app.media.ic_swiper_1'),
    $r('app.media.ic_swiper_2'),
    $r('app.media.ic_swiper_3'),
    $r('app.media.bg_eventinfo')
  ]

  //是否显示选择组件
  @State selectShow: boolean = false
  //页面加载时，默认获取位置的权限，如果获取到权限，就更新event的位置信息
  aboutToAppear() {
    //获取活动信息
    let param = router.getParams() as RouterParam<Event>
    this.type = param.message
    if (this.type == 'edit'){  //页面编辑功能
      this.event = param.data
      
      // 如果有已选图片，尝试根据图片路径找到对应的索引
      if (this.event.headImg && this.event.headImg.includes('event_image_')) {
        try {
          // 从路径中提取索引
          const match = this.event.headImg.match(/event_image_(\d+)/)
          if (match && match[1]) {
            const index = parseInt(match[1]) - 1
            if (index >= 0 && index < this.imageList.length) {
              this.selectedImageIndex = index
            }
          }
        } catch (error) {
          console.error('解析图片路径出错:', error)
        }
      }
    }else if (this.type == 'add') {
      //初始化活动定位
      reqLocationPermissions(getContext(this) as common.UIAbilityContext)
        .then(res =>{
          if (res) {
            //获取当前位置信息
            let requestInfo:geoLocationManager.LocationRequest={
              "scenario": geoLocationManager.LocationRequestScenario.NAVIGATION,
              'timeInterval':0,
              'distanceInterval':0,
              "maxAccuracy":0
            }
            geoLocationManager.getCurrentLocation(requestInfo)
              .then(location =>{
                this.autoLocation = true
                this.event.lat = location.latitude
                this.event.lon = location.longitude

                //如果地理编码可用，则通过地理编码获取当前位置
                if (geoLocationManager.isGeocoderAvailable()) {
                  let reverseGeocodeRequest:geoLocationManager.ReverseGeoCodeRequest = {
                    latitude:location.latitude,
                    longitude:location.longitude,
                    maxItems:1
                  }
                  geoLocationManager.getAddressesFromLocation(reverseGeocodeRequest,(err,data)=>{
                    if (err) {
                      console.log("getAddressFromLocation err:"+JSON.stringify(err))
                    }else {
                      this.event.addr = data[0].placeName as string
                    }
                  })
                }
              })
          }else {
            showText("获取位置信息失败，请手动填写活动位置")
          }
        }).catch(()=>{
        showText("获取位置信息失败，请手动填写活动位置")
      })
    }
    //页面加载时，初始化当前用户id
    loginUserAction((user:User)=>{
      this.event.userId = user.id
    })
  }
  //添加及活动信息
  addEvent(){
    let event = this.event
    // 如果选择了图片，设置图片路径
    if (this.selectedImageIndex >= 0 && this.selectedImageIndex < this.imageList.length) {
      // 这里我们使用索引作为图片标识，实际应用中可能需要更复杂的处理
      event.headImg = `/upload/event_image_${this.selectedImageIndex + 1}.jpg`
    }
    
    //检查活动必填信息
    if (event.name != ''&&event.intro!='' && event.addr!="") {
      //发起后端创建活动请求
      addEvent(event).then(res => {
        if (res.success == true) {
          showText("添加成功",()=>{
            // 返回首页，使用全局AppStorage存储刷新标志
            globalThis.needRefreshHome = true;
            router.back();
          })
        }else{
          showText(res.errorMessage)
        }
      })
        .catch((e:BusinessError)=>{
          showText(e.message)
        })
    }else {
      showText('请填写完整活动信息')
    }
  }

  //编辑方法
  editEvent(){
    let event = this.event
    // 如果选择了图片，设置图片路径
    if (this.selectedImageIndex >= 0 && this.selectedImageIndex < this.imageList.length) {
      // 这里我们使用索引作为图片标识，实际应用中可能需要更复杂的处理
      event.headImg = `/upload/event_image_${this.selectedImageIndex + 1}.jpg`
    }
    
    //检查活动必填信息
    if (event.name != ''&&event.intro!='' && event.addr!="") {
      //发起后端创建活动请求
      editEvent(event).then(res => {
        if (res.success == true) {
          showText("修改成功",()=>{
            // 返回首页，使用全局AppStorage存储刷新标志
            globalThis.needRefreshHome = true;
            router.back();
          })
        }else{
          showText(res.errorMessage)
        }
      })
        .catch((e:BusinessError)=>{
          showText(e.message)
        })
    }else {
      showText('请填写完整活动信息')
    }
  }

  //活动预制图片选择组件
  @Builder
  ImageListComponent(){
    Column(){
      Row() {
        Text("选择活动图片")
          .fontSize(18)
          .fontWeight(FontWeight.Medium)
        Blank()
        Button("关闭")
          .backgroundColor(Color.Transparent)
          .fontColor("#ff3551")
          .onClick(() => {
            this.selectShow = false
          })
      }
      .width("90%")
      .padding(10)

      WaterFlow(){
        ForEach(this.imageList, (item: Resource, index: number) => {
          FlowItem(){
            Column() {
              Image(item)
                .width(160)
                .height(120)
                .objectFit(ImageFit.Cover)
                .borderRadius(8)
                .onClick(() => {
                  this.selectedImageIndex = index
                  this.event.headImg = index.toString() // 存储索引而非资源路径
                  this.selectShow = false
                })
              
              // 如果当前图片已选中，显示选中标记
              if (this.selectedImageIndex === index) {
                Row() {
                  Text("已选择")
                    .fontSize(14)
                    .fontColor("#ff3551")
                }
                .width("100%")
                .justifyContent(FlexAlign.Center)
                .margin({top: 5})
              }
            }
            .padding(5)
          }
        })
      }
      .layoutDirection(FlexDirection.Column)
      .columnsTemplate('1fr 1fr')
      .columnsGap(10)
      .rowsGap(10)
      .width("90%")
      .height("90%")
    }
    .width("100%")
    .height("100%")
    .backgroundColor(Color.White)
    .borderRadius({topLeft: 16, topRight: 16})
  }

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
              if (this.type == 'add') {
                this.addEvent()
              }
              if (this.type == 'edit') {
                this.editEvent()
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
          if (this.selectedImageIndex >= 0 && this.selectedImageIndex < this.imageList.length){
            Image(this.imageList[this.selectedImageIndex])
              .width(80)
              .height(80)
              .margin({right:20})
              .borderRadius(8)
              .onClick(()=>{
                //点击已选图片，显示大图预览
                this.showImage = true
              })
          }
          Image($r('app.media.ic_img_add'))
            .width(80)
            .height(80)
            .onClick(()=>{
              //显示底部图片选择面板
              this.selectShow = true
            })
        }.padding(10)
        .width('100%')
        .bindSheet(this.selectShow, this.ImageListComponent(), {
          height: '50%',
          backgroundColor: Color.White
        })
      }.width('100%')
      .height('100%')

      // 图片预览层
      if (this.showImage && this.selectedImageIndex >= 0 && this.selectedImageIndex < this.imageList.length){
        Column(){
          // 添加关闭按钮
          Row() {
            Blank()
            Button() {
              Image($r('app.media.ic_back'))
                .width(24)
                .height(24)
            }
            .backgroundColor(Color.Transparent)
            .width(40)
            .height(40)
            .margin(10)
            .onClick(() => {
              this.showImage = false
            })
          }
          .width('100%')
          
          // 显示大图
          Image(this.imageList[this.selectedImageIndex])
            .width("90%")
            .objectFit(ImageFit.Contain)
        }
        .backgroundColor(Color.Black)
        .opacity(0.9)
        .width('100%')
        .height("100%")
        .justifyContent(FlexAlign.Center)
        .onClick(()=>{
          this.showImage = false
        })
      }
    }

  }
}