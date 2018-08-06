export default function install (Vue, options) {
    var webapi = {
        ws:null,
        wsOpen:false,
        wsID:1,
        wsCB:[],
        progress:null,
        init(){
            this.ws = new WebSocket("ws://127.0.0.1:31687");
            this.ws.onopen = ()=>{
              console.log("webSocket 连接成功")
              this.wsOpen = true;
              //this.getFolderList();
            }
            this.ws.onmessage = (e)=>{
              e = e.data.replace(/\n$/gi,"");
              e = JSON.parse(e);
              if (e.method == "return"){
                // console.log("<==",e.data,this.wsCB[e.id])
                if (typeof this.wsCB[e.id] == "function"){
                  this.wsCB[e.id](e.data);
                  this.wsCB[e.id] = undefined;
                }
              }
              if(e.method == "progress"){
                this.progress(e);
              }
              
            }
            this.ws.onclose = ()=>{
              console.log("WebSocket 被关闭");
              this.wsOpen = false;
            }
            this.ws.onerror = (e)=>{
              console.log("WebSocket 出错" + e);
              this.wsOpen = false;
            }
        },
        webSocketSend(method,data,cb){
            if(this.wsOpen){
              var temp_id = this.wsID++;
              this.wsCB[temp_id] = cb;
              if (typeof data == "object"){data = JSON.stringify(data)}
              this.ws.send(JSON.stringify({"id":temp_id,"method":method,"data":data}))
              // console.log("==>",JSON.stringify({"id":temp_id,"method":method,"data":data}))
            }
        },
    }
    webapi.init();

    Vue.prototype.$webapi = webapi;

};