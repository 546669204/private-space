<template>
<transition type="fade">
        <el-card class="box-card" v-show="show">
                <div v-for="item in tareaList" :key="item.ID">
                    <!-- <transition-group type="fade"> -->
                        <div>{{item.Name}}-{{item.Import?"加密":"解密"}}</div>
                        <el-progress :percentage="item.Percentage" :show-text="false"></el-progress>
                    <!-- </transition-group> -->
                    
                </div>
        </el-card>
</transition>

</template>
<script>
export default {
    name:"tarea",
    data(){
        return {
            willHide:[],
            show:false,
            hideTimer:null,asfarf:[],
            tareaList:[]
        }
    },
    methods:{
        progress(e){
            if(this.tareaList.map((v)=>v.ID).indexOf(e.data.ID) == -1){
                this.tareaList.push(e.data);
            }else{
                if (this.tareaList[this.tareaList.map((v)=>v.ID).indexOf(e.data.ID)].Percentage != e.data.Percentage){
                this.$set(this.tareaList, this.tareaList.map((v)=>v.ID).indexOf(e.data.ID), e.data)
                // this.tareaList[this.tareaList.map((v)=>v.ID).indexOf(e.data.ID)].Percentage = e.data.Percentage;
                }
            }
        }
    },
    created() {
        console.log(this.tareaList)
    },
    watch:{
        tareaList(cur,old){
            cur.forEach((v)=>{
                if(v.Percentage == 100){
                    if (this.willHide.indexOf(v.ID) == -1){
                        this.willHide.push(v.ID);
                        setTimeout(()=>{
                            this.tareaList.splice(this.tareaList.map((v2)=>v2.ID).indexOf(v.ID),1);
                        },3000)
                    }
                }
            })


            clearTimeout(this.hideTimer)
            if (this.tareaList.length == 0){
                this.hideTimer = setTimeout(()=>{this.show = false;},3000)
            }else{
                if (!this.show){
                    this.show = true;
                }
            }
        }
    }

}
</script>
<style lang="scss" scoped>
.box-card{
    position: absolute;
    left: 15px;
    bottom: 15px;
    width: 170px;
}
</style>
