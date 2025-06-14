import { User } from '../viewmodel/User'
import { BaseUrl } from '../common/Constans'
import { router } from '@kit.ArkUI'

// 定义好友信息接口
interface Friend {
  id: number;
  username: string;
  name: string;
  headImg: string;
  addr: string;
  gender: string;
  lat: number;
  lon: number;
}

@Component
export struct FriendsViewComponent {
  @StorageProp("isLogin") isLogin: boolean = false
  @StorageProp("user") user: User = new User(0, "", "", "", "", "", 0, 0)
  @State friendsList: Friend[] = []

  aboutToAppear() {
    this.loadFriends()
  }

  async loadFriends() {
    // 这里应该调用获取好友列表的API
    // 由于可能没有实现该API，我们先使用模拟数据
    this.loadMockFriends()
  }

  loadMockFriends() {
    this.friendsList = [
      {
        id: 1,
        username: "失忆",
        name: "失忆",
        headImg: "user_avatar_1.jpg",
        addr: "山东省青岛市黄岛区",
        gender: "男",
        lat: 36.06,
        lon: 120.38
      },
      {
        id: 2,
        username: "qst022",
        name: "qst022",
        headImg: "user_avatar_2.jpg",
        addr: "山东省济南市历城区",
        gender: "男",
        lat: 36.67,
        lon: 117.12
      }
    ]
  }

  // 获取头像资源
  private getAvatarResource(path: string): Resource | string {
    if (!path) return $r('app.media.ic_user_icon_default')
    
    // 如果是user_avatar_格式的路径，转换为本地资源
    if (path.includes('user_avatar_')) {
      const match = path.match(/user_avatar_(\d+)/)
      if (match && match[1]) {
        const index = parseInt(match[1])
        switch (index) {
          case 1: return $r('app.media.ic_user_icon_default')
          case 2: return $r('app.media.ic_user_icon_default')
          default: return $r('app.media.ic_user_icon_default')
        }
      }
    }
    
    // 其他情况返回完整的远程URL路径
    return `${BaseUrl}${path}`
  }

  build() {
    Column() {
      // 标题栏
      Row() {
        Image($r('app.media.ic_friends_on'))
          .width(30)
          .height(30)
        Text("好友")
          .fontSize(20)
          .fontWeight(FontWeight.Medium)
          .margin({ left: 10 })
        Blank()
        Image($r('app.media.user_icon'))
          .width(24)
          .height(24)
      }
      .width('100%')
      .padding({ left: 20, right: 20, top: 15, bottom: 15 })
      .alignItems(VerticalAlign.Center)
      .backgroundColor(Color.White)
      
      // 好友列表
      List() {
        ForEach(this.friendsList, (friend: Friend) => {
          ListItem() {
            Row() {
              // 好友头像
              Image(this.getAvatarResource(friend.headImg))
                .width(60)
                .height(60)
                .borderRadius(30)
                .objectFit(ImageFit.Cover)
              
              // 好友信息
              Column() {
                // 好友名称
                Text(friend.name)
                  .fontSize(16)
                  .fontWeight(FontWeight.Medium)
                  .width('100%')
                  .textAlign(TextAlign.Start)
                
                // 好友地址
                Row() {
                  Image($r('app.media.event_position'))
                    .width(14)
                    .height(14)
                  Text(friend.addr)
                    .fontSize(14)
                    .fontColor('#666666')
                    .margin({ left: 4 })
                    .maxLines(1)
                    .textOverflow({ overflow: TextOverflow.Ellipsis })
                }
                .margin({ top: 8 })
                .width('100%')
                .alignItems(VerticalAlign.Center)
              }
              .layoutWeight(1)
              .margin({ left: 12 })
              .alignItems(HorizontalAlign.Start)
              .justifyContent(FlexAlign.Start)
              
              // 操作按钮
              Row() {
                Column() {
                  Image($r('app.media.ic_chat'))
                    .width(24)
                    .height(24)
                  Text("聊天")
                    .fontSize(12)
                    .fontColor('#666666')
                    .margin({ top: 4 })
                }
                .onClick(() => {
                  // 聊天功能
                })
                .margin({ right: 20 })
                
                Column() {
                  Image($r('app.media.ic_phone'))
                    .width(24)
                    .height(24)
                  Text("电话")
                    .fontSize(12)
                    .fontColor('#666666')
                    .margin({ top: 4 })
                }
                .onClick(() => {
                  // 电话功能
                })
                .margin({ right: 20 })
                
                Column() {
                  Image($r('app.media.ic_input_captcha'))
                    .width(24)
                    .height(24)
                  Text("邮箱")
                    .fontSize(12)
                    .fontColor('#666666')
                    .margin({ top: 4 })
                }
                .onClick(() => {
                  // 邮箱功能
                })
              }
              .alignItems(VerticalAlign.Center)
            }
            .width('100%')
            .padding(12)
            .backgroundColor(Color.White)
            .borderRadius(8)
          }
          .margin({ bottom: 10 })
        })
      }
      .width('100%')
      .layoutWeight(1)
      .padding({ left: 16, right: 16 })
      
      // 底部导航栏占位
      Blank().height('10%')
    }
    .width('100%')
    .height('100%')
    .backgroundColor('#f5f5f5')
  }
} 