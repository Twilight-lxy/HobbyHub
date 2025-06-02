//前台的数据模型
// 用户类 包含ID 用户名 昵称，头像地址，住址，性别，经纬度信息
// addr 地址 string
// createTime string(date-time)
// gender 性别 string
// headImg 头像的路径 string
// id integer(int64)
// lat 维度 number(double)
// lon 经度 number(double)
// name 姓名 string
// password 密码 string
// username ⽤⼾名,登录的唯⼀凭证 string
export class User {
  id: number
  username: string
  name: string
  headImg: string
  addr: string
  gender: string
  lat: number
  lon: number

  constructor(
    id: number,
    username: string,
    name: string,
    headImg: string,
    addr: string,
    gender: string,
    lat: number,
    lon: number,
  ) {
    this.id = id
    this.username = username
    this.name = name
    this.headImg = headImg
    this.addr = addr
    this.gender = gender
    this.lat = lat
    this.lon = lon
  }
}