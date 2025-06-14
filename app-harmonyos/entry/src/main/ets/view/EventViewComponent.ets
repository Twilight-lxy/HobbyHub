import { EventLoginViewComponent } from './EventLoginViewComponent'
import { MyEventsListComponent } from './MyEventsListComponent'

@Component
export struct EventViewComponent {
  @StorageProp("isLogin") isLogin: boolean = false

  build() {
    Column() {
      if (this.isLogin) {
        // 已登录状态下显示我的活动列表
        MyEventsListComponent()
      } else {
        // 未登录状态显示登录提示
        EventLoginViewComponent()
      }
    }
    .width('100%')
    .height('100%')
  }
} 