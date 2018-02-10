<template>
    <div>
        <!-- Content Header (Page header) -->
        <section class="content-header">
            <h1>
                任务状态
                <small>Task Status</small>
                <a v-if="sharable" href="javascript:;" @click="share"><i class="icon fa fa-share-alt pull-right"></i></a>
            </h1>
        </section>
        <!-- Content -->
        <section class="content">
            <div class="row">
                <div class="col-md-12">
                    <div id="msg-success" class="alert alert-success alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-info"></i> 成功!</h4>

                        <p id="msg-success-p"></p>
                    </div>
                    <div id="msg-warning" class="alert alert-warning alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-warning-p">{{ warning }}</p>
                    </div>
                    <div id="msg-error" class="alert alert-danger alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-error-p">电波无法送达哦～<a href="javascript:;" @click="startClocking">重试</a></p>
                    </div>
                </div>
            </div>
            <div class="table-responsive" style="margin-top: 20px" data-slideout-ignore>
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <th>运行 ID</th>
                        <th v-if="!isManager">节点/服务器名</th>
                        <th>任务类型</th>
                        <th v-if="!noIpVer">IP 协议</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th v-if="isAdmin">重置</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr v-if="task">
                        <td><a href="javascript:;" @click="location('/task/' + task.ID)">{{ '#' + task.ID }}</a></td>
                        <td v-if="!isManager">{{ task.Node.Name || task.ServerName }}</td>
                        <td>{{ task.Class }}</td>
                        <td v-if="!noIpVer">{{ task.Class === 'tester' ? task.IPVer : '' }}</td>
                        <td>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-if="task.State === 'Queuing'" class="btn btn-info">{{
                                task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Passing'" class="btn btn-success">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Shiny☆'" class="btn btn-shiny">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Failing'" class="btn btn-danger">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else class="btn btn-warning">{{ task.State }}</a>
                        </td>
                        <td>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-if="task.Worker" class="btn btn-danger">{{ task.Worker
                                }}</a>
                            <p v-else>未指定</p>
                        </td>
                        <td v-if="isAdmin"><a href="javascript:;" @click="resetTask(task.ID)" class="btn btn-danger">重置</a></td>
                        <td>{{ formatDateTimeFromDatetimeString(task.CreatedAt) }}</td>
                        <td>{{ formatDateTimeFromDatetimeString(task.UpdatedAt) }}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
            <div v-if="!isManager" class="row">
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
            <div v-if="!isManager" class="row">
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
            <div v-if="isManager" class="table-responsive" style="margin-top: 20px" data-slideout-ignore>
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <th>运行 ID</th>
                        <th>节点/服务器名</th>
                        <th>类型/IP 协议</th>
                        <th>运行结果</th>
                        <th>运行 Host</th>
                        <th v-if="isAdmin">重置</th>
                        <th>创建时间</th>
                        <th>更新时间</th>
                    </tr>
                    <tr v-for="task in data">
                        <td><a href="javascript:;" @click="location('/task/' + task.ID)">{{ '#' + task.ID }}</a></td>
                        <td>{{ task.Node.Name || task.ServerName }}</td>
                        <td>{{ task.Class + '/' + task.IPVer }}</td>
                        <td>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-if="task.State === 'Queuing'" class="btn btn-info">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Passing'" class="btn btn-success">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Shiny☆'" class="btn btn-shiny">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else-if="task.State === 'Failing'" class="btn btn-danger">{{ task.State }}</a>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-else class="btn btn-warning">{{ task.State }}</a>
                        </td>
                        <td>
                            <a href="javascript:;" @click="location('/task/' + task.ID)" v-if="task.Worker" class="btn btn-danger">{{ task.Worker }}</a>
                            <p v-else>未指定</p>
                        </td>
                        <td v-if="isAdmin"><a href="javascript:;" @click="resetTask(task.ID)" class="btn btn-danger">重置</a></td>
                        <td>{{ formatDateTimeFromDatetimeString(task.CreatedAt) }}</td>
                        <td>{{ formatDateTimeFromDatetimeString(task.UpdatedAt) }}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
            <pagination v-if="isManager" :data="laravelData" :limit=2 v-on:pagination-change-page="onPageChange"></pagination>
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
    import urlParam from './buildUrlParam'
    import formatDateTimeFromDatetimeString from "./datetimeUtil"

    export default {
        props: [
            "auth"
        ],
        data() {
            return {
                jsonSource: {},
                jsonSourceChildren: {},
                isManager: false,
                page: this.$route.hash.substr(1) || 1,
                perPage: 10,
                clock: null,
                warning: "",
                error: ""
            }
        },
        mounted() {
            this.startClocking()
        },
        beforeDestroy() {
            clearInterval(this.clock)
        },
        watch: {
            $route() {
                this.jsonSource = {};
                this.jsonSourceChildren = {};
                this.isManager = false;
                this.page = this.$route.hash.substr(1) || 1;
                this.updateData();
            }
        },
        computed: {
            sharable() {
                return navigator.share;
            },
            isAdmin() {
                return this.auth.isLogin && this.auth.user.privilege === "admin"
            },
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
            location(href) {
                this.$router.push({path: href});
            },
            startClocking() {
                $("#msg-error").hide(10);
                this.clock = setInterval(this.updateData, 5000);
                this.updateData();
            },
            dismissAlert() {
                $("#msg-success").hide(10);
                $("#msg-warning").hide(10);
                $("#msg-error").hide(10);
            },
            updateData() {
                let vm = this;
                fetch(config.urlPrefix + '/task/' + this.$route.params.id)
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSource = res;
                                    if (res.data.Class === 'manager')
                                        this.updateDataChildren()
                                }
                            }
                        )
                    })
                    .catch(error => {
                        $("#msg-error").hide(10).show(100);
                        clearInterval(this.clock);
                    });
            },
            onPageChange(page) {
                this.page = page;
                this.$router.push({hash: '#' + page});
                this.updateDataChildren();
            },
            updateDataChildren() {
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
                    })
                    .catch(error => {
                        $("#msg-error").hide(10).show(100);
                        clearInterval(this.clock);
                    });
            },
            resetTask(id) {
                let r = confirm("确定要重置任务吗？");
                if (!r) return;
                let vm = this;
                fetch(config.urlPrefix + '/task/' + id, {
                    credentials: 'include',
                    method: "DELETE",
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    $("#msg-warning").hide(100);
                                    $("#msg-success").show(100);
                                    this.updateData()
                                } else if (res.code === 401) {
                                    vm.warning = "登录超时：" + res.msg;
                                    $("#msg-warning").hide(10).show(100);
                                    this.$emit('check-auth');
                                    setTimeout(()=>{ vm.location("/admin") }, 2000);
                                } else {
                                    vm.warning = "发生错误：" + res.msg;
                                    $("#msg-warning").hide(10).show(100);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        vm.warning = "发生错误：" + res.msg;
                        $("#msg-warning").hide(10).show(100);
                    });
            },
            share() {
                navigator.share({
                    title: 'Task #' + this.$route.params.id,
                    text: config.shareText,
                    url: document.location.href,
                })
                    .then(() => console.log('Successful share'))
                    .catch((error) => console.log('Error sharing', error));
            },
            formatDateTimeFromDatetimeString(time) {
                return formatDateTimeFromDatetimeString(time);
            }
        },
    }
</script>

<style lang="scss" scoped>
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        white-space: -moz-pre-wrap;
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
    .content {
        padding-top: 0;
    }
</style>