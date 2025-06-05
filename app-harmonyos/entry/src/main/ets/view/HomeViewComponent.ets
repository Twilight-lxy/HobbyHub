

@Preview
@Component
export default struct HomeViewComponent{
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
    }
  }
}