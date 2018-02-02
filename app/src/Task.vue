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
                        <th>任务类型</th>
                        <th v-if="!noIpVer">IP 协议</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr v-if="task">
                        <td><a :href="'#/task/' + task.ID">{{ '#' + task.ID }}</a></td>
                        <td>{{ task.Node.Name || task.ServerName }}</td>
                        <td>{{ task.Class }}</td>
                        <td v-if="!noIpVer">{{ task.Class === 'tester' ? task.IPVer : '' }}</td>
                        <td>
                            <a :href="'#/task/' + task.ID" v-if="task.State === 'Queuing'" class="btn btn-info">{{
                                task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Passing'" class="btn btn-success">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Shiny☆'" class="btn btn-shiny">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Failing'" class="btn btn-danger">{{ task.State }}</a>
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
                            <i class="fa fa-comment"></i>
                            <h3 class="box-title">任务结果</h3>
                        </div>
                        <div v-if="task" class="box-body">
                            <tree-view :data="result" :options="{maxDepth: 3}"></tree-view>
                        </div>
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col-xs-12">
                    <div class="box box-primary">
                        <div class="box-header">
                            <i class="fa fa-newspaper-o"></i>
                            <h3 class="box-title">任务日志</h3>
                        </div>
                        <div v-if="task" class="box-body">
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
                            <tree-view :data="jsonSource" :options="{maxDepth: 0}"></tree-view>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="isManager" class="table-responsive" style="margin-top: 20px">
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <th>运行 ID</th>
                        <th>节点/服务器名</th>
                        <th>类型/IP 协议</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr v-for="task in data">
                        <td><a :href="'#/task/' + task.ID">{{ '#' + task.ID }}</a></td>
                        <td>{{ task.Node.Name || task.ServerName }}</td>
                        <td>{{ task.Class + '/' + task.IPVer }}</td>
                        <td>
                            <a :href="'#/task/' + task.ID" v-if="task.State === 'Queuing'" class="btn btn-info">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Passing'" class="btn btn-success">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Shiny☆'" class="btn btn-shiny">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else-if="task.State === 'Failing'" class="btn btn-danger">{{ task.State }}</a>
                            <a :href="'#/task/' + task.ID" v-else class="btn btn-warning">{{ task.State }}</a>
                        </td>
                        <td>
                            <a :href="'#/task/' + task.ID" v-if="task.Worker" class="btn btn-danger">{{ task.Worker }}</a>
                            <p v-else>未指定</p>
                        </td>
                        <td>{{ task.CreatedAt }}</td>
                        <td>{{ task.UpdatedAt }}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
            <pagination v-if="isManager" :data="laravelData" :limit=2 v-on:pagination-change-page="updateDataChildren"></pagination>
            <div v-if="isManager" class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-list"></i>
                            <h3 class="box-title">详细信息</h3>
                        </div>
                        <div class="box-body">
                            <tree-view :data="jsonSourceChildren" :options="{maxDepth: 0}"></tree-view>
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
                jsonSourceChildren: {},
                isManager: false,
                page: 1,
                perPage: 10,
            }
        },
        mounted() {
            this.updateData()
        },
        computed: {
            noIpVer() {
                return this.jsonSource.result ? this.jsonSource.data.Class !== 'tester' : false
            },
            task() {
                return this.jsonSource.data
            },
            result() {
                try {
                    return JSON.parse(this.jsonSource.data.Result)
                } catch(e) {
                    return {}
                }
            },
            data() {
                return this.jsonSourceChildren.result ? this.jsonSourceChildren.data.data : []
            },
            total() {
                return this.jsonSourceChildren.result ? this.jsonSourceChildren.data.total : 0;
            },
            laravelData() {
                return {
                    current_page: this.page,
                    data: [],
                    from: (this.page - 1) * this.perPage + 1,
                    last_page: Math.ceil(this.total / this.perPage),
                    next_page_url: null,
                    per_page: this.perPage,
                    prev_page_url: null,
                    to: (this.page) * this.perPage,
                    total: this.total,
                }
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
                                    vm.jsonSource = res;
                                    if (res.data.Class === 'manager')
                                        this.updateDataChildren(1)
                                }
                            }
                        )
                    });
            },
            updateDataChildren(page) {
                this.page = page;
                let vm = this;
                fetch(config.urlPrefix + '/task?' + urlParam({
                    callback_id: this.task.ID,
                    page: this.page,
                    order: 'desc'
                }) )
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSourceChildren = res;
                                    vm.isManager = true;
                                }
                            }
                        )
                    });
            },
        },
        watch: {
            '$route'(to, from) {
                this.updateData()
            }
        }
    }
</script>

<style lang="scss" scoped>
    td {
        vertical-align: middle !important;
    }
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        white-space: -moz-pre-wrap;
    }
    .pagination {
        margin-top: 0;
    }
    .btn-shiny {
        color: #fff;
        background-color: #b76add;
        border-color: #a761ca;
        &:hover {
             color: #fff;
             background-color: #a761ca;
             border-color: #a761ca;
        }
    }
</style>