//将提示信息在app前端显示
export function showText(text: string, action: () => void = () => {
}) {
  AlertDialog.show({
    message: text, //提示内容
    alignment: DialogAlignment.Center, //显示位置
    primaryButton: {
      //确认按钮
      value: "确定",
      action: action//确认按钮的回调函数
    }
  })
}