import { register } from './global'
import packageInfo from '../../package.json'

export default {
  install: (app) => {
    register(app)
    console.log(`
       欢迎使用 pcfarm-admin
       当前版本:v${packageInfo.version}
       默认接口文档地址:http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
       默认前端运行地址:http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
    `)
  }
}
