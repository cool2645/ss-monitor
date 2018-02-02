<template>
    <div>
        <!-- Content Header (Page header) -->
        <section class="content-header">
            <h1>
                任务记录
                <small>Task History</small>
            </h1>
        </section>
        <!-- Content -->
        <section class="content">
            <div class="row">
                <div class="col-sm-6">
                    <label for="worker_type" class="control-label">任务类型</label>
                    <select @change="updateData(1)" v-model="workerType" class="form-control" id="worker_type">
                        <option value="%">全部</option>
                        <option value="manager">Manager</option>
                        <option value="watcher">Watcher</option>
                        <option value="tester">Tester</option>
                        <option value="cleaner">Cleaner</option>
                    </select>
                </div>
                <div v-if="!noIpVer" class="col-sm-6">
                    <label v-if="!noIpVer" for="ip_ver" class="control-label">IP 协议</label>
                    <select @change="updateData(1)" v-model="ipVer" class="form-control" id="ip_ver">
                        <option value="%">不限</option>
                        <option value="4">IPv4</option>
                        <option value="6">IPv6</option>
                    </select>
                </div>
                <div class="col-sm-6">
                    <label for="node" class="control-label">节点</label>
                    <select @change="updateData(1)" v-model="nodeId" class="form-control" id="node">
                        <option value="%">不限</option>
                        <option v-for="node in nodes" :value="node.ID">{{ node.Name }}</option>
                    </select>
                </div>
            </div>
            <div class="table-responsive" style="margin-top: 20px">
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <th>运行 ID</th>
                        <th>节点/服务器名</th>
                        <th v-if="!noIpVer">{{ workerType === 'tester' ? 'IP 协议' : '类型/IP 协议' }}</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr v-for="task in data">
                        <td><a :href="'#/task/' + task.ID">{{ '#' + task.ID }}</a></td>
                        <td>{{ task.Node.Name || task.ServerName }}</td>
                        <td v-if="!noIpVer">{{ task.Class === 'tester' ? task.IPVer : task.Class + '/' + task.IPVer }}</td>
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
            <pagination :data="laravelData" :limit=2 v-on:pagination-change-page="updateData"></pagination>
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
        </section>
        <!-- /.content -->
    </div><!-- /.content-wrapper -->
</template>

<script>
    import config from './config'
    import urlParam from './buildUrlParam'

    export default {
        data() {
            return {
                jsonSource: {},
                nodes: [],
                workerType: '%',
                ipVer: '%',
                nodeId: '%',
                page: 1,
                perPage: 10,
                noIpVer: false
            }
        },
        computed: {
            data() {
                return this.jsonSource.result ? this.jsonSource.data.data : []
            },
            total() {
                return this.jsonSource.result ? this.jsonSource.data.total : 0;
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
        mounted() {
            this.getNodes();
            this.updateData(this.page)
        },
        methods: {
            updateData(page) {
                this.page = page;
                this.noIpVer = this.workerType !== '%' && this.workerType !== 'tester';
                let vm = this;
                fetch(config.urlPrefix + '/task?' + urlParam({
                    class: this.workerType,
                    ip_ver: this.ipVer,
                    node_id: this.nodeId,
                    page: this.page,
                    order: 'desc'
                }) )
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSource = res
                                }
                            }
                        )
                    });
            },
            getNodes() {
                let vm = this;
                fetch(config.urlPrefix + '/node?')
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.nodes = res.data
                                }
                            }
                        )
                    });
            }
        }
    }
</script>

<style lang="scss">
    .col-sm-6 {
        margin-top: 5px;
        margin-bottom: 5px;
    }
    td {
        vertical-align: middle !important;
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