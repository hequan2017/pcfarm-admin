import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'pcfarm-admin',
  showViteLogo: true,
  keepAliveTabs: false,
  logs: []
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    console.log(greenText('> 欢迎使用 pcfarm-admin'))
    console.log(greenText(`> 当前版本:v${packageInfo.version}`))
    console.log(
      greenText(
        `> 默认接口文档地址:http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`
      )
    )
    console.log(
      greenText(`> 默认前端运行地址:http://127.0.0.1:${env.VITE_CLI_PORT}`)
    )
    console.log('\n')
  }
}

export default config
