<template>
  <div class="projects block" v-loading="loading">
    <el-button @click="createProjectDialogVisible = true">创建项目</el-button>
    <el-dialog
      title="创建项目"
      :visible.sync="createProjectDialogVisible"
      @open="createProjectDialogOpen"
      width="30%"
      center>
      <el-form :model="form">
        <el-form-item label="项目名称" label-width="120px">
          <el-input v-model="form.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="项目域名" label-width="120px">
          <el-input v-model="form.domain" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="启用插件" label-width="120px">
          <el-checkbox v-model="form.checkAll" @change="handleCheckAllChange">全选</el-checkbox>
          <el-checkbox-group v-model="form.plugins" @change="handleCheckedPluginsChange">
            <el-checkbox v-for="plugin in pluginOptions" :label="plugin" :key="plugin">{{plugin}}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="createProjectDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="createProjectDialogVisible = false">确 定</el-button>
      </span>
    </el-dialog>
    <el-table 
      ref="multipleTable"
      border
      :data="tableData" 
      :row-class-name="tableRowClassName"
      :default-sort = "{prop: 'date', order: 'descending'}"
      @selection-change="handleSelectionChange"
      style="width: 100%" >
      <el-table-column type="selection" fixed width="55"></el-table-column>
      <el-table-column prop="name" label="名称" fixed></el-table-column>
      <el-table-column prop="domain" label="域名"></el-table-column>
      <el-table-column prop="config" label="配置文件"></el-table-column>
      <el-table-column prop="listen" label="监听端口"></el-table-column>
      <el-table-column prop="process_id" label="进程ID"></el-table-column>
      <el-table-column label="操作" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="primary" @click="handleDetail(scope.$index, scope.row)">详情</el-button>
          <el-button size="mini" 
            :type="isRunning(scope.row)?'danger':'success'"
            @click="handleToggle(scope.$index, scope.row)">
            {{isRunning(scope.row)?'Stop':'Start'}}
          </el-button>
          <el-button size="mini" type="primary" @click="handleResult(scope.$index, scope.row)">结果</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination v-if="pager.total > pager.limit"
      :current-page="pager.current"
      :page-size="pager.limit"
      :total="pager.total"
      @current-change="handleCurrentChange"
      background
      layout="prev, pager, next">
    </el-pagination>
  </div>
</template>

<script>
const pluginOptions = [
  'xss', 'baseline', 'cmd_injection', 
  'crlf_injection','dirscan','jsonp',
  'path_traversal','redirect','sqldet',
  'ssrf','xxe','upload','brute_force',
  'struts','thinkphp','fastjson','phantasm'
];

import {
  fetchProjects,
  startProject,
  stopProject
} from '@/service'

export default {
  name: 'Projects',
  data(){
    return {
      pluginOptions:pluginOptions,
      createProjectDialogVisible: false,
      form: {},
      pager:{
        limit:20,
        current:1,
        total:0
      },
      tableData: [],
      loading: false,
      multipleSelection: []
    }
  },
  created() {
    this.form = this.newProjectForm()
    this.getProjects()
  },
  methods: {
    getProjects(){
      this.loading = true
      fetchProjects(this.pager.current).then(resp => {
        var newData = []
        resp.data.code == 0 && resp.data.data.forEach((e) => {
          newData.push({
            id: e.ID,
            name: e.name,
            domain: e.domain,
            config: e.config,
            listen: e.listen,
            process_id: e.process_id
          })
        })
        this.tableData = newData
        this.loading = false
      })
    },
    newProjectForm(){
      return {
        name: '',
        domain: '',
        plugins: [],
        checkAll: false
      }
    },
    toggleSelection(rows) {
      if (rows) {
        rows.forEach(row => {
          this.$refs.multipleTable.toggleRowSelection(row);
        });
      } else {
        this.$refs.multipleTable.clearSelection();
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    tableRowClassName({row}) {
      return this.isRunning(row) ? 'success-row' : '';
    },
    isRunning(row){
      return row.listen > 30000 && row.process_id > 0
    },
    handleDetail(_, row) {
      console.log(row);
    },
    handleResult(_, row) {
      this.$router.push({ path: `/vuls/${row.id}`})
    },
    handleToggle(_, row) {
      let p;
      if(this.isRunning(row)){
        p = stopProject(row.id)
      }else{
        p = startProject(row.id)
      }
      p.then(resp => {
        resp.data.code == 0 && this.getProjects()
      })
    },
    handleCurrentChange(val) {
      console.log(`当前页: ${val}`);
    },
    createProjectDialogOpen(){
      this.form = this.newProjectForm()
    },
    handleCheckAllChange(val) {
      this.form.plugins = val ? pluginOptions : [];
    },
    handleCheckedPluginsChange(value) {
      this.form.checkAll = value.length === this.pluginOptions.length;
    }
  }
}
</script>

<style >
  .projects {
    background: #fff;
    padding: 15px;
  }
  .el-table .success-row {
    background: #cdeebb;
  }
</style>