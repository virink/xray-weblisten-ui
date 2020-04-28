<template>
  <el-container>
    <el-header>
      <el-menu router :default-active="this.$route.path" class="el-menu-demo" mode="horizontal">
        <el-menu-item index="/projects">项目管理</el-menu-item>
        <!--<el-menu-item index="/vuls">漏洞详情</el-menu-item>-->
      </el-menu>
    </el-header>
    <el-main>
      <transition name="fade-transform" mode="out-in">
        <router-view :key="this.$route.path" />
      </transition>
    </el-main>
  </el-container>
</template>

<script>
  export default {
    name: 'Layout'
  }
</script>

<style>
  .el-header {
    background-color: #B3C0D1;
    color: #333;
    text-align: center;
    line-height: 60px;
  }
  
  .el-main {
    background-color: #E9EEF3;
    color: #333;
    text-align: center;
  }
  
  body > .el-container {
    margin-bottom: 40px;
  }
  
</style>