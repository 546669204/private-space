<template>
  <div class="flex-box">
    <div style="width:100%;">
    <el-checkbox v-model="allcheck" :indeterminate="indeterminate" @change="handleAllCheck">全选</el-checkbox>

    </div>
    <div class="item" v-for="(item,index) in ListData" :key="index" :class="item.checked?'hover':''" @mousedown.right="rightMenu(item)" @dblclick="itemDbClick(item)">
      <file-icon :icon="item.Name" :isdir="item.IsDir" :height="40" draggable="true" @dragstart.native.prevent="handledrag(item)"/>
      <p>{{ item.Name }}</p>
      <el-checkbox v-model="item.checked" class="checkbox" @change="handleCheck"></el-checkbox>
    </div>
  </div>
</template>
<script>
import fileIcon from './FileIcon'
  const Electron = require('electron');
  const remote = Electron.remote;
  export default {
    name: "Index",
    props:["ListData"],
    components:{fileIcon},
    data() {
      return {
        menu:null,
        allcheck:false,
        indeterminate:false
      }
    },
    methods: {
      handleNodeClick(e) {
        e.preventDefault()
      },
      rightMenu(e) {
        this.$emit("itemRightMenu",e)
      },
      handledrag(e){
        //e.preventDefault()
        var shuzu = this.ListData.filter((v)=>v.checked && !v.IsDir).map((v)=>{v.ID});
        if (shuzu.length<=0){
          shuzu = [e.ID];
        }
       
        Electron.ipcRenderer.send('ondragstart', shuzu)
        
      },
      itemDbClick(e){
        if (e.IsDir){
          this.$emit("itemDbClick",e)
        }
      },
      handleCheck(e){
        var i = 0;
        var t = 0;
        for (const key in this.ListData) {
          if (this.ListData.hasOwnProperty(key)) {
            const element = this.ListData[key];
            i++;
            if(!element.checked){
              t++;
            }
          }
        }
        console.log(i,t)
        if(i == t){
          this.allcheck = false;
          this.indeterminate = false;
        }
        if(i > t){
          this.allcheck = true;
          this.indeterminate = true;
        }
        if(t == 0){
          this.allcheck = true;
          this.indeterminate = false;
        }
        console.log(this.allcheck,this.indeterminate)
      },
      handleAllCheck(e){
        for (const key in this.ListData) {
          if (this.ListData.hasOwnProperty(key)) {
            const element = this.ListData[key];
            element.checked = e
          }
        }
        this.indeterminate = false;
      }
    },
  }
</script>

<style lang="scss" scoped>
  .flex-box {
    display: flex;
    justify-content: flex-start;
    flex-wrap: wrap;
  }

  .checkbox {
    display: none;
  }
  p{
    word-break: break-all;
    margin-left: 0
  }
  .item {
    display: flex;
    width: 80px;
    flex-direction: column;
    position: relative;
    justify-content: flex-start;
    align-items: center;
    padding: 15px;
    &:hover,
    &.hover {
      background: #f2f2f2;
      &>.checkbox {
        display: block;
        position: absolute;
        right: 5px;
        top: 5px;
      }
    }
  }
</style>