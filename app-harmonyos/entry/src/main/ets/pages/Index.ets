import HomeViewComponent from '../view/HomeViewComponent'
import { UserInfoViewComponent } from '../view/UserInfoViewComponent';

Entry
@Component
struct Index {
  //tab Icon 的当前焦点
  @State currentIndex: number = 0;
  //是否可见（用于判断是否登录，如果登录可见）
  @State pageVisible : boolean = false;

  //title 显示的内容
  //selectIcon 获得焦点时显示的图
  //unSelectIcon 失去焦点时显示的图片
  //index tab的索引
  @Builder//动态构造组件
  EndBar(title:string,selectIcon:Resource,unSelectIcon:Resource,index:number){
    Column(){
      Image(this.currentIndex == index ? selectIcon:unSelectIcon) //判断当前索引值与Tab索引值比较
        .width(26)
        .height(26)
        .margin(5)

      Text(title)
        .fontColor(this.currentIndex == index ? `#ff3551`:`#000000`)
    }
  }

  build() {
    Row() {
      Column() {
        //Tab 组件，首页下面的App工具栏
        Tabs(){
          TabContent(){
            HomeViewComponent({pageVisible:this.pageVisible)
          }.tabBar(this.EndBar('首页',$r('app.media.ic_home_on'),$r('app.media.ic_home_off'),0))
          TabContent(){
            //todo 点击活动要跳转的页面
          }.tabBar(this.EndBar('活动',$r('app.media.ic_event_on'),$r('app.media.ic_event_off'),1))
          TabContent(){
            //todo 点击好友要跳转的页面
          }.tabBar(this.EndBar('好友',$r('app.media.ic_friends_on'),$r('app.media.ic_friends_off'),2))
          TabContent(){
            UserInfoViewComponent()
          }.tabBar(this.EndBar('我的',$r('app.media.ic_mine_on'),$r('app.media.ic_mine_off'),3))
        }.barPosition(BarPosition.End) //工具栏放置的位置，页面下方
        .barHeight('10%')//工具栏的高度
        .width('100%')
        .height('100%')
        .onChange((index:number)=>{
          this.currentIndex = index
        })
      }
      .width('100%')
      .height('100%')
    }
    .height('100%')
  }
}