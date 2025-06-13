//活动类 包含：活动ID、活动名、活动内容、地址、图片、创建时间、活动的经纬度及其状态
@Observed
export class Event {
  id: number
  userId: number
  name: string
  intro: string
  addr: string
  headImg: string
  starTime: string
  lat: number
  lon: number
  state: number
  createTime: string
  ifDelete: number
  updateTime: string

  constructor(
    id: number,
    userId: number,
    name: string,
    intro: string,
    addr: string,
    headImg: string,
    starTime: string,
    lat: number,
    lon: number,
    state: number,
    createTime:string,
    ifDelete:number,
    updateTime:string) {
    this.id = id
    this.userId = userId
    this.name = name
    this.intro = intro
    this.addr = addr
    this.headImg = headImg
    this.starTime = starTime
    this.lat = lat
    this.lon = lon
    this.state = state
    this.createTime = createTime
    this.ifDelete = ifDelete
    this.updateTime = updateTime
  }
}