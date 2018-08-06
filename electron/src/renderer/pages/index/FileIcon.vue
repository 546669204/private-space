<template >
    <div class="file" :class="className" v-if="!isdir" :style="{height:height +'px',width:height+'px'}" ></div>
    <div class="folder" v-else  :style="{height:height +'px',width:height+'px'}" ></div>
</template>
<script>
export default {
    name:"FIleIcon",
    props:["icon",'height',"isdir"],
    data(){
        return {
            className:""
        }
    },
    created() {
        if (this.icon.indexOf(".") == -1){
            this.className = this.icon;
            return
        }
        this.className = this.icon.slice(this.icon.lastIndexOf(".")+1);
    },
}
</script>
<style lang="scss" scoped >

.file{
    background-image: url("../../images/file.png");
    background-size: 100%;
    height: 15px;
    width: 15px;
    display: inline-block;
}
.folder{
    background-image: url("../../images/folder.png");
    background-size: 100%;
    height: 15px;
    width: 15px;
    display: inline-block;
}
.ai{
    background-image: url("../../images/ai.png")
}
.avi{
    background-image: url("../../images/avi.png")
}
.css{
    background-image: url("../../images/css.png")
}
.csv{
    background-image: url("../../images/csv.png")
}
.dbf{
    background-image: url("../../images/dbf.png")
}
.doc{
    background-image: url("../../images/doc.png")
}
.dwg{
    background-image: url("../../images/dwg.png")
}
.exe{
    background-image: url("../../images/exe.png")
}
.fla{
    background-image: url("../../images/fla.png")
}
.html{
    background-image: url("../../images/html.png")
}
.iso{
    background-image: url("../../images/iso.png")
}
.js{
    background-image: url("../../images/javascript.png")
}
.jpg{
    background-image: url("../../images/jpg.png")
}
.json{
    background-image: url("../../images/json-file.png")
}
.mp3{
    background-image: url("../../images/mp3.png")
}
.mp4{
    background-image: url("../../images/mp4.png")
}
.pdf{
    background-image: url("../../images/pdf.png")
}
.png{
    background-image: url("../../images/png.png")
}
.ppt{
    background-image: url("../../images/ppt.png")
}
.psd{
    background-image: url("../../images/psd.png")
}
.rtf{
    background-image: url("../../images/rtf.png")
}
.svg{
    background-image: url("../../images/svg.png")
}
.txt{
    background-image: url("../../images/txt.png")
}
.xls{
    background-image: url("../../images/xls.png")
}
.xml{
    background-image: url("../../images/xml.png")
}
.zip{
    background-image: url("../../images/zip.png")
}
</style>

