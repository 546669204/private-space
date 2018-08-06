import { app, BrowserWindow,ipcMain,Menu,dialog } from 'electron'
import fs from 'fs'
import path from 'path'
import os from 'os'
import child_process from 'child_process'
/**
 * Set `__static` path to static files in production
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-static-assets.html
 */
if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

var upass = null;

if(process.env.NODE_ENV!== 'development'){
  try {
    var ext = os.platform() == "win32" ? ".exe" : "";
    upass = child_process.spawn("./uPass" + ext);
  } catch (error) {
    alert("找不到服务端文件")
  }
  
}

let mainWindow
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`


var drag_info = {};

function createWindow () {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    height: 563,
    useContentSize: true,
    width: 1000
  })

  mainWindow.loadURL(winURL)

  mainWindow.on('closed', () => {
    mainWindow = null
  })

  const menu = Menu.buildFromTemplate(MenuTemplate)
  Menu.setApplicationMenu(menu)

  ipcMain.on('ondragstart', (event, filePath) => {
    var crypto=require('crypto');
    var md5=crypto.createHash("md5");
    md5.update("xiaoming" + Math.random()*10000000);
    var str=md5.digest('hex');
    drag_info.targe = filePath;
    drag_info.file_name = str.toUpperCase();
    drag_info.file_path = path.join(path.resolve("/"),drag_info.file_name);
    fs.appendFileSync(drag_info.file_path,"")
    event.sender.startDrag({
      file: drag_info.file_path,
      icon: ""//process.cwd() + '/src/main/file.png'
    })
    //wat.close();
    // fs.unlinkSync(drag_info.file_path)
  })




}

app.on('ready', createWindow)

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
    if(process.env.NODE_ENV!== 'development'){
      upass.kill();
    }
  }
})

app.on('activate', () => {
  if (mainWindow === null) {
    createWindow()
  }
})


const MenuTemplate = [
  {
    label: 'Edit',
    submenu: [
      {role: 'undo'},
      {role: 'redo'},
      {type: 'separator'},
      {role: 'cut'},
      {role: 'copy'},
      {role: 'paste'},
      {role: 'pasteandmatchstyle'},
      {role: 'delete'},
      {role: 'selectall'}
    ]
  },
  {
    label: 'View',
    submenu: [
      {role: 'reload'},
      {role: 'forcereload'},
      {role: 'toggledevtools'},
      {type: 'separator'},
      {role: 'resetzoom'},
      {role: 'zoomin'},
      {role: 'zoomout'},
      {type: 'separator'},
      {role: 'togglefullscreen'}
    ]
  },
  {
    role: 'window',
    submenu: [
      {role: 'minimize'},
      {role: 'close'}
    ]
  },
  {
    role: 'help',
    submenu: [
      {
        label: 'Learn More',
        click () { 
          dialog.showMessageBox({
            type : "none",
            title:"关于我们",
            buttons :[],
            message :"图标来自于:https://www.flaticon.com/packs/file-types(如有侵权删)"
          })
          //require('electron').shell.openExternal('https://electronjs.org') 
        }
      }
    ]
  }
]


function fsExistsSync(path) {
  try{
  fs.accessSync(path,fs.F_OK);
  }catch(e){
  return false;
  }
  return true;
}

var drag_timer = null;
var wat =  fs.watch("/",{ persistent: true, recursive: true },(event, filename)=>{
  //if (event === 'rename') {

    //console.log(event,filename)
    
    filename = path.resolve(filename)
    if (filename != drag_info.file_path && drag_info.file_name == path.basename(filename)){
      console.log(filename,drag_info.file_path)
      clearTimeout(drag_timer);
      drag_timer = setTimeout(()=>{
        var obj = [drag_info.targe,path.dirname(filename)];
        drag_info = Object.assign({});
        fs.unlinkSync(filename);
        mainWindow.webContents.send('dragExprot', obj);

        console.log("选中的目录是",obj)
      },200);
    }
    
    // if(!fsExistsSync(filename)){
    //   console.log("add",filename)
    // }


  //}
})




