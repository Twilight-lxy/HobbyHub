import {router} from '@kit.ArkUI'
import {Event} from '../viewmodel/Event'
import {addEvent} from '../api/EventApi'
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

@Entry
@Component
struct EventEditPage{
  @State type:string='add'
  @State event : Event = new Event(0,0,"","","","","",120.39,36.30,0)
  @State showImage : boolean = false

  autoLocation:boolean = false
  @State imageList:string[] = [
    '/upload/ic_event1.jpg',
    '/upload/ic_event2.jpg',
    '/upload/ic_event3.jpg',
    '/upload/ic_event4.jpg',
  ]

  //是否显示选择组件
  @State selectShow: boolean= false
  //页面加载时，默认获取位置的权限，如果获取到权限，就更新event的位置信息
  aboutToAppear() {
    //获取活动信息
    let param = router.getParams() as RouterParam<Event>
    this.type = param.message
    if (this.type == 'add') {
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

  //活动预制图片选择组件
  @Builder
  ImageListComponent(){
    Column(){
      Text("选择活动图片")
        .width("90%")
        .height("10%")

      WaterFlow(){
        ForEach(this.imageList,(item:string) =>{
          FlowItem(){
            Image(`${BaseUrl}${item}`)
              .width(160)
              .height(80)
              .onClick(()=>{
                this.event.headImg = item
                this.showImage = false
              })
          }
        })
      }.layoutDirection(FlexDirection.Column)
      .columnsTemplate('1fr,1fr')
      .rowsGap(5)
      .width("90%")
      .height("100%")
    }
    .width("100%")
    .height("100%")
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
        .bindSheet($$this.selectShow,this.ImageListComponent(),{
          height:500,
          backgroundColor:Color.White
        })
      }.width('100%')
      .height('100%')

      if (this.showImage){
        Column(){
          Image(`${BaseUrl}${this.event.headImg}`)
            .width("100%")
            .objectFit(ImageFit.Cover)
        }.backgroundColor(Color.White)
        .width('100%')
        .height("100%")
        .onClick(()=>{
          this.showImage = false
        })
      }
    }

  }
}