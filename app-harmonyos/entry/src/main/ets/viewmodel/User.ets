//用户类 包含：用户ID、用户名、昵称、头像的路径、住址、性别、经纬度信息
@Observed
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
    lon: number) {
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