//获取是否登录

import { User } from '../viewmodel/User'

//通过该方法获取用户是否登录
export function loginUserAction(loginAction: (User: User) => void = () => {
}, noLoginAction: () => void = () => {
}) {
  //从App的Storage获取当前登录状态
  let isLogin = AppStorage.get("isLogin") as boolean
  if (!isLogin) {
    noLoginAction();
  } else {
    let user = AppStorage.get("user") as User
    loginAction(user)
  }
}