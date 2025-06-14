import http from '@ohos.net.http';
import { ServerPath } from "../common/Constans";
import { httpRequest, ResponseData } from '../utils/HttpRequest';

// 加入活动
export function joinEvent(eventId: number, userId: number): Promise<ResponseData> {
  let url = `${ServerPath.join}?eventId=${eventId}&userId=${userId}`;
  let option: http.HttpRequestOptions = {
    method: http.RequestMethod.POST
  };
  return httpRequest(url, option);
}

// 获取活动成员列表
export function getEventMembers(eventId: number): Promise<ResponseData> {
  let url = `${ServerPath.eventMemberList}?eventId=${eventId}`;
  let option: http.HttpRequestOptions = {
    method: http.RequestMethod.GET
  };
  return httpRequest(url, option);
}

// 退出活动
export function quitEvent(eventId: number, userId: number): Promise<ResponseData> {
  let url = `${ServerPath.join}/quit?eventId=${eventId}&userId=${userId}`;
  let option: http.HttpRequestOptions = {
    method: http.RequestMethod.DELETE
  };
  return httpRequest(url, option);
} 