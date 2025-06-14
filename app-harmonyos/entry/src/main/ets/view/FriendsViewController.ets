import { FriendsViewComponent } from './FriendsViewComponent'
import { FriendsLoginViewComponent } from './FriendsLoginViewComponent'

@Component
export struct FriendsViewController {
  @StorageProp("isLogin") isLogin: boolean = false

  build() {
    Column() {
      if (this.isLogin) {
        // 已登录状态下显示好友列表
        FriendsViewComponent()
      } else {
        // 未登录状态显示登录提示
        FriendsLoginViewComponent()
      }
    }
    .width('100%')
    .height('100%')
  }
} 