#!/usr/bin/env node

'use strict'

if (!process.env.NODE_ENV) {
  process.env.NODE_ENV = 'production'
}

function merge(target, src, prefix = '') {
  Object.keys(target).forEach(k => {
    if (typeof src[`${prefix}${k}`] !== 'undefined')
      target[k] = src[`${prefix}${k}`]
  })
}

const defaultCfg = {
  // For server
  DB_USER: '',
  DB_PASSWD: '',
  DB_HOST: '127.0.0.1',
  DB_NAME: 'mirror',
  DB_PORT: 27017,
  API_PORT: 9999,
  DOCKERD_PORT: 2375,
  DOCKERD_HOST: '127.0.0.1',
  DOCKERD_SOCKET: '/var/run/docker.sock',
  BIND_ADDRESS: '',
  CT_LABEL: 'syncing',
  CT_NAME_PREFIX: 'syncing',
  LOGDIR_ROOT: '/var/log/ustcmirror',
  IMAGES_UPGRADE_INTERVAL: '1 * * * *',
  OWNER: `${process.getuid()}:${process.getgid()}`,

  // For client
  API_ROOT: '',
}

defaultCfg.API_ROOT = `http://localhost:${defaultCfg.API_PORT}/`
const path = require('path')
const fps = ['/etc/ustcmirror/config', path.join(process.env['HOME'], '.ustcmirror/config')]

for (const fp of fps) {
  let cfg
  try {
    cfg = require(fp)
  } catch (e) {
    if (e.code !== 'MODULE_NOT_FOUND') {
      throw e
    }
    continue
  }
  merge(defaultCfg, cfg)
}

merge(defaultCfg, process.env, 'YUKI_')

defaultCfg.isDev = process.env.NODE_ENV.startsWith('dev')
defaultCfg.isProd = process.env.NODE_ENV.startsWith('prod')
defaultCfg.isTest = process.env.NODE_ENV.startsWith('test')

if (!(defaultCfg.isTest ||
    process.argv[2] !== 'daemon' ||
    defaultCfg.BIND_ADDRESS))
{
  console.error('Need to specify <BIND_ADDRESS> in configuration')
  process.exit(1)
}

// should be lower case
defaultCfg['TOKEN_NAME'] = 'x-mirror-token'

defaultCfg._images = [
  'ustcmirror/gitsync:latest',
  'ustcmirror/rsync:latest',
  'ustcmirror/lftpsync:latest',
]

if (defaultCfg.isDev) {
  console.log('Configuration:', JSON.stringify(defaultCfg, null, 4))
}

module.exports = defaultCfg
