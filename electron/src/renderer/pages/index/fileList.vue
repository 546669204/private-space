<template>
  <div style="min-height: 100%;" ref="full">
    <div ref="datatable">
    <el-table ref="table" :data="ListData" style="width: 100%" @row-contextmenu="rightMenu" @row-dblclick="dbClick" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55">
      </el-table-column>
      <el-table-column prop="Name" label="名称" width="180">
        <template slot-scope="scope">
          <file-icon :icon="scope.row.Name" :isdir="scope.row.IsDir" />
          <span style="margin-left: 10px">{{ scope.row.Name }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="update_time" label="上次修改时间" width="180">
      </el-table-column>
      <el-table-column prop="Size" label="大小">
      </el-table-column>
    </el-table>
    </div>

    <div  @mousedown.right="rightMenu" ref="black">

    </div>

  </div>
</template>
<script>
import fileIcon from './FileIcon'
  export default {
    name: "Index",
    props:["ListData"],
    components:{fileIcon},
    data() {
      return {
        multipleSelection:[]
      }
    },
    methods: {
      rightMenu(e) {
        this.$emit("itemRightMenu",e)
      },
      dbClick(e){
        if (e.IsDir){
          this.$emit("itemDbClick",e)
        }
      },
      handleSelectionChange(val){
         this.multipleSelection = val;
         this.$emit("update:mulitItemList",val);
      },
      hana(e){
        console.log(e)
      }
    },
    mounted() {
      this.$refs.black.style.height ="300px";
    },
  }
</script>

<style scoped>
</style>