<template>
    <div>
        <!-- Content Header (Page header) -->
        <section class="content-header">
            <h1>
                任务状态
                <small>Task Status</small>
            </h1>
        </section>
        <!-- Content -->
        <section class="content">
            <div class="table-responsive" style="margin-top: 20px">
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <th>运行 ID</th>
                        <th>节点/服务器名</th>
                        <th v-if="!noIpVer">IP 协议</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr>
                        <td><a :href="'#/task/' + task.ID">{{ '#' + task.ID }}</a></td>
                        <td>{{ task.Node.Name || task.ServerName }}</td>
                        <td v-if="!noIpVer">{{ task.Class === 'tester' ? task.IPVer : '' }}</td>
                        <td>
                            <a :href="'#/task/' + task.ID" v-if="task.State === 'Queuing'" class="btn btn-info">{{
                                task.State }}</a>
                            <a :href="'#/task/' + task.ID"
                               v-else-if="task.State === 'Passing' || task.State === 'Shiny☆'" class="btn btn-success">{{
                                task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Failing'" class="btn btn-danger">{{
                                task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else class="btn btn-warning">{{ task.State }}</a>
                        </td>
                        <td>
                            <a :href="'#/task/' + task.ID" v-if="task.Worker" class="btn btn-danger">{{ task.Worker
                                }}</a>
                            <p v-else>未指定</p>
                        </td>
                        <td>{{ task.CreatedAt }}</td>
                        <td>{{ task.UpdatedAt }}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
            <div class="row">
                <div class="col-xs-12">
                    <div class="box box-primary">
                        <div class="box-header">
                            <i class="fa fa-list"></i>
                            <h3 class="box-title">任务日志</h3>
                        </div>
                        <div class="box-body">
                            <pre>{{ task.Log }}</pre>
                        </div>
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-list"></i>
                            <h3 class="box-title">详细信息</h3>
                        </div>
                        <div class="box-body">
                            <tree-view :data="jsonSource" :options="{maxDepth: 3}"></tree-view>
                        </div>
                    </div>
                </div>
            </div>
        </section>
        <!-- /.content -->
    </div><!-- /.content-wrapper -->
</template>

<script>
    import config from './config'

    export default {
        data() {
            return {
                jsonSource: {},
            }
        },
        mounted() {
            this.updateData()
        },
        computed: {
            task() {
                return this.jsonSource.data
            }
        },
        methods: {
            updateData() {
                let vm = this;
                fetch(config.urlPrefix + '/task/' + this.$route.params.id)
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSource = res
                                }
                            }
                        )
                    });
            }
        },
        watch: {
            '$route'(to, from) {
                this.updateData()
            }
        }
    }
</script>

<style scoped>
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        white-space: -moz-pre-wrap;
    }
</style>