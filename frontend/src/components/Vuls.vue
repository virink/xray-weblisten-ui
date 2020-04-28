<template>
  <div class="vuls">
    <div v-if="loading" class="loading">Loading...</div>
    <div v-if="error" class="error">{{ error }}</div>
    <el-select v-model="currentProject" clearable placeholder="请选择项目">
      <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value">
      </el-option>
    </el-select>
    <el-card class="box-card" v-for="vul in vuls" :key="vul.id" >
      <div slot="header" class="clearfix">
        <span>{{vul.name}}</span>
      </div>
      <div v-for="o in 4" :key="o" class="text item">
        {{'列表内容 ' + o }}
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Vuls',
  watch: {
    '$route': 'fetchData'
  },
  data(){
    return {
      id: this.$route.params.id,
      loading: false,
      error: null,
      currentProject: '',
      options: [{
        value: '选项1',
        label: '黄金糕'
      }],
      vuls:[
        {
          id:1,
          name:'aaa'
        },
        {
          id:2,
          name:'bbb'
        },
        {
          id:3,
          name:'ccc'
        }]
    }
  },
  created() {
    console.log("id",this.id)
    // this.fetchData()
  },
  methods: {
    fetchData () {
      this.error = this.post = null
      this.loading = true
      // replace getPost with your data fetching util / API wrapper
      // getPost(this.$route.params.id, (err, post) => {
      //   this.loading = false
      //   if (err) {
      //     this.error = err.toString()
      //   } else {
      //     this.post = post
      //   }
      // })
    }
  },
}
</script>

<style scoped>
.box-card {
  width: 480px;
  padding: 5px;
  margin: 10px;
}
</style>
