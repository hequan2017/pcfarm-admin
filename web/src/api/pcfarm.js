import service from '@/utils/request'

export const createPcfarmServer = (data) => {
  return service({ url: '/pcfarm/server/create', method: 'post', data })
}

export const getPcfarmServerList = (params) => {
  return service({ url: '/pcfarm/server/list', method: 'get', params })
}

export const updatePcfarmBootPolicy = (data) => {
  return service({ url: '/pcfarm/server/bootPolicy', method: 'put', data })
}

export const executePcfarmPowerAction = (data) => {
  return service({ url: '/pcfarm/server/powerAction', method: 'post', data })
}

export const createPcfarmIPPool = (data) => {
  return service({ url: '/pcfarm/ipPool/create', method: 'post', data })
}

export const getPcfarmIPPoolList = (params) => {
  return service({ url: '/pcfarm/ipPool/list', method: 'get', params })
}

export const refreshPcfarmPXE = () => {
  return service({ url: '/pcfarm/pxe/refresh', method: 'post' })
}

export const getPcfarmPXEStatus = () => {
  return service({ url: '/pcfarm/pxe/status', method: 'get' })
}
