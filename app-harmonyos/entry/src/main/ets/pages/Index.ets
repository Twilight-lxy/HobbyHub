//主页面框架

import HomeViewComponent from '../view/HomeViewComponent'
import { UserInfoViewComponent } from '../view/UserInfoViewComponent';
import { EventViewComponent } from '../view/EventViewComponent';
import { router } from '@kit.ArkUI';

@Entry
@Component
struct Index {
  //tab Icon 的当前焦点
  @State currentIndex: number = 0;
  //是否可见，用于判断是否登录，登录后可见
  @State pageVisible: boolean = false;
  
  // 添加页面显示生命周期方法
  onPageShow() {
    try {
      // 检查全局刷新标志
      if (globalThis.needRefreshHome === true) {
        // 重置标志
        globalThis.needRefreshHome = false;
        // 如果需要刷新，则切换pageVisible状态触发刷新
        this.pageVisible = !this.pageVisible;
        // 短暂延迟后再切换回来，确保触发两次变化
        setTimeout(() => {
          this.pageVisible = !this.pageVisible;
        }, 100);
      }
    } catch (error) {
      // 处理可能的异常
      console.error('刷新首页出错:', error);
    }
  }

  //动态构造组件
  @Builder
  EndBar(title: string, selectIcon: Resource, unSelectIcon: Resource, index: number) {
    //title:显示的名称;selectIcon:获得焦点时显示的图片;unselectIcon:失去焦点时显示的图片;index:tab的索引
    Column() {
      //将当前索引值与Tab的索引值进行比较，判断是否选中了该组件
      //图标组件:选中显示为on，未选中显示为off
      Image(this.currentIndex == index ? selectIcon : unSelectIcon)
        .width(26)
        .height(26)
        .margin(5)
      //文本组件:选中显示为红色，未选中显示为黑色
      Text(title)
        .fontColor(this.currentIndex == index ? '#ff3551' : '#000000')
    }
  }

  build() {
    Row() {
      Column() {
        //定义四个 Tab 组件
        //即首页下面的App工具栏，实现四个不同的页面索引
        Tabs() {
          TabContent() {
            HomeViewComponent({ pageVisible: this.pageVisible })
          }.tabBar(this.EndBar("首页", $r('app.media.ic_home_on'), $r('app.media.ic_home_off'), 0))

          TabContent() {
            EventViewComponent()
          }.tabBar(this.EndBar("活动", $r('app.media.ic_event_on'), $r('app.media.ic_event_off'), 1))

          TabContent() {
            //todo 点击好友跳转页面
          }.tabBar(this.EndBar("好友", $r('app.media.ic_friends_on'), $r('app.media.ic_friends_off'), 2))

          TabContent() {
            UserInfoViewComponent()
          }.tabBar(this.EndBar("我的", $r('app.media.ic_mine_on'), $r('app.media.ic_mine_off'), 3))
        }
        .barPosition(BarPosition.End) //工具栏放置的位置，页面下方
        .barHeight('10%') //工具栏高度
        .barWidth('100%') //工具栏宽度
        .height('100%')
        .onChange((index: number) => {
          this.currentIndex = index
        })
      }
      .width('100%')
      .height('100%')
    }
    .height('100%')
  }
}