<template>
  <div class="flex-box">
    <el-card class="box-card">
      <el-form ref="form" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="Username"></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input type="password" v-model="password" placeholder="Password"></el-input>
        </el-form-item>
        <div style="text-align:center">
          <el-button type="primary" @click="onSubmit">登录</el-button>

        </div>
      </el-form>
    </el-card>
  </div>
</template>
<script>
  export default {
    name: "Login",
    data() {
      return {
        username: "",
        password: "",
      }
    },
    methods: {
      onSubmit() {
        this.$webapi.webSocketSend("login",{user:this.username,pass:this.password},(res)=>{
          if(res.code == 0){
            this.$message({
              message: '登录成功',
              type: 'success'
            });
            this.$router.push("Index")
          }else{
            this.$message({
              message: res.msg,
              type: 'fail'
            });
          }
        })

      }
    }
  }
</script>

<style scoped>
  html {
    padding: 0;
  }

  .flex-box {
    display: flex;
    height: 100vh;
    align-content: center;
    justify-content: center;
    align-items: center;
    background: #ccc;
  }

  .box-card {
    background: #fff;
  }
</style>