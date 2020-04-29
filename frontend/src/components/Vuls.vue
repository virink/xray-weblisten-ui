<template>
  <div class="vuls block">
    <div v-if="loading" class="loading">Loading...</div>
    <el-dialog
      title="Raw"
      :visible.sync="dialogVisible"
      @opened="vulRawDialogOpened"
      width="80%">
      <div id="editor"></div>
    </el-dialog>

    <el-card class="box-card" v-for="vul in vuls" :key="vul.id" >
      <div slot="header" class="clearfix">
        <span>{{vul.title}}</span>
        <el-button style="float: right; padding: 3px 0" @click="openRawDialog(vul)" type="text">Raw</el-button>
      </div>
      <el-row>
        <el-col :span="4"><div class="grid-title">URL</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.url}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">Type</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.type}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">Payload</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.payload}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">Params</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.aprmas}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">Plugin</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.plugin}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">VulnClass</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.vuln_class}}</div></el-col>
      </el-row>
      <el-row>
        <el-col :span="4"><div class="grid-title">Found</div></el-col>
        <el-col :span="20"><div class="grid-content">{{vul.found}}</div></el-col>
      </el-row>
    </el-card>

  </div>
</template>

<script>

import {
  fetchVuls
} from '@/service'
import JSONEditor from 'jsoneditor'
import 'jsoneditor/dist/jsoneditor.min.css'
import moment from 'moment';
import 'moment/locale/zh-cn'
moment.locale('zh-cn')

export default {
  name: 'Vuls',
  watch: {
    '$route': 'fetchData'
  },
  data(){
    return {
      timer: null,
      id: this.$route.params.id,
      loading: false,
      dialogVisible:false,
      currentRaw: null,
      error: null,
      editor:null,
      vuls:[]
    }
  },
  created() {
    this.fetchData()
    this.timer = setInterval(this.fetchData,10*1000)
  },
  beforeDestroy() {
    clearInterval(this.timer)
  },
  methods: {
    fetchData () {
      this.error = this.post = null
      this.loading = true
      fetchVuls(this.id).then(resp => {
        this.loading = false
        var vuls = []
        resp.data.code == 0 && resp.data.data.forEach(e=>{
          vuls.push({
            id: e.ID,
            title: e.title,
            url: e.url,
            params: e.params,
            payload: e.payload,
            plugin: e.plugin,
            type: e.type,
            found: moment(e.create_time).format('YYYY-MM-DD HH:mm:ss'),
            vuln_class: e.vuln_class,
            raw: e.raw,
          })
        })
        this.vuls = vuls
      })
    },
    openRawDialog(vul){
      try{
        let obj = JSON.parse(vul.raw)
        this.currentRaw = obj
        this.dialogVisible = true
      }catch(e){
        this.$message.error(e.toString());
      }
    },
    vulRawDialogOpened(){
      if(!this.editor){
        this.editor = new JSONEditor(document.getElementById("editor"),{})
      }
      this.editor.set(this.currentRaw)
    }
  },
}
</script>

<style scoped>

.box-card {
  width: 800px;
  padding: 10px;
  margin: 20px auto;
}

.grid-title {
  border-radius: 4px;
  min-height: 36px;
  font-size: 18px;
  text-align: left;
}
.grid-content {
  border-radius: 4px;
  min-height: 36px;
  text-align: left;
}
</style>
