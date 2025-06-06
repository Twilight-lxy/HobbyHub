/*
{
  "addr": "string",
"createTime": "string",
"headImg": "string",
"id": 0,
"ifDelete": 0,
"intro": "string",
"lat": 0,
"lon": 0,
"name": "string",
"startTime": "string",
"state": 0,
"updateTime": "string",
"userId": 0
}
*
* 活动id 活动名 活动内容 地址 图片 创建时间 活动的经纬度及状态
*/
@Observed
export class Event{
  id:number
  userId:number
  name:string
  intro:string
  addr:string
  headImg:string
  startTime:string
  lon:number
  lat:number
  state:number
  updateTime:string

  constructor(id:number,userId:number,name:string,intro:string,addr:string,headImg:string,startTime:string,lon:number,lat:number,state:number,updateTime:string) {
    this.id = id
    this.userId = userId
    this.name = name
    this.intro = intro
    this.addr = addr
    this.headImg = headImg
    this.startTime = startTime
    this.lon = lon
    this.lat = lat
    this.state = state
    this.updateTime = updateTime
  }
}