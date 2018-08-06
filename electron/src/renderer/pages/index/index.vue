<template>
  <el-container style="height:100vh">
    <el-aside width="200px">
      <el-tree :data="fileDirTree" @node-click="handleNodeClick" :default-expand-all="true" :highlight-current="true"></el-tree>
    </el-aside>
    <el-container>
      <el-header height="80px">
        <div class="header_menu_but">
          <div>
            <a title="拖拽上传" @click="DragDialogVisible = true">
              <div><font-awesome-icon icon="upload" size="2x" /></div>
              <div>拖拽上传</div>
            </a>
            <a title="导入" @click="handlerfileImport">
              <div><font-awesome-icon icon="file-import" size="2x" /></div>
              <div>导入</div>
            </a>
            <a title="导出" @click="handlerfileExport">
              <div><font-awesome-icon icon="file-export" size="2x"  /></div>
              <div>导出</div>
            </a>
            <a title="删除" @click="handlerFileDelete">
              <div><font-awesome-icon icon="trash-alt" size="2x"  /></div>
              <div>删除</div>
            </a>

            <!-- <a title="复制">
              <font-awesome-icon icon="copy" size="2x"  />
            </a>
            <a title="压缩">
              <font-awesome-icon icon="file-archive" size="2x"  />
            </a> -->

            <a title="刷新" @click="getFolderList">
              <div><font-awesome-icon icon="sync-alt" size="2x"  /></div>
              <div>刷新</div>
            </a>
          </div>


          <div>
            <a title="列表" @click="switchFileShowTypeNum = 0">
              <div><font-awesome-icon icon="th-list" size="2x"  /></div>
              <div>列表</div>
            </a>
            <a title="图块" @click="switchFileShowTypeNum = 1">
              <div><font-awesome-icon icon="th-large" size="2x"  /></div>
              <div>图块</div>
            </a>
          </div>

          

        </div>


        



        <div>
          <breadcrumb :currentPath="currentPath" @changePath="changePath"/>
        </div>
      </el-header>
      <el-main >
        <div>
          <fileList v-if="switchFileShowTypeNum ==0" :ListData="folerFileList" @itemRightMenu="rightMenu" @itemDbClick="itemDbClick" :mulitItemList.sync="mulitItemList"/>
          <fileRow v-else :ListData="folerFileList" @itemRightMenu="rightMenu" @itemDbClick="itemDbClick" :mulitItemList.sync="mulitItemList"/>
        </div>


        <el-dialog
          title="提示"
          :visible.sync="DragDialogVisible"
          width="50%"
          center>
          <div ref="mymain" style="border:3px solid #fff;height:300px;text-align:center;line-height:300px;font-size:32px;" @drop.stop.prevent="fileDrage" @dragleave="handledragenter(2)" @dragover.stop.prevent @dragenter.stop.prevent="handledragenter(1)">
            拖拉文件到这里
          </div>
        </el-dialog>
      </el-main>
      <!-- <el-footer>Footer</el-footer> -->
    </el-container>
    <tarea ref="progress" />
  </el-container>
</template>
<script>

  const Electron = require('electron');
  const remote = Electron.remote;
  const Menu = remote.Menu;
  const MenuItem = remote.MenuItem;
  const Dialog = remote.dialog;
  
