
//router (页面跳转) 传参的封装类
export class RouterParam<T>{
  data:T
  message:string
  constructor(data:T,message:string) {
    this.data = data
    this.message = message
  }
}