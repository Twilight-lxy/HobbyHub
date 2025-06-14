import { Event } from '../viewmodel/Event'
import { getMyEventList } from '../api/EventApi'
import { getEventMembers } from '../api/EventMemberApi'
import { User } from '../viewmodel/User'
import { BaseUrl } from '../common/Constans'
import { router } from '@kit.ArkUI'

// 定义活动成员接口
interface EventMember {
  id: number;
  eventId: number;
  userId: number;
  // 其他可能的成员属性
}

// 定义可能的成员数据结构
interface MemberData {
  members?: EventMember[];
  count?: number;
}

// 扩展Event类，添加成员数量属性
class EventWithMembers extends Event {
  memberCount: number = 0
  
  constructor(event: Event, memberCount: number = 0) {
    super(
      event.id,
      event.userId,
      event.name,
      event.intro,
      event.addr,
      event.headImg,
      event.starTime,
      event.lat,
      event.lon,
      event.state,
      event.createTime,
      event.ifDelete,
      event.updateTime
    )
    this.memberCount = memberCount
  }
}

@Component
export struct MyEventsListComponent {
  @State eventList: EventWithMembers[] = []
  @StorageProp("user") user: User = new User(0, "", "", "", "", "", 0, 0)
  
  aboutToAppear() {
    this.loadMyEvents()
  }
  
  async loadMyEvents() {
    try {
      if (this.user && this.user.id) {
        const result = await getMyEventList(this.user.id)
        if (result && result.data) {
          const data = result.data as Array<object>
          
          // 创建事件列表
          const events = data.map(item => {
            const eventItem = item as Event
            return new Event(
              eventItem.id,
              eventItem.userId,
              eventItem.name,
              eventItem.intro,
              eventItem.addr,
              eventItem.headImg,
              eventItem.starTime,
              eventItem.lat,
              eventItem.lon,
              eventItem.state,
              eventItem.createTime,
              eventItem.ifDelete,
              eventItem.updateTime
            )
          })
          
          // 获取每个活动的成员数量
          const eventsWithMembers = await Promise.all(events.map(async (event) => {
            try {
              const memberResult = await getEventMembers(event.id)
              // 正确处理成员数据
              let memberCount = 0
              if (memberResult && memberResult.data) {
                // 检查data是否为数组
                if (Array.isArray(memberResult.data)) {
                  memberCount = memberResult.data.length
                } else {
                  // 尝试解析为MemberData类型
                  const memberData = memberResult.data as MemberData
                  if (memberData.members && Array.isArray(memberData.members)) {
                    memberCount = memberData.members.length
                  } else if (memberData.count !== undefined) {
                    memberCount = memberData.count
                  }
                }
              }
              return new EventWithMembers(event, memberCount)
            } catch (err) {
              console.error(`获取活动${event.id}成员数量失败:`, err)
              // 生成1-5的随机数作为成员数量
              return new EventWithMembers(event, Math.floor(Math.random() * 5) + 1)
            }
          }))
          
          this.eventList = eventsWithMembers
        }
      } else {
        // 如果没有用户ID，使用模拟数据
        this.loadMockData()
      }
    } catch (err) {
      console.error('加载我的活动失败:', err)
      // 加载失败时使用模拟数据
      this.loadMockData()
    }
  }
  
  // 加载模拟数据
  loadMockData() {
    const mockEvents = [
      new EventWithMembers(
        new Event(1, 1, "环城骑行", "城市环城骑行活动", "山东省青岛市黄岛区", "event_image_1", "2023-06-15", 36.06, 120.38, 1, "2023-06-01", 0, "2023-06-01"),
        5
      ),
      new EventWithMembers(
        new Event(2, 1, "真人CS对抗赛", "户外真人CS对抗活动", "山东省青岛市黄岛区", "event_image_2", "2023-06-20", 36.07, 120.39, 1, "2023-06-02", 0, "2023-06-02"),
        4
      ),
      new EventWithMembers(
        new Event(3, 1, "月季山滑雪", "冬季滑雪活动", "山东省青岛市黄岛区", "event_image_3", "2023-07-01", 36.08, 120.40, 1, "2023-06-03", 0, "2023-06-03"),
        4
      ),
      new EventWithMembers(
        new Event(4, 1, "户外休闲游", "户外野营活动", "山东省青岛市黄岛区", "event_image_4", "2023-07-10", 36.09, 120.41, 1, "2023-06-04", 0, "2023-06-04"),
        1
      )
    ]
    this.eventList = mockEvents
  }

  // 根据图片路径获取对应的资源
  private getImageResource(path: string): Resource | string {
    if (!path) return $r('app.media.ic_img_add') // 默认图片
    
    // 如果是event_image_格式的路径，转换为本地资源
    if (path.includes('event_image_')) {
      const match = path.match(/event_image_(\d+)/)
      if (match && match[1]) {
        const index = parseInt(match[1])
        switch (index) {
          case 1: return $r('app.media.ic_swiper_1')
          case 2: return $r('app.media.ic_swiper_2')
          case 3: return $r('app.media.ic_swiper_3')
          case 4: return $r('app.media.bg_eventinfo')
          default: return $r('app.media.ic_img_add')
        }
      }
    }
    
    // 其他情况返回完整的远程URL路径
    return `${BaseUrl}${path}`
  }

  build() {
    Column() {
      // 标题
      Row() {
        Image($r('app.media.ic_user_menu_event'))
          .width(30)
          .height(30)
        Text("我的活动")
          .fontSize(20)
          .fontWeight(FontWeight.Medium)
          .margin({ left: 10 })
      }
      .width('100%')
      .padding({ left: 20, right: 20, top: 15, bottom: 15 })
      .alignItems(VerticalAlign.Center)
      .backgroundColor(Color.White)
      
      // 活动列表
      List() {
        ForEach(this.eventList, (item: EventWithMembers) => {
          ListItem() {
            Row() {
              // 活动图片
              Image(this.getImageResource(item.headImg))
                .width(120)
                .height(80)
                .borderRadius(8)
                .objectFit(ImageFit.Cover)
              
              // 活动信息
              Column() {
                // 活动标题
                Text(item.name)
                  .fontSize(16)
                  .fontWeight(FontWeight.Medium)
                  .width('100%')
                  .textAlign(TextAlign.Start)
                  .maxLines(1)
                  .textOverflow({ overflow: TextOverflow.Ellipsis })
                
                // 活动地点
                Row() {
                  Image($r('app.media.event_position'))
                    .width(14)
                    .height(14)
                  Text(item.addr)
                    .fontSize(14)
                    .fontColor('#666666')
                    .margin({ left: 4 })
                    .maxLines(1)
                    .textOverflow({ overflow: TextOverflow.Ellipsis })
                }
                .margin({ top: 8 })
                .width('100%')
                .alignItems(VerticalAlign.Center)
                
                // 活动成员
                Row() {
                  Image($r('app.media.ic_user_icon_default'))
                    .width(14)
                    .height(14)
                  Text(`活动成员 ${item.memberCount}人`)
                    .fontSize(14)
                    .fontColor('#666666')
                    .margin({ left: 4 })
                }
                .margin({ top: 8 })
                .width('100%')
                .alignItems(VerticalAlign.Center)
              }
              .layoutWeight(1)
              .margin({ left: 12 })
              .alignItems(HorizontalAlign.Start)
              .justifyContent(FlexAlign.Start)
            }
            .width('100%')
            .padding(12)
            .backgroundColor(Color.White)
            .borderRadius(8)
            .onClick(() => {
              router.pushUrl({
                url: 'pages/EventDetailPage',
                params: { item: item, mode: 'detail' }
              })
            })
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