import { it } from "@ohos/hypium"
import { loginUserAction } from "../utils/UserUtil"
import {Event} from '../viewmodel/Event'
import {User} from '../viewmodel/User'
import {getEventList} from '../api/EventApi'
import {EventComponent} from '../component/EventConponent'

class EvertType{
  name : string
  icon : Resource

  constructor(name:string,icon:Resource) {
    this.name = name
    this.icon = icon
  }
}

//轮播图定义 使用在下面
const typeList:EvertType[]=[
  new EvertType("附近",$r("app.media.ic_home_recommend")),
  new EvertType("运动",$r("app.media.ic_home_motion")),
  new EvertType("骑行",$r("app.media.ic_home_ride")),
  new EvertType("郊游",$r("app.media.ic_home_excursion")),
]




@Preview
@Component
export default struct HomeViewComponent{

  //活动信息数据
  @State eventList:Event[] = []

  //监控父组件显示状态 //父组件向子组件传递值
  @Prop @Watch('onPageVisibilityChange') pageVisible:boolean

  //监听目标函数
  async onPageVisibilityChange(){
    //获取活动列表
    this.getEvents()
  }

  //在组件加载时获取活动列表
  async getEvents(){
    //已登录
    loginUserAction(async (user:User) =>{
      let result = (await getEventList(user.id)).obj as Event[]
      this.eventList = []

      //通过构造方法添加，实现属性的数据监听
      result.forEach(item =>{
        let event = new Event(
          item.id,
          item.userId,
          item.name,
          item.intro,
          item.addr,
          item.headImg,
          item.startTime,
          item.lon,
          item.lat,
          item.state,
          item.updateTime
        )
        this.eventList.push(event)
      }) //通用 默认缺省
    },async () =>{
      let result = (await getEventList(0)).obj as Event[]
      this.eventList = []

      //通过构造方法添加，实现属性的数据监听
      result.forEach(item =>{
        let event = new Event(
          item.id,
          item.userId,
          item.name,
          item.intro,
          item.addr,
          item.headImg,
          item.startTime,
          item.lon,
          item.lat,
          item.state,
          item.updateTime
        )
        this.eventList.push(event)
      })
    })
  }

  build() {
    Column(){
      Search({placeholder:'搜索活动',value:''})
        .width("90%")
        .placeholderColor("#999999") //默认提示颜色
        .backgroundColor(Color.White)
        .margin({top:'10%'})
        .fontColor('#999999')
        .height('5%')
        .textAlign(TextAlign.Center)
        .onSubmit(value => {
          //todo 跳转到活动搜索页面
        })

      Swiper(){
        Image($r('app.media.ic_swiper_1'))
        Image($r('app.media.ic_swiper_2'))
        Image($r('app.media.ic_swiper_3'))
      }.width('95%')
      .height('25%')
      .margin({top:20})
      .borderRadius(20)//设置圆角数
      .autoPlay(true)//自动轮播
      .loop(true)//是否循环

      Row(){
        ForEach(typeList,(item:EvertType) =>{
          Column(){
            Image(item.icon)
              .width(40)
              .height(40)
            Text(item.name)
              .padding({top:5})
          }.layoutWeight(1) //边距1
          .onClick(()=>{
            //todo 各按钮的动作
          })
        })
      }
      .width('90%')
      .height('10%')
      .margin({top:10})

      Blank()

      WaterFlow(){
        ForEach(this.eventList,(item:Event)=>{
          //获取活动列表
          FlowItem(){
            //获取单个活动的组件
            EventComponent({item:item})
          }
        })
      }
      .layoutDirection(FlexDirection.Column)//布局的主轴方向为列方向
      .columnsTemplate("lfr lfr")
      .columnsGap(5)
      .rowsGap(5)
      .width('90%')
      .height('50%')
    }
    .width('100%')
    .height('100%')
    .backgroundImage($r('app.media.bg_home'))
    .backgroundImageSize({
      width:'100%',
      height:'20%'
    })
  }
}