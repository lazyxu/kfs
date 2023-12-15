import { is } from '@electron-toolkit/utils'
import remoteMain from '@electron/remote/main'
import { spawn } from 'child_process'
import { app, BrowserWindow, dialog, ipcMain, Menu, shell } from 'electron'
import reloader from 'electron-reloader'
import fs from 'fs'
import iconvLite from 'iconv-lite'
import path from 'path'
import { getProcesses } from './processManager'

if (!app.isPackaged) {
  reloader(module)
}

if (app.isPackaged) {
  let configFilename = 'kfs-config.json'
  let configPath = path.join(process.resourcesPath, configFilename)
  console.log('spawn kfs-electron', __dirname, path.join(__dirname, '../../resources/kfs-electron'))
  let child = spawn(path.join(__dirname, '../../resources/kfs-electron'), ['-p', '11234'], {
    shell: true
  })
  let regex = new RegExp(/^Websocket server listening at: .+:(\d+)\n$/)
  child.stdout.on('data', function (chunk) {
    let stdout = iconvLite.decode(chunk, 'cp936')
    console.log('kfs-electron stdout', stdout)
    let results = regex.exec(stdout)
    if (results && results[1]) {
      const port = results[1]
      console.log('port', results[1])
      const config = fs.readFileSync(configPath).toString()
      const json = JSON.parse(config)
      json.port = port
      fs.writeFileSync(configPath, JSON.stringify(json, undefined, 2), { flag: 'w+' })
    }
  })
  child.stderr.on('data', function (chunk) {
    console.log('kfs-electron stderr', iconvLite.decode(chunk, 'cp936'))
  })
}

const publicPath = 'electron-' + (app.isPackaged ? 'production' : 'development')
let mainWindow

app.setName('考拉云盘')

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1280,
    height: 800,
    // title: "考拉云盘",
    // titleBarStyle: 'hidden',
    // backgroundColor: '#FFF',
    webPreferences: {
      preload: path.join(__dirname, '../preload/index.js'),
      nodeIntegration: true,
      contextIsolation: false,
      nativeWindowOpen: true,
      remote: true,
      sandbox: false,
      nodeIntegrationInSubFrames: true // for subContent nodeIntegration Enable
      // webviewTag:true //for webView
    },
    // from out\main\index.js
    icon: path.join(__dirname, '../../src/renderer/icon512.png')
  })

  remoteMain.initialize()
  remoteMain.enable(mainWindow.webContents)

  global.mainWindow = mainWindow

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(path.join(__dirname, '../renderer/index.html'))
  }

  const isMac = process.platform === 'darwin'

  const template = [
    // { role: 'appMenu' }
    ...(isMac
      ? [
          {
            label: app.name,
            submenu: [
              { role: 'about' },
              { type: 'separator' },
              { role: 'services' },
              { type: 'separator' },
              { role: 'hide' },
              { role: 'hideOthers' },
              { role: 'unhide' },
              { type: 'separator' },
              { role: 'quit' }
            ]
          }
        ]
      : []),
    // { role: 'fileMenu' }
    {
      label: 'File',
      submenu: [isMac ? { role: 'close' } : { role: 'quit' }]
    },
    // { role: 'editMenu' }
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
        ...(isMac
          ? [
              { role: 'pasteAndMatchStyle' },
              { role: 'delete' },
              { role: 'selectAll' },
              { type: 'separator' },
              {
                label: 'Speech',
                submenu: [{ role: 'startSpeaking' }, { role: 'stopSpeaking' }]
              }
            ]
          : [{ role: 'delete' }, { type: 'separator' }, { role: 'selectAll' }])
      ]
    },
    // { role: 'viewMenu' }
    {
      label: 'View',
      submenu: [
        { role: 'reload' },
        { role: 'forceReload' },
        { role: 'toggleDevTools' },
        { type: 'separator' },
        { role: 'resetZoom' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { type: 'separator' },
        { role: 'togglefullscreen' }
      ]
    },
    // { role: 'windowMenu' }
    {
      label: 'Window',
      submenu: [
        { role: 'minimize' },
        { role: 'zoom' },
        ...(isMac
          ? [{ type: 'separator' }, { role: 'front' }, { type: 'separator' }, { role: 'window' }]
          : [{ role: 'close' }])
      ]
    },
    {
      role: 'help',
      submenu: [
        {
          label: 'Learn More',
          click: async () => {
            await shell.openExternal('https://electronjs.org')
          }
        }
      ]
    }
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

app.whenReady().then(() => {
  createWindow()
})

app.on('window-all-closed', () => {
  getProcesses().forEach((p) => p.kill())
  app.quit()
})

ipcMain.on('select-dirs', async (event, arg) => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  })
  console.log('directories selected', result.filePaths)
})
