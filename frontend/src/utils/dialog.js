export const dialogError = (dialog, content1, content2) => {
  dialog.error({
    title: '错误',
    content: content1 + '，' + content2,
    positiveText: '那好吧',
  })
}

export const dialogInfo = (dialog, content) => {
  dialog.info({
    title: '提示',
    content: content,
    positiveText: '确认',
  })
}