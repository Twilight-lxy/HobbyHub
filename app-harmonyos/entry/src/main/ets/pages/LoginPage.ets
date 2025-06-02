
import { login } from '../api/LoginApi'
import {User} from '../viewmodel/User'
import { router } from '@kit.ArkUI';//arktx 路由
import { showText } from '../utils/ShowUtil';

@Entry
@Component
struct LoginPage{
  @State username : string = ''
  @State password : string = ''

  //登录的异步方法
  async clickLogin(){
    //发送登录请求
    let result = await login(this.username,this.password);
    //登录成功
    if (result.success) {
      try {
      //解析后台相应数据
      let data = result.obj as User
      let user = new User(data.id,data.username,data.name,data.headImg,data.addr,data.gender,data.lat,data.lon)


      //将当前用户信息写入appStorage
      AppStorage.setOrCreate("user",user)//创建或更新
      AppStorage.setOrCreate("isLogin",true)

      //跳转到主页面
      router.replaceUrl({url:'page/Index'})
      }catch (e){
        console.log(e.message)
      }
    }else{
      //登陆失败
      showText(result.info)
    }
  }

  build() { //构建页面
    Column(){
      Row(){
        Blank() //空行
        Text('随便看看')
          .fontSize(16)
          .fontColor('#999999')
          //todo 不登陆跳转到个人首页
      }.height('10%')
      .width('90%')
      .margin({top:10})

      Text("欢迎来到兴趣小队")
        .width("80%")
        .fontSize(22)
        .fontWeight(FontWeight.Medium)
        .padding(10)

      Column(){
        TextInput({placeholder:'用户名', text:'' })
          .onChange((value:string)=>{
            this.username = value
          })
          .shadow(ShadowStyle.OUTER_DEFAULT_SM)
          .backgroundColor('#ffffff')
          .height(45)

        Blank()
        TextInput({placeholder:'密码', text:'' })
          .onChange((value:string)=>{
            this.password = value
          })
          .type(InputType.Password)
          .shadow(ShadowStyle.OUTER_DEFAULT_SM)
          .backgroundColor('#ffffff')
          .height(45)
      }.height('18%')
      .width('90%')
      .padding(10)//行间距
      .borderRadius(10)//文本框边缘设置为原型

      Button("登录")
        .width('75%')
        .height(40)
        .margin({top:12})//上边距
        .borderRadius(25)//设置圆角
        .backgroundColor('#ff3351')
        .onClick(()=>{ //单击事件
          //todo 判断用户名和密码是否为空
          //非空的话，发起登录请求
          if (this.username != '' && this.password != '') {
            this.clickLogin()
          }
        })

      Row(){
        Text('去注册->')
          .height('100%')
          .fontColor('#666666')
          .onClick(()=>{
            //todo
          })
        Blank()
        Text('忘记密码')
          .height("100%")
          .fontColor('#666666')
          .onClick(()=>{
            //todo
          })
      }.height('6%')
      .width('80%')
      .margin({top:10})
    }.width('100%')
    .backgroundImage($r('app.media.ic_bg_login'))
    .backgroundImageSize({
      height:'20%',
      width:'100%'
    })
  }
}