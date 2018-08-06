<template>
    <div>
        <template v-for="(item) in currentPath" >
            <span @click="handlerClickBreadcrumb(item)" :key="1000+item.ID">{{item.Name}}</span>
            <i class="el-icon-arrow-right" :key="item.ID"/>
        </template>
    </div>
</template>
<script>
export default {
    name:"breadcrumb",
    props:["currentPath"],
    data(){
        return {

        }
    },
    methods:{
        handlerClickBreadcrumb(e){
            this.$emit("changePath",e)
        },
    },
}
</script>
