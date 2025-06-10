//用户编辑界面

import { router } from '@kit.ArkUI'
import { User } from '../viewmodel/User'
import { updateUserInfo } from '../api/UserApi'
import { showText } from '../utils/ShowUtil'
import { BusinessError } from '@kit.BasicServicesKit'
import { BaseUrl } from '../common/Constans'
import { selectImage } from '../utils/UploadUtil'

@Entry
@Component
struct UserEditPage {
  @StorageProp("user") user: User = new User(0, "", "", "", "", "", 0, 0)
  //后端图片库中的预设图片
  @State imageList: string[] = [
    '/upload/ic_user1.jpg',
    '/upload/ic_user2.jpg',
    '/upload/ic_user3.jpg',
    '/upload/ic_user4.jpg',
    '/upload/ic_user5.jpg',
    '/upload/ic_user6.jpg',
  ]
  //是否显示头像
  @State selectShow: boolean = false

  //更新用户信息的方法
  updateUserInfo() {
    updateUserInfo(this.user).then(
      res => {
        if (res.success) {
          showText("更新成功")
          AppStorage.setOrCreate('user', this.user)
        } else {
          showText("更新失败")
        }
      })
      .catch((e: BusinessError) => {
        showText("网络异常")
      })
  }

  //头像选择组件
  @Builder
  ImageListComponent() {
    Column() {
      Text("选择头像")
        .width('90%')
        .height('10%')
      WaterFlow() {
        ForEach(this.imageList, (item: string) => {
          FlowItem() {
            Image(`${BaseUrl}${item}`)
              .width(80)
              .height(80)
              .onClick(() => {
                this.user.headImg = item
                this.selectShow = false
              })
          }
        })
      }
      .layoutDirection(FlexDirection.Column)
      .columnsTemplate('1fr 1fr 1fr 1fr')
      .columnsGap(5)
      .rowsGap(5)
      .width('90%')
      .height('100%')
    }.width('100%')
    .height('100%')
  }

  build() {
    Column() {
      Row() {
        Image($r('app.media.ic_return'))
          .width(30)
          .height(30)
          .onClick(() => {
            router.back()
          })
        Blank();
        Text("用户设置")
          .fontSize(18)
        Blank();
        Button("保存")
          .backgroundColor('#ff3551')
          .onClick(() => {
            this.updateUserInfo()
          })
      }.padding(15)
      .width('100%')
      .backgroundColor(Color.White)

      Row() {
        Image(`${BaseUrl}${this.user.headImg}`)
          .alt($r('app.media.ic_user_icon_default'))
          .width(50)
          .height(50)
          .borderRadius(35)
          .margin({ left: 20 })
          .onClick(() => {
            //此处代码为模拟器中，通过半模态窗口模拟头像选择
            this.selectShow = true
            /*此代码通过Picker选择器去选择头像图片
            selectImage(getContext(this) as common.UIAbilityContext)
              .then(res => {
                this.user.headImg = res
            })*/
          })
        Text("点击修改头像")
          .fontColor('#999999')
          .fontSize(16)
          .margin({ left: 10 })
      }.width('100%')
      .height('10%')
      .backgroundColor(Color.White)
      .bindSheet($$this.selectShow, this.ImageListComponent(), {
        height: 500,
        backgroundColor: Color.White
      })

      Column() {
        Text("个人信息")
          .fontColor('#999999')
          .fontSize(14)
        Row() {
          Text("昵称")
          TextInput({ text: this.user.name, placeholder: "请输入昵称" })
            .onChange(value => {
              this.user.name = value
            })
            .borderStyle(0)
            .backgroundColor(Color.White)
        }.padding({ top: 15, left: 15 })

        Row() {
          Text("地区")
          TextInput({ text: this.user.addr, placeholder: "请填写地区" })
            .onChange(value => {
              this.user.addr = value
            })
            .borderStyle(0)
            .backgroundColor(Color.White)
        }.padding({ top: 15, left: 15 })
      }
      .alignItems(HorizontalAlign.Start)
      .width('100%')
      .backgroundColor(Color.White)
      .padding(10)
      .margin({ top: 10 })
    }
  }
}