import fileList from './fileList'
import fileRow from './fileRow'
import breadcrumb from './breadcrumb'
import tarea from './tarea'

  export default {
    name: "Index",
    components:{fileList,fileRow,breadcrumb,tarea},
    data() {
      return {
        fileDirTree: [],
        currentPath:[{Name:"根目录",ID:1}],
        folerFileList:[],
        switchFileShowTypeNum:0,
        ws:null,
        wsOpen:false,
        wsID:1,
        wsCB:[],
        DragDialogVisible:false,
        mulitItemList:[],
        tareaList:[],
      }
    },
    methods: {
      handleNodeClick(e) {
        this.currentPath = e.path;
        this.getFolderFileList(this.getFather())
        console.log(e)
      },
      switchFileShowType(i){
        this.switchFileShowTypeNum = i;
      },
      getFolderList(){
        this.$webapi.webSocketSend("getAllFolder","",(res)=>{
          if(res.code == 0){
            var tree = res.data.map((v)=>{
              v.label = v.Name;
              v.children = [];
              return v;
            });
            var treeID = res.data.map((v)=>{
              return v.ID;
            });
            tree.forEach((v,i)=>{
              if(treeID.indexOf(v.FatherID) != -1 && i != treeID.indexOf(v.FatherID)){
                tree[treeID.indexOf(v.FatherID)].children.push(v);
              }
            })          
            this.fileDirTree = [tree[0]];
            // console.log(this.fileDirTree)
            var ddpath = [];
            var tree = this.fileDirTree;
            function digui(father) {
              father.children.forEach((v)=>{
                v.path = Object.assign([],father.path);
                v.path.push(v);
                if (v.children.length!=0){
                  digui(v)
                }
              })
            }
            
              tree.forEach((v)=>{
                v.path = [v];
                digui(v)
              })
            // this.$message({
            //   message: "获取成功",
            //   type: 'success'
            // });
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
        this.getFolderFileList(this.getFather())
      },
      getFolderFileList(id){
        this.$webapi.webSocketSend("getFolderFile",{id},(res)=>{
          if(res.code == 0){
            this.folerFileList = res.data;
            // this.$message({
            //   message: "获取成功",
            //   type: 'success'
            // });
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      createFolder(name,father){
        this.$webapi.webSocketSend("createFolder",{name,father},(res)=>{
          if(res.code == 0){
            this.getFolderFileList(father);
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      deleteFolder(id){
        this.$webapi.webSocketSend("deleteFolder",{id},(res)=>{
          if(res.code == 0){
            this.getFolderFileList(this.currentPath[this.currentPath.length-1].ID);
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      deleteFolder(id){
        this.$webapi.webSocketSend("deleteFolder",{id},(res)=>{
          if(res.code == 0){
            this.getFolderFileList(this.currentPath[this.currentPath.length-1].ID);
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      fileImport(file){
        this.$webapi.webSocketSend("fileImport",{file,father:this.getFather()},(res)=>{
          if(res.code == 0){
            this.getFolderFileList(this.getFather());
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      fileExport(file,id){
        this.$webapi.webSocketSend("fileExport",{dstFile:file,file:id},(res)=>{
          if(res.code == 0){
            this.getFolderFileList(this.currentPath[this.currentPath.length-1].ID);
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      fileExportTempOpen(id){
        this.$webapi.webSocketSend("fileExportTemp",{file:id},(res)=>{
          if(res.code == 0){
            remote.shell.openItem(res.data.path);
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })
      },
      getFather(){
        return this.currentPath[this.currentPath.length-1].ID;
      },

      changePath(e){
        console.log(e,this.currentPath)
        this.currentPath.splice(this.currentPath.map((v)=>v.ID).indexOf(e.ID)+1,999)
        this.getFolderFileList(e.ID)
      },
      rightMenu(e) {
        this.menu = new Menu();
        this.menu.append(new MenuItem({
          label: "打开",
          click: ()=>{
            if(e.IsDir){
              this.currentPath.push(e);
              this.getFolderFileList(this.currentPath[this.currentPath.length-1].ID)
            }else{
              this.fileExportTempOpen(e.ID)
            }

          }
        }));
        this.menu.append(new MenuItem({
          label: "刷新",
          click: ()=>{
            this.getFolderFileList(this.currentPath[this.currentPath.length-1].ID)
          }
        }));
        this.menu.append(new MenuItem({
          label: "删除",
          click: ()=>{
            this.deleteFolder(e.ID);
          }
        }));
        this.menu.append(new MenuItem({
          label: "创建文件夹",
          click:()=>{
            this.$prompt('请输入文件夹名称', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              //inputPattern: /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
              //inputErrorMessage: '邮箱格式不正确'
            }).then(({ value }) => {
              this.createFolder(value,this.currentPath[this.currentPath.length-1].ID)
            }).catch(() => {
              this.$message({
                type: 'info',
                message: '取消输入'
              });       
            });
          }
        }));
       this.menu.append(new MenuItem({
          label: "资源管理器打开",
          click: ()=>{
            remote.shell.showItemInFolder("C:\\Users\\Administrator\\Desktop\\golang\\uPass")
          }
        }));
        

        this.menu.popup(remote.getCurrentWindow());
      },
      itemDbClick(e){
        console.log(e)
        this.currentPath.push(e);
        this.getFolderFileList(e.ID);
      },
      handlerfileImport(){
        Dialog.showOpenDialog(remote.getCurrentWindow(),{
          properties: ['openFile', 'multiSelections'],
          filters: [
            {name: 'All Files', extensions: ['*']},
          ]
        },(file)=>{
          this.fileImport(file);
        })
      },
      handlerfileExport(){
        if (this.mulitItemList.filter((v)=>!v.IsDir).map((v)=>v.ID).length<=0){
          this.$message({
            type:"info",
            message:"没有选择文件呢"
          })
          return
        }
        Dialog.showOpenDialog(remote.getCurrentWindow(),{
          properties: ['openDirectory'],
          filters: [
            {name: 'All Files', extensions: ['*']},
            {name: 'Images', extensions: ['jpg', 'png', 'gif']},
            {name: 'Movies', extensions: ['mkv', 'avi', 'mp4']},
            {name: 'Custom File Type', extensions: ['ming']},
          ]
        },(file)=>{
          this.fileExport(file,this.mulitItemList.filter((v)=>!v.IsDir).map((v)=>v.ID));
        })
      },
      handlerFileDelete(){
        this.mulitItemList.filter((v)=>!v.IsDir).forEach((v)=>{
          this.deleteFolder(v.ID);
        })
      },
      fileDrage(e,p){
        console.log(e,p)
        console.log(e.dataTransfer.files)
        this.handledragenter(2);
        this.DragDialogVisible = false;
        this.fileImport(Array.from(e.dataTransfer.files).map((v)=>v.path));
        //e.preventDefault()
        //ipcRenderer.send('ondragstart', '/path/to/item')
      },
      handledragenter(e){
        if (e == 1){
          this.$refs.mymain.className = "border-play";
        }else{
          this.$refs.mymain.className = "";
        }
      },
    },
    created(){
      this.$webapi.progress = (e)=>{
         this.$refs.progress.progress(e);
      }
      Electron.ipcRenderer.on('dragExprot', (event, message) => {
        this.fileExport(message[1],message[0]);
      })
    },
    mounted() {
      this.getFolderList();
    },
  }
</script>
<style lang="scss" scoped>
.el-main{
  height: 1px;
}
.border-play{
  animation:myborder 0.3s infinite;
}
@keyframes myborder
{
    from {border: 3px dashed #ccc;}
    to {border: 3px dotted #ccc;}
}

.header_menu_but{
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  text-align: center;
  &>div{
    display: flex;
    flex-direction: row;
  }
  & a>div{
    margin: 10px;

  }
}


</style